package chaincode

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/hyperledger-labs/cc-tools-demo/ccapi/common"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	cb "github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
)

// fakeGateway emulates the peer Gateway service, responding with a fixed chaincode response.
type fakeGateway struct {
	gateway.UnimplementedGatewayServer
	publicKey *ecdsa.PublicKey
	response  *peer.Response
	commit    peer.TxValidationCode
}

func (f *fakeGateway) chaincodeError() error {
	if f.response.Status < 400 {
		return nil
	}

	msg := fmt.Sprintf("chaincode response %d, %s", f.response.Status, f.response.Message)
	st, _ := status.New(codes.Unknown, msg).WithDetails(&gateway.ErrorDetail{Address: "peer0:7051", MspId: "orgMSP", Message: msg})
	return st.Err()
}

func (f *fakeGateway) Evaluate(ctx context.Context, req *gateway.EvaluateRequest) (*gateway.EvaluateResponse, error) {
	digest := sha256.Sum256(req.GetProposedTransaction().GetProposalBytes())
	if !ecdsa.VerifyASN1(f.publicKey, digest[:], req.GetProposedTransaction().GetSignature()) {
		return nil, status.Error(codes.PermissionDenied, "invalid proposal signature")
	}

	if err := f.chaincodeError(); err != nil {
		return nil, err
	}

	return &gateway.EvaluateResponse{Result: f.response}, nil
}

func (f *fakeGateway) Endorse(ctx context.Context, req *gateway.EndorseRequest) (*gateway.EndorseResponse, error) {
	if err := f.chaincodeError(); err != nil {
		return nil, err
	}

	envelope, err := preparedEnvelope("mychannel", f.response)
	if err != nil {
		return nil, err
	}

	return &gateway.EndorseResponse{PreparedTransaction: envelope}, nil
}

func (f *fakeGateway) Submit(ctx context.Context, req *gateway.SubmitRequest) (*gateway.SubmitResponse, error) {
	return &gateway.SubmitResponse{}, nil
}

func (f *fakeGateway) CommitStatus(ctx context.Context, req *gateway.SignedCommitStatusRequest) (*gateway.CommitStatusResponse, error) {
	return &gateway.CommitStatusResponse{Result: f.commit, BlockNumber: 1}, nil
}

func preparedEnvelope(channelName string, response *peer.Response) (*cb.Envelope, error) {
	marshal := func(m proto.Message, err error) ([]byte, error) {
		if err != nil {
			return nil, err
		}
		return proto.Marshal(m)
	}

	b, err := marshal(&peer.ChaincodeAction{Response: response}, nil)
	b, err = marshal(&peer.ProposalResponsePayload{Extension: b}, err)
	b, err = marshal(&peer.ChaincodeActionPayload{Action: &peer.ChaincodeEndorsedAction{ProposalResponsePayload: b}}, err)
	b, err = marshal(&peer.Transaction{Actions: []*peer.TransactionAction{{Payload: b}}}, err)
	data := b
	b, err = marshal(&cb.ChannelHeader{ChannelId: channelName}, err)
	b, err = marshal(&cb.Payload{Header: &cb.Header{ChannelHeader: b}, Data: data}, err)
	if err != nil {
		return nil, err
	}

	return &cb.Envelope{Payload: b}, nil
}

type gatewayFixture struct {
	conn     *grpc.ClientConn
	contract *client.Contract
	sign     identity.Sign
	gateway  *fakeGateway
}

func newGatewayFixture(t *testing.T) *gatewayFixture {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "Admin"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatal(err)
	}

	id, err := identity.NewX509Identity("orgMSP", certificate)
	if err != nil {
		t.Fatal(err)
	}
	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	fake := &fakeGateway{publicKey: &privateKey.PublicKey, commit: peer.TxValidationCode_VALID}

	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	gateway.RegisterGatewayServer(server, fake)
	go server.Serve(listener)
	t.Cleanup(server.Stop)

	conn, err := grpc.Dial("bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return listener.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })

	gw, err := client.Connect(id, client.WithSign(sign), client.WithClientConnection(conn))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { gw.Close() })

	return &gatewayFixture{
		conn:     conn,
		contract: gw.GetNetwork("mychannel").GetContract("cc"),
		sign:     sign,
		gateway:  fake,
	}
}

