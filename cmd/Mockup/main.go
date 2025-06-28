package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type RequestPayload struct {
	Chainid int    `json:"chainID"`
	Txid    string `json:"txid"`
}

type GraphData struct {
	Nodes []NodeData `json:"nodes"`
	Edges []EdgeData `json:"edges"`
}

type ResponsePayload struct {
	Status      int       `json:"status"`
	Description string    `json:"description"`
	GraphData   GraphData `json:"graphdata"`
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

func main() {
	// Port Number
	const lensAddr string = ":7814"

	http.HandleFunc("/", proccessHandler)

	http.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling slow request...")
		time.Sleep(5 * time.Second) // 일부러 지연 발생
		fmt.Fprintf(w, "This was a slow request!")
		log.Println("Finished slow request.")
	})

	server := &http.Server{
		Addr: lensAddr,
	}

	// 3. 서버를 별도의 고루틴에서 시작
	go func() {
		log.Printf("Server starting on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt) // Ctrl+C (SIGINT) 시그널을 수신

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // context 리소스 해제

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited gracefully.")
}

func proccessHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w) // CORS 헤더 설정

	// Preflight 요청 (OPTIONS 메서드) 처리
	// 실제 요청 전에 브라우저가 CORS 정책을 확인하기 위해 보냄
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK) // 200 OK 응답
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only Post requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqPayload RequestPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqPayload); err != nil {
		http.Error(w, "Invalid JSON formet", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	a1Node := GraphNode{
		Id:   "0x559432...d53",
		Name: "0x559432...d53",
	}
	a1NodeData := NodeData{
		Data: a1Node,
	}

	a2Nodedata := NodeData{
		Data: GraphNode{
			Id:   "0xabB87A...AD8",
			Name: "0xabB87A...AD8",
		},
	}

	nodelist := []NodeData{
		a1NodeData, a2Nodedata,
	}

	edgelist := []EdgeData{
		{
			Data: GraphEdge{
				Source: "0x559432...d53",
				Target: "0xabB87A...AD8",
				Name:   "[1] 115,284.302665 Tether: USDT Stablecoin",
			},
		},
	}

	graphData := GraphData{
		Nodes: nodelist,
		Edges: edgelist,
	}

	responsePayload := ResponsePayload{
		Status:      1,
		Description: "",
		GraphData:   graphData,
	}

	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(responsePayload); err != nil {
		log.Printf("Error encoding response JSON: %s\n", err)
	}

}

func enableCORS(w *http.ResponseWriter) {
	// 모든 Origin을 허용: 개발 단계에서만 사용하고, 운영 환경에서는 특정 Origin으로 제한해야 합니다.
	(*w).Header().Set("Access-Control-Allow-Origin", "*")

	// 허용할 HTTP 메서드
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	// 허용할 헤더 (프론트엔드에서 커스텀 헤더를 보낼 경우 여기에 추가해야 함)
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// 자격 증명(쿠키, HTTP 인증 등)을 포함한 요청 허용 여부
	// (*w).Header().Set("Access-Control-Allow-Credentials", "true") // Access-Control-Allow-Origin이 '*'일 때는 사용할 수 없음
}
