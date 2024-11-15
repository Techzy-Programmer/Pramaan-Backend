package utils

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifySignature(pubAddr string, message string, sig string) bool {
	hash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)))

	sigBytes := common.FromHex(sig)
	if len(sigBytes) != 65 {
		return false
	}

	r := sigBytes[:32]
	s := sigBytes[32:64]
	v := sigBytes[64]

	if v == 27 || v == 28 {
		v -= 27
	}

	publicKey, err := crypto.SigToPub(hash.Bytes(), append(r, append(s, v)...))
	if err != nil {
		fmt.Printf("Failed to recover public key: %v\n", err)
		return false
	}

	recoveredAddress := crypto.PubkeyToAddress(*publicKey)
	return bytes.Equal(recoveredAddress.Bytes(), common.HexToAddress(pubAddr).Bytes())
}
