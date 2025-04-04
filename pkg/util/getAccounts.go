package util

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"main/pkg/types"
)

func privateKeyToAddress(privateKey *ecdsa.PrivateKey) (*common.Address, error) {
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKey)
	return &address, nil
}

func GetAccounts(
	inputs []string,
	onlyKeys bool,
) ([]types.AccountData, error) {
	var accounts []types.AccountData

	for _, input := range inputs {
		input = RemoveHexPrefix(input)

		if common.IsHexAddress("0x" + input) {
			if onlyKeys {
				log.Printf("%s | Address, Not Private Key", input)
			} else {
				accounts = append(accounts, types.AccountData{
					AccountLogData: "0x" + input,
					AccountKeyHex:  "",
					AccountKey:     nil,
					AccountAddress: common.HexToAddress("0x" + input),
				})
			}

			continue
		}

		privateKey, err := crypto.HexToECDSA(input)
		if err != nil {
			log.Printf("%s | Invalid Private Key", input)
			continue
		}

		sweepedAddress, err := privateKeyToAddress(privateKey)
		if err != nil {
			log.Printf("%s | Failed To Derive Address", input)
			continue
		}

		accounts = append(accounts, types.AccountData{
			AccountLogData: input,
			AccountKeyHex:  "",
			AccountKey:     privateKey,
			AccountAddress: *sweepedAddress,
		})
	}

	return accounts, nil
}
