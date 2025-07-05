package fetcher

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

type payload struct {
	ChainID int    `json:"chainID"`
	TxnHash string `json:"txnHash"`
	Blocked bool   `json:"blocked"`
}

var profileURL = url.URL{
	Scheme: "https",
	Host:   "app.blocksec.com",
	Path:   "api/v1/onchain/tx/profile",
}

var labelURL = url.URL{
	Scheme: "https",
	Host:   "app.blocksec.com",
	Path:   "api/v1/onchain/tx/address-label",
}

func setHttpHeaders(req *http.Request) {
	req.Header.Add("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Thunder Client (https://www.thunderclient.com) ")
}

// 트랜잭션의 정보 읽어와야 함, 현재 chainId는 무시 Ethereum으로 고정되나 향후 조정 필요
func FetchTransaction(txhash string) (*Profile, *AddressLabels, error) {
	client := &http.Client{}

	payloadBody := &payload{
		ChainID: 1,
		TxnHash: txhash,
		Blocked: false,
	}

	marshaledPayload, err := json.Marshal(payloadBody)
	if err != nil {
		return nil, nil, err
	}
	bufPayload := bytes.NewBuffer(marshaledPayload)

	profileReq, err := http.NewRequest("POST", profileURL.String(), bufPayload)
	if err != nil {
		return nil, nil, err
	}

	setHttpHeaders(profileReq)

	profileResp, err := client.Do(profileReq)
	if err != nil {
		return nil, nil, err
	}
	defer profileResp.Body.Close() // Ensure the response body is closed

	var profile Profile
	profileDecoder := json.NewDecoder(profileResp.Body)

	if err := profileDecoder.Decode(&profile); err != nil {
		return nil, nil, err
	}

	// Request Address Label
	bufPayload = bytes.NewBuffer(marshaledPayload)
	labelReq, err := http.NewRequest("POST", labelURL.String(), bufPayload)
	if err != nil {
		return nil, nil, err
	}

	setHttpHeaders(labelReq)

	labelResq, err := client.Do(labelReq)
	if err != nil {
		return nil, nil, err
	}
	defer labelResq.Body.Close()

	var addresslabel AddressLabels
	labelDecoder := json.NewDecoder(labelResq.Body)

	if err = labelDecoder.Decode(&addresslabel); err != nil {
		return nil, nil, err
	}

	return &profile, &addresslabel, nil
}
