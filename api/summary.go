package api

import (
	"bytes"
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

// Transaction Summary를 작성 API
func summuryHandler(w http.ResponseWriter, r *http.Request) {

	/* TODO List
	[ ] Verifing chainID
	[ ] Verifing Ethereum Transaction
	*/

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

	// NOTE: Better Logger
	log.Printf("New access from %s. Start summary generate\n", r.RemoteAddr)
	log.Printf("Requested TXid : %s   ChainID: %d\n", reqPayload.Txid, reqPayload.Chainid)

	// Fetch Transaction Infomation
	profile, addressLabels, err := fetcher.FetchTransaction(reqPayload.Txid)
	if err != nil {
		log.Printf("Error in Transaction summary: %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	summarizedObj := fetcher.Summarizer(*profile, *addressLabels)

	summarizedString, err := json.Marshal(summarizedObj)
	if err != nil {
		log.Printf("On summuring transaction has error")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var intendJSON bytes.Buffer
	json.Indent(&intendJSON, summarizedString, "", "  ")

	log.Printf("Transaction Summary is :\n%s\n", intendJSON.String())

	// NOTE: client 생성 위치? 매번 HTTP Requests?
	// 아직 접속이 많이 없어 고려하지 않도록 한다
	lc, err := llmclient.NewLLMClient(context.Background())
	if err != nil {
		log.Printf("Error on Create LLM Client %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println(string(summarizedString))

	llmresponse, err := lc.Summary(context.Background(), string(summarizedString))

	if err != nil {
		log.Printf("Error on generate summary %s\n", err.Error())
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

	log.Printf("Summary String is : \n%s\n", responseSummury.Description)
}

// Transaction FundFlow와 Address Label을 기준으로 Cytoscape Graph Data 생성
func buildGraphData(profile *fetcher.Profile, addressLabels *fetcher.AddressLabels) (*GraphData, error) {

	// Address Labels to map
	labelMap := make(map[string]string)

	for _, val := range addressLabels.Labels {
		// 기존에 키가있는지 검사 하는것보다 무지성으로 때려 박는게 더 빠르지 않을까?
		labelMap[val.Address] = val.Label
	}

	// node를 중복 없이 저장할 필요가 있지만 Set이 없기 때문에 Map으로 대체
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

// Amount 수량 문자열 정규화
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
