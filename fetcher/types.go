// 외부 RPC-Node 혹은 Explorer의 Response schema 정의
package fetcher

import "time"

// BlockSec Json Response
type Profile struct {
	BasicInfo struct {
		BaseFee               string    `json:"baseFee"`
		BlockNumber           int       `json:"blockNumber"`
		CallData              string    `json:"callData"`
		DebugEnable           bool      `json:"debugEnable"`
		ErrorInfo             string    `json:"errorInfo"`
		EventCount            int       `json:"eventCount"`
		FeeRecipient          string    `json:"feeRecipient"`
		FormatValue           string    `json:"formatValue"`
		GasLimit              int       `json:"gasLimit"`
		GasPrice              string    `json:"gasPrice"`
		GasUsed               int       `json:"gasUsed"`
		IntTxnCount           int       `json:"intTxnCount"`
		InternalRevert        bool      `json:"internalRevert"`
		InternalRevertMessage string    `json:"internalRevertMessage"`
		IsContractCreation    bool      `json:"isContractCreation"`
		IsFlashloan           bool      `json:"isFlashloan"`
		MaxFee                string    `json:"maxFee"`
		Nonce                 int       `json:"nonce"`
		OverwritedBlockNumber int       `json:"overwritedBlockNumber"`
		PriorityFee           string    `json:"priorityFee"`
		Receiver              string    `json:"receiver"`
		RevertMessage         string    `json:"revertMessage"`
		Sender                string    `json:"sender"`
		SimulationEnable      bool      `json:"simulationEnable"`
		SourceTxnHash         string    `json:"sourceTxnHash"`
		Status                bool      `json:"status"`
		Timestamp             time.Time `json:"timestamp"`
		TransactionFee        string    `json:"transactionFee"`
		TxnHash               string    `json:"txnHash"`
		TxnIndex              int       `json:"txnIndex"`
		Type                  int       `json:"type"`
		Value                 string    `json:"value"`
	} `json:"basicInfo"`
	SecurityEvents SecurityEvent `json:"securityEvent,omitempty"`
	FundFlows      []FundFlow    `json:"fundFlow"`
	TokenInfos     []TokenInfo   `json:"tokenInfos"`
}

type FundFlow struct {
	Amount     string `json:"amount"`
	From       string `json:"from"`
	ID         int    `json:"id"`
	IsERC1155  bool   `json:"isERC1155"`
	IsERC721   bool   `json:"isERC721"`
	IsReverted bool   `json:"isReverted"`
	Order      int    `json:"order"`
	To         string `json:"to"`
	Token      string `json:"token"`
}

type TokenInfo struct {
	Address  string `json:"address"`
	Decimals int    `json:"decimals"`
}

type AddressLabels struct {
	Labels []AddressLabel `json:"labels"`
}

type AddressLabel struct {
	Address string `json:"address"`
	Label   string `json:"label"`
}
type SecurityEvent struct {
	ID          int    `json:"id"`
	Project     string `json:"project"`
	ProjectLogo string `json:"projectLogo"`
	Loss        int    `json:"loss"`
	Media       string `json:"media"`
	RootCause   string `json:"rootCause"`
	Poc         string `json:"poc"`
	Rescued     int    `json:"rescued"`
}
