package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hmluck83/txlens-srv/llmclient"
	"github.com/hmluck83/txlens-srv/tracer"
	"github.com/shopspring/decimal"
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

	fundFlows, addrLabels, err := tracer.FundFlowFromTx(common.HexToHash(reqPayload.Txid))

	if err != nil {
		log.Printf("Error in Transaction summary: %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// NOTE: client 생성 위치? 매번 HTTP Requests?
	// 아직 접속이 많이 없어 고려하지 않도록 한다
	lc, err := llmclient.NewLLMClient(context.Background())
	if err != nil {
		log.Printf("Error on Create LLM Client %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 각 주소의 Asistant Prompt 생성
	var wg sync.WaitGroup
	var mu sync.Mutex
	var sb strings.Builder

	for keyAddr := range addrLabels {
		if keyAddr.Cmp(tracer.EthAddress) != 0 {
			wg.Add(1)

			go func() {
				defer wg.Done()
				// 현재 addr Prompt 생성중 에러는 무시
				addrPrompt, _ := lc.AddressPrompting(context.Background(), keyAddr.Hex())

				mu.Lock()
				sb.WriteString(*addrPrompt)
				mu.Unlock()
			}()
		}
	}

	txFlowSummary := flowSummary(&fundFlows, &addrLabels)
	summarized, err := json.Marshal(txFlowSummary)

	txType, err := lc.Classifier(context.Background(), sb.String(), string(summarized))
	if err != nil {
		log.Printf("Error on generate summary %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	prompt := lc.GetSummaryPrompt(*txType)

	llmresponse, err := lc.Summary(context.Background(), prompt+sb.String(), string(summarized))

	if err != nil {
		log.Printf("Error on generate summary %s\n", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	graphData, _ := buildGraphData(fundFlows, addrLabels)
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
func buildGraphData(fundFlows tracer.FundFlows, addrLabels tracer.AddrLabels) (*GraphData, error) {

	// node를 중복 없이 저장, addrLabel에는 Token Address의 정보도 함께 포함되어 있음
	// Graph Node의 정보만 nodes에 추가
	nodeMap := make(map[common.Address]struct{})
	edges := []EdgeData{}

	for idx, fundFlow := range fundFlows {
		nodeMap[fundFlow.From] = struct{}{}
		nodeMap[fundFlow.To] = struct{}{}

		// edge에서 거래량과 Token 종류 표시
		name := fmt.Sprintf("[%d] %s %s",
			idx,
			shortenAmount(decimal.NewFromBigInt(&fundFlow.Value, -(int32(addrLabels[fundFlow.Token].Decimals))).String()),
			addrLabels[fundFlow.Token].Symbol,
		)

		edges = append(edges, EdgeData{
			GraphEdge{
				Source: fundFlow.From.Hex(),
				Target: fundFlow.To.Hex(),
				Name:   name,
			}})
	}

	nodes := []NodeData{}

	for key := range nodeMap {
		var name string
		if label, exist := addrLabels[key]; exist {
			name = fmt.Sprintf("%s(%s)", label.Label, shortenAddress(key.Hex()))
		} else {
			name = shortenAddress(key.Hex())
		}

		nodes = append(nodes, NodeData{
			Data: GraphNode{
				Id:   key.Hex(),
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

// FundFlows와 Address Label을 LLM에 넘겨줄 Summary 형태로 정리
func flowSummary(fundFlows *tracer.FundFlows, addrLabels *tracer.AddrLabels) *summaryDesc {
	// !!!Under construction

	summaryProfile := summaryProfile{
		Sender:        "",
		Status:        true,
		RevertMessage: "",
	}

	fundDescs := []FundDesc{}

	for idx, fundflow := range *fundFlows {
		var fromLabel string
		var toLabel string
		var token string
		var tokenDecimal uint8

		if v, exist := (*addrLabels)[fundflow.From]; exist {
			fromLabel = v.Label
		} else {
			fromLabel = shortenAddress(fundflow.From.Hex())
		}

		if v, exsit := (*addrLabels)[fundflow.To]; exsit {
			toLabel = v.Label
		} else {
			toLabel = shortenAddress(fundflow.From.Hex())
		}

		if v, exist := (*addrLabels)[fundflow.Token]; exist {
			token = v.Symbol
			tokenDecimal = v.Decimals
		}

		fundDesc := FundDesc{
			Amount:     decimal.NewFromBigInt(&fundflow.Value, -(int32(tokenDecimal))).String(),
			From:       fundflow.From.Hex(),
			FromLabel:  fromLabel,
			ID:         idx,
			IsERC1155:  false,
			IsERC721:   false,
			IsReverted: false,
			Order:      idx,
			To:         fundflow.To.Hex(),
			ToLabel:    toLabel,
			Token:      token,
		}
		fundDescs = append(fundDescs, fundDesc)

	}

	return &summaryDesc{
		SummaryProfile: summaryProfile,
		FundDesc:       fundDescs,
	}
}