func (f *gatewayFixture) evaluate(t *testing.T) ([]byte, int, error) {
	proposal, err := f.contract.NewProposal("getBooksByAuthor", client.WithArguments(`{"authorName":"x"}`))
	if err != nil {
		t.Fatal(err)
	}

	return evaluate(f.conn, "mychannel", f.sign, proposal)
}

func TestEvaluateStatus(t *testing.T) {
	f := newGatewayFixture(t)

	cases := []struct {
		response *peer.Response
		status   int
		payload  string
	}{
		{&peer.Response{Status: 200, Payload: []byte(`[]`)}, 200, `[]`},
		{&peer.Response{Status: 201, Payload: []byte(`{"a":1}`)}, 201, `{"a":1}`},
		{&peer.Response{Status: 204}, 204, ``},
	}

	for _, c := range cases {
		f.gateway.response = c.response
		payload, st, err := f.evaluate(t)
		if err != nil {
			t.Fatalf("status %d: unexpected error %v", c.response.Status, err)
		}
		if st != c.status || string(payload) != c.payload {
			t.Errorf("expected %d %q, got %d %q", c.status, c.payload, st, payload)
		}
	}
}

func TestEvaluateError(t *testing.T) {
	f := newGatewayFixture(t)
	f.gateway.response = &peer.Response{Status: 404, Message: "asset not found"}

	_, _, err := f.evaluate(t)
	if err == nil {
		t.Fatal("expected error")
	}

	err, st := common.ParseError(err)
	if st != 404 || err.Error() != "asset not found" {
		t.Errorf("expected 404 %q, got %d %q", "asset not found", st, err.Error())
	}
}

func TestSubmitStatus(t *testing.T) {
	f := newGatewayFixture(t)
	f.gateway.response = &peer.Response{Status: 201, Payload: []byte(`{"@key":"book:1"}`)}

	payload, st, err := submit(f.contract, "createAsset", client.WithArguments(`{}`))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if st != 201 || string(payload) != `{"@key":"book:1"}` {
		t.Errorf("expected 201, got %d %q", st, payload)
	}
}

func TestSubmitErrors(t *testing.T) {
	f := newGatewayFixture(t)

	f.gateway.response = &peer.Response{Status: 409, Message: "asset already exists"}
	_, _, err := submit(f.contract, "createAsset", client.WithArguments(`{}`))
	if err == nil {
		t.Fatal("expected endorse error")
	}
	if err, st := common.ParseError(err); st != 409 || err.Error() != "asset already exists" {
		t.Errorf("expected 409 %q, got %d %q", "asset already exists", st, err.Error())
	}

	f.gateway.response = &peer.Response{Status: 200, Payload: []byte(`{}`)}
	f.gateway.commit = peer.TxValidationCode_MVCC_READ_CONFLICT
	_, _, err = submit(f.contract, "createAsset", client.WithArguments(`{}`))
	commitErr, ok := err.(*client.CommitError)
	if !ok || commitErr.Code != peer.TxValidationCode_MVCC_READ_CONFLICT {
		t.Fatalf("expected commit error, got %v", err)
	}
}

func TestChaincodeStatusFromPreparedTransaction(t *testing.T) {
	envelope, err := preparedEnvelope("mychannel", &peer.Response{Status: 202})
	if err != nil {
		t.Fatal(err)
	}

	b, err := proto.Marshal(&gateway.PreparedTransaction{TransactionId: "tx1", Envelope: envelope})
	if err != nil {
		t.Fatal(err)
	}

	st, err := chaincodeStatusFromPreparedTransaction(b)
	if err != nil || st != 202 {
		t.Errorf("expected 202, got %d %v", st, err)
	}

	_, err = chaincodeStatusFromPreparedTransaction([]byte("invalid"))
	if err == nil {
		t.Error("expected error for invalid transaction")
	}
}
