package fetcher

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hmluck83/txlens-srv/schemas"
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

func FetchTransaction(txhash string) (*schemas.Profile, *schemas.AddressLabels, error) {
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

	var profile schemas.Profile
	profileDecoder := json.NewDecoder(profileResp.Body)

	if err := profileDecoder.Decode(&profile); err != nil {
		return nil, nil, err
	}

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

	var addresslabel schemas.AddressLabels
	labelDecoder := json.NewDecoder(labelResq.Body)

	if err = labelDecoder.Decode(&addresslabel); err != nil {
		return nil, nil, err
	}

	return &profile, &addresslabel, nil
}
