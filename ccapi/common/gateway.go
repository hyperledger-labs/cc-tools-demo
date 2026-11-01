package common

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var (
	gatewayTLSCredentials *credentials.TransportCredentials
)

// Default timeouts for different gRPC calls
const (
	GatewayEvaluateTimeout     = 5 * time.Second
	GatewayEndorseTimeout      = 15 * time.Second
	GatewaySubmitTimeout       = 5 * time.Second
	GatewayCommitStatusTimeout = 1 * time.Minute
)

func CreateGrpcConnection(endpoint string) (*grpc.ClientConn, error) {
	// Check TLS credential was created
	if gatewayTLSCredentials == nil {
		gatewayServerName := os.Getenv("FABRIC_GATEWAY_NAME")

		cred, err := createTransportCredential(GetTLSCACert(), gatewayServerName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create tls credentials")
		}

		gatewayTLSCredentials = &cred
	}

	// Create client grpc connection
	return grpc.Dial(endpoint, grpc.WithTransportCredentials(*gatewayTLSCredentials))
}

func CreateGatewayConnection(grpcConn *grpc.ClientConn, user string) (*client.Gateway, error) {
	// Create identity
	id, err := newIdentity(getSignCert(user), GetMSPID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new identity")
	}
	gatewayId := id

	// Create sign function
	sign, err := NewGatewaySign(user)
	if err != nil {
		return nil, err
	}

	gatewaySign := sign

	// Create a Gateway connection for a specific client identity.
	return client.Connect(
		gatewayId,
		client.WithSign(gatewaySign),
		client.WithClientConnection(grpcConn),

		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(GatewayEvaluateTimeout),
		client.WithEndorseTimeout(GatewayEndorseTimeout),
		client.WithSubmitTimeout(GatewaySubmitTimeout),
		client.WithCommitStatusTimeout(GatewayCommitStatusTimeout),
	)
}

// NewGatewaySign creates the function used to sign gateway requests on behalf of the user.
func NewGatewaySign(user string) (identity.Sign, error) {
	sign, err := newSign(getSignKey(user))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new sign function")
	}

	return sign, nil
}

// Create transport credential
func createTransportCredential(tlsCertPath, serverName string) (credentials.TransportCredentials, error) {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	return credentials.NewClientTLSFromCert(certPool, serverName), nil
}

// Creates a client identity for a gateway connection using an X.509 certificate.
func newIdentity(certPath, mspID string) (*identity.X509Identity, error) {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		return nil, err
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		return nil, err
	}

	return id, nil
}

// Creates a function that generates a digital signature from a message digest using a private key.
func newSign(keyPath string) (identity.Sign, error) {
	privateKeyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read private key file")
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create signer function")
	}

	return sign, nil
}

// Returns error and status code
func ParseError(err error) (error, int) {
	var errMsg string

	switch err := err.(type) {
	case *client.EndorseError:
		errMsg = "endorse error for transaction"
	case *client.SubmitError:
		errMsg = "submit error for transaction"
	case *client.CommitStatusError:
		if errors.Is(err, context.DeadlineExceeded) {
			errMsg = "timeout waiting for transaction commit status"
		} else {
			errMsg = "error obtaining commit status for transaction"
		}
	case *client.CommitError:
		errMsg = fmt.Sprintf("transaction failed to commit with status code %d (%s)", int32(err.Code), err.Code.String())
	default:
		errMsg = "unexpected error type:" + err.Error()
	}

	statusErr := status.Convert(err)

	details := statusErr.Details()
	if len(details) == 0 {
		return errors.New(errMsg), http.StatusInternalServerError
	}

	// Each endorser that failed adds a detail. Prefer the status returned by the
	// chaincode over errors from peers that could not be reached.
	var firstMsg string
	for _, detail := range details {
		switch detail := detail.(type) {
		case *gateway.ErrorDetail:
			status, msg, ok := extractStatusAndMessage(detail.Message)
			if ok {
				return errors.New(msg), status
			}
			if firstMsg == "" {
				firstMsg = msg
			}
		}
	}

	if firstMsg != "" {
		return errors.New(firstMsg), http.StatusInternalServerError
	}

	return errors.New(errMsg), http.StatusInternalServerError
}

// extractStatusAndMessage parses messages in the "chaincode response <status>, <message>" format
// produced by the Fabric Gateway. ok is false when the message is not in that format.
func extractStatusAndMessage(msg string) (status int, errMsg string, ok bool) {
	pattern := `chaincode response (\b(\d{3})\b), `
	reg := regexp.MustCompile(pattern)
	matches := reg.FindStringSubmatch(msg)

	if len(matches) == 0 {
		return http.StatusInternalServerError, msg, false
	}

	errMsg = strings.Replace(msg, matches[0], "", 1)
	status, err := strconv.Atoi(matches[1])
	if err != nil {
		return http.StatusInternalServerError, errMsg, false
	}

	return status, errMsg, true
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

func getSignCert(user string) string {
	cryptoPath := GetCryptoPath()
	filename := user + "@" + os.Getenv("ORG") + "." + os.Getenv("DOMAIN") + "-cert.pem"

	return strings.Replace(cryptoPath, "{username}", user, 1) + "/signcerts/" + filename
}

func getSignKey(user string) string {
	cryptoPath := GetCryptoPath()
	filename := "priv_sk"

	return strings.Replace(cryptoPath, "{username}", user, 1) + "/keystore/" + filename
}
