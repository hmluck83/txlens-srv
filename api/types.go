package api

type ReqeustSummury struct {
	Chainid int    `json:"chainID"`
	Txid    string `json:"txid"`
}

type ResponseSummury struct {
	StatusMessage string    `json:"statusMessage"`
	Status        int       `json:"status"`
	Description   string    `json:"description"`
	GraphData     GraphData `json:"graphdata"`
}

type GraphData struct {
	Nodes []NodeData `json:"nodes"`
	Edges []EdgeData `json:"edges"`
}

type NodeData struct {
	Data GraphNode `json:"data"`
}

type GraphNode struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Tag  string `json:"tag,omitempty"`
}

type EdgeData struct {
	Data GraphEdge `json:"data"`
}

type GraphEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Name   string `json:"name"`
}

// LLM에 질의 하기 위한 FundFlow Description -> FundFlow Structure 를 수정하여 생성
type FundDesc struct {
	Amount     string `json:"amount"`
	From       string `json:"from"`
	FromLabel  string `json:"fromLabel"`
	ID         int    `json:"id"`
	IsERC1155  bool   `json:"isERC1155"`
	IsERC721   bool   `json:"isERC721"`
	IsReverted bool   `json:"isReverted"`
	Order      int    `json:"order"`
	To         string `json:"to"`
	ToLabel    string `json:"toLabel"`
	Token      string `json:"token"`
}

type summaryProfile struct {
	Sender        string `json:"sender"`
	Status        bool   `json:"status"`
	RevertMessage string `json:"revertMessage"`
}

type summaryDesc struct {
	SummaryProfile summaryProfile `json:"summaryProfile"`
	FundDesc       []FundDesc     `json:"fundDesc"`
}
