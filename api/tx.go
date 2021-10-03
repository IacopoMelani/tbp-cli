package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/IacopoMelani/the-blockchain-pub/database"
	"github.com/IacopoMelani/the-blockchain-pub/wallet"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
)

type NextNonceRes struct {
	Nonce uint `json:"nonce"`
}

func SendTx(key *keystore.Key, to string, amount uint) error {

	nextNonceRawBody, err := makeRequest("http://localhost:8110/address/nonce/next", "POST", map[string]interface{}{
		"account": key.Address.Hex(),
	})
	if err != nil {
		return err
	}

	var nextNonceRes NextNonceRes
	err = json.Unmarshal(nextNonceRawBody, &nextNonceRes)
	if err != nil {
		return err
	}

	tx := database.NewTx(key.Address, database.NewAccount(to), amount, nextNonceRes.Nonce, "")

	signedTx, err := wallet.SignTxWithKey(tx, key)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rawBytes, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rawTx := hexutil.Encode(rawBytes)

	_, err = makeRequest("http://localhost:8110/tx/add", "POST", map[string]interface{}{
		"tx": rawTx,
	})
	if err != nil {
		return err
	}

	return nil
}
