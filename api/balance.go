package api

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

type BalanceRes struct {
	Balance uint `json:"balance"`
}

const endtpointAddressBalance = "/address/balance"

func GetBalance(account common.Address) (uint, error) {

	res := BalanceRes{}

	balanceRawBody, err := makeRequest(defaultNode+endtpointAddressBalance, "POST", map[string]interface{}{
		"account": account.Hex(),
	})
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(balanceRawBody, &res)
	if err != nil {
		return 0, err
	}

	return res.Balance, nil
}
