package main

import (
	"fmt"
	"os"

	"github.com/hmluck83/txlens-srv/tracer"
	"github.com/joho/godotenv"
	"github.com/lmittmann/w3"
)

func main() {
	_ = godotenv.Load(".env")

	if len(os.Args) != 2 {
		panic("TXID is required")
	}

	txHash := w3.H(os.Args[1])

	fundFlows, _, err := tracer.FundFlowFromTx(txHash)
	if err != nil {
		panic(err)
	}

	for _, ff := range fundFlows {
		fmt.Printf("From: %s To: %s | Value %s | Token %s\n", ff.From.Hex(), ff.To.Hex(), ff.Value.String(), ff.Token.Hex())
	}

}
