package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const defaultNode = "http://localhost:8110"

func makeRequest(url string, method string, data map[string]interface{}) ([]byte, error) {

	if data == nil {
		data = make(map[string]interface{})
	}

	jsonBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
