package tracer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
)

var (
	funcName     = w3.MustNewFunc("name()", "string")
	funcDecimals = w3.MustNewFunc("decimals()", "uint8")
	funcSymbol   = w3.MustNewFunc("symbol()", "string")

	noERC20Error = errors.New("Probably not an ERC-20")
)

func addressLabeler(ff *fundFlow, addrLabels AddrLabels) {

	addresses := []common.Address{ff.From, ff.To, ff.Token}
	for i := range addresses {
		if _, exist := addrLabels[addresses[i]]; !exist {
			// Handling the error later
			addrLabel, _ := labelAddress(addresses[i])
			addrLabels[addresses[i]] = addrLabel
		}

	}

}

func labelAddress(addr common.Address) (*addrLabel, error) {

	var label *addrLabel
	label, err := requestArkham(addr)
	// TODO: Check the error for arkham api availability.
	if err != nil {
		label = &addrLabel{}
	}

	var client = w3.MustDial(os.Getenv("RPCNODE"))

	var addressName string
	err = client.Call(
		eth.CallFunc(addr, funcName).Returns(&addressName),
	)

	if err != nil {
		return label, fmt.Errorf("%w, %w", noERC20Error, err)
	}

	var decimals uint8
	err = client.Call(
		eth.CallFunc(addr, funcDecimals).Returns(&decimals),
	)

	if err != nil {
		return label, fmt.Errorf("%w, %w", noERC20Error, err)
	}

	var addressSymbol string
	err = client.Call(
		eth.CallFunc(addr, funcSymbol).Returns(&addressSymbol),
	)

	if err != nil {
		return label, fmt.Errorf("%w, %w", noERC20Error, err)
	}

	label.Name = addressName
	label.Decimals = decimals
	label.Symbol = addressSymbol

	return label, nil
}

func requestArkham(addr common.Address) (*addrLabel, error) {

	// Chech Ethereum
	if addr.Cmp(EthAddress) == 0 {
		return &addrLabel{
			IsContract: false,
			Label:      "Ethereum",
			Name:       "Ethereym",
			Decimals:   18,
			Symbol:     "ETH",
		}, nil
	}

	urls := fmt.Sprintf("https://api.arkm.com/intelligence/address/%s?chain=ethereum", addr.Hex())
	req, err := http.NewRequest("GET", urls, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("API-Key", os.Getenv("ARKHAM"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)

	var arkhamObj arkhamResp
	err = json.Unmarshal(bytes, &arkhamObj)
	if err != nil {
		return nil, err
	}

	label := addrLabel{
		IsContract: arkhamObj.Contract,
	}

	if label.IsContract {
		label.Label = arkhamObj.ArkhamEntity.Name + ":" + arkhamObj.ArkhamLabel.Name
	}

	return &label, nil
}
