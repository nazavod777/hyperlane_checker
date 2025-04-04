package util

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"main/pkg/types"
)

func SignMessage(
	accountData types.AccountData,
	signText string,
) string {
	signature, err := crypto.Sign(accounts.TextHash([]byte(signText)), accountData.AccountKey)

	if err != nil {
		log.Panicf("%s | Failed To Sign Message: %v", accountData.AccountLogData, err)
	}

	signature[64] += 27
	signHash := hexutil.Encode(signature)

	return signHash
}
