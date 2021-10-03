package api

import "encoding/json"

type StatusRes struct {
	BlockHash   string `json:"block_hash"`
	BlockNumber int64  `json:"block_number"`
}

func GetStatus() (StatusRes, error) {
	var res StatusRes

	rawBody, err := makeRequest("http://localhost;:8110/node/status", "GET", nil)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(rawBody, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
