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
