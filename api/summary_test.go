package api

import (
	"testing"

	"github.com/hmluck83/txlens-srv/tracer"
	"github.com/lmittmann/w3"
)

func TestSummary(t *testing.T) {

	txHash := w3.H("0xff8d3d66bd1c24130554a61796acccee4f21422ddafd26999138aa41606dba6f")

	fundFlows, addrLabels, err := tracer.FundFlowFromTx(txHash)

	if err != nil {
		panic(err)
	}

	for key, val := range addrLabels {
		t.Logf("Address is %s\n", key.Hex())
		t.Logf("Label: %s, Name: %s, Symbol: %s", val.Label, val.Name, val.Symbol)
	}
	suumaryDesc := flowSummary(&fundFlows, &addrLabels)

	t.Log(suumaryDesc)
}
