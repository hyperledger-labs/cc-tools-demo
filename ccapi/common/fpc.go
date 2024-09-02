/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"fmt"
	"os"
	"strconv"

	pkgFpc "github.com/hyperledger-labs/ccapi/fpcUtils"
)

var (
	fpcconfig *pkgFpc.Config
)

func InitFpcEnclave() error {
	admin := pkgFpc.NewAdmin(fpcconfig)
	defer admin.Close()
	return admin.InitEnclave(fpcconfig.CorePeerId)
}

func InitFpcConfig() {

	getStrEnv := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			panic(fmt.Sprintf("%s not set", key))
		}
		return val
	}

	getBoolEnv := func(key string) bool {
		val := getStrEnv(key)
		ret, err := strconv.ParseBool(val)
		if err != nil {
			if val == "" {
				panic(fmt.Sprintf("invalid bool value for %s", key))
			}
		}
		return ret
	}

	fpcconfig = &pkgFpc.Config{
		CorePeerAddress:         getStrEnv("CORE_PEER_ADDRESS"),
		CorePeerId:              getStrEnv("CORE_PEER_ID"),
		CorePeerLocalMSPID:      getStrEnv("CORE_PEER_LOCALMSPID"),
		CorePeerMSPConfigPath:   getStrEnv("CORE_PEER_MSPCONFIGPATH"),
		CorePeerTLSCertFile:     getStrEnv("CORE_PEER_TLS_CERT_FILE"),
		CorePeerTLSEnabled:      getBoolEnv("CORE_PEER_TLS_ENABLED"),
		CorePeerTLSKeyFile:      getStrEnv("CORE_PEER_TLS_KEY_FILE"),
		CorePeerTLSRootCertFile: getStrEnv("CORE_PEER_TLS_ROOTCERT_FILE"),
		OrdererCA:               getStrEnv("ORDERER_CA"),
		ChaincodeId:             getStrEnv("CCNAME"),
		ChannelId:               getStrEnv("CHANNEL"),
		GatewayConfigPath:       getStrEnv("GATEWAY_CONFIG"),
	}

}

func NewFpcClient() *pkgFpc.Client {
	// fpcconfig.ChaincodeId=
	// fpcconfig.ChannelId=
	return pkgFpc.NewClient(fpcconfig)
}
