package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hmluck83/txlens-srv/fetcher"
	"github.com/hmluck83/txlens-srv/llmclient"
)

// Transaction Summury를 작성 API
func summuryHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: Logger  적용
	// CORS 헤더 설정
	enableCORS(&w)

	if strings.ToUpper(r.Method) == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only Post requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqPayload ReqeustSummury
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPayload); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Fetch Transaction Infomation
	// TODO: verify chainId, Transaction ID
	profile, addressLabels, err := fetcher.FetchTransaction(reqPayload.Txid)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	summarizedObj := fetcher.Summarizer(*profile, *addressLabels)

	summarizedString, err := json.Marshal(summarizedObj)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	lc := llmclient.NewLLMClient()
	llmresponse, err := lc.Request(context.Background(), llmclient.GetScaffoldTemplate(), string(summarizedString))

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	graphData, _ := buildGraphData(profile, addressLabels)
	responseSummury := ResponseSummury{
		StatusMessage: "success",
		Status:        1,
		Description:   *llmresponse,
		GraphData:     *graphData,
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(responseSummury); err != nil {
		log.Printf("Error encoding response JSON: %s\n", err)
	}
}

// Transaction FundFlow와 Address Label을 기준으로 Cytoscape Graph Data 생성
func buildGraphData(profile *fetcher.Profile, addressLabels *fetcher.AddressLabels) (*GraphData, error) {

	// Address Labels to map
	labelMap := make(map[string]string)

	log.Printf("log start\n")
	for _, val := range addressLabels.Labels {

		log.Printf("%s, %s", val.Address, val.Label)
		// 기존에 키가있는지 검사 하는것보다 무지성으로 때려 박는게 더 빠르지 않을까?
		labelMap[val.Address] = val.Label
	}

	// node를 중복 없이 저장할 필요가 있지만 Set이 없기 때문에 Map으로 대ㅈ
	nodeMap := make(map[string]struct{})
	edges := []EdgeData{}

	for _, val := range profile.FundFlows {
		nodeMap[val.From] = struct{}{}
		nodeMap[val.To] = struct{}{}

		name := fmt.Sprintf("%s %s %s",
			strconv.Itoa(val.ID),
			shortenAmount(val.Amount),
			labelMap[val.Token],
		)

		edges = append(edges, EdgeData{
			GraphEdge{
				Source: val.From,
				Target: val.To,
				Name:   name,
			}})
	}

	nodes := []NodeData{}

	for key := range nodeMap {
		var name string
		if label, exist := labelMap[key]; exist {
			name = fmt.Sprintf("%s(%s)", label, shortenAddress(key))
		} else {
			name = shortenAddress(key)
		}

		nodes = append(nodes, NodeData{
			Data: GraphNode{
				Id:   key,
				Name: name,
			},
		})
	}

	return &GraphData{
		Nodes: nodes,
		Edges: edges,
	}, nil
}

// Amount 수량 줄이기
func shortenAmount(Amount string) string {
	parts := strings.Split(Amount, ".")

	if len(parts) < 2 || len(parts[1]) == 0 || len(parts[1]) <= 5 {
		return Amount
	}

	result := parts[1][0:4] + "..."

	return fmt.Sprintf("%s.%s", parts[0], result)
}

// Address 주소 줄이기
func shortenAddress(address string) string {
	if len(address) < 11 {
		return address
	} else {
		return fmt.Sprintf("%s...%s", address[0:7], address[len(address)-3:])
	}
}
