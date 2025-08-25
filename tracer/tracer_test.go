package tracer

import (
	"testing"

	"github.com/lmittmann/w3"
)

func TestTracer(t *testing.T) {

	txHash := w3.H("0xff8d3d66bd1c24130554a61796acccee4f21422ddafd26999138aa41606dba6f")

	fundFlows, addrLabels, err := FundFlowFromTx(txHash)
	if err != nil {
		panic(err)
	}

	for _, ff := range fundFlows {
		t.Log(ff)
	}

	for _, label := range addrLabels {
		t.Log(label)
	}
}
