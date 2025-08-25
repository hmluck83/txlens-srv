package tracer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type fundFlow struct {
	From  common.Address
	To    common.Address
	Value big.Int
	Token common.Address
}

type FundFlows []fundFlow

type addrLabel struct {
	IsContract bool
	Label      string
	Name       string
	Decimals   uint8
	Symbol     string
}

type AddrLabels map[common.Address]*addrLabel

type arkhamResp struct {
	Address      string `json:"address"`
	Chain        string `json:"chain"`
	ArkhamEntity struct {
		Name       string `json:"name"`
		Note       string `json:"note"`
		ID         string `json:"id"`
		Type       string `json:"type"`
		Service    string `json:"service,omitzero"`
		Addresses  string `json:"addresses,omitzero"`
		Website    string `json:"website"`
		Twitter    string `json:"twitter"`
		Crunchbase string `json:"crunchbase"`
		Linkedin   string `json:"linkedin"`
	} `json:"arkhamEntity,omitzero"`
	ArkhamLabel struct {
		Name      string `json:"name"`
		Address   string `json:"address"`
		ChainType string `json:"chainType"`
	} `json:"arkhamLabel,omitzero"`
	IsUserAddress bool `json:"isUserAddress"`
	Contract      bool `json:"contract"`
}
