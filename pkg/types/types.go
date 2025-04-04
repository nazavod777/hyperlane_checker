package types

import "crypto/ecdsa"
import "github.com/ethereum/go-ethereum/common"

type AccountData struct {
	AccountKeyHex  string
	AccountKey     *ecdsa.PrivateKey
	AccountAddress common.Address
	AccountLogData string
}
