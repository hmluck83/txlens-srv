package llmclient

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hmluck83/txlens-srv/fetcher"
	"github.com/joho/godotenv"
)

func Loadenv() {
	_ = godotenv.Load("../.env")
}

func Test_LLMrequest(t *testing.T) {
	Loadenv()

	inquiry, err := os.ReadFile("test/zkBridgeTransaction.txt")

	// Run LLM
	l, err := NewLLMClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	prompt := l.GetSummaryPrompt("Deposit_To_CentralizedExchange")

	result, err := l.Summary(context.Background(), prompt, string(inquiry))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*result)
}

func Test_classification(t *testing.T) {
	Loadenv()

	l, err := NewLLMClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	inquiry, _ := os.ReadFile("test/zkBridgeTransaction.txt")

	result, err := l.Classifier(context.Background(), string(inquiry))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(*result)
}

func Test_classificationWithFetch(t *testing.T) {
	Loadenv()

	profile, addresslabes, err := fetcher.FetchTransaction("0x04f9941cf871f8e11e49f53b8c276a9c930db7c0dfd8a42c08687b333c73da49")
	if err != nil {
		t.Fatal(err)
	}

	summaryObj := fetcher.SummarizerClassification(*profile, *addresslabes)
	summanryString, err := json.Marshal(summaryObj)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(summanryString))
	l, err := NewLLMClient(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	result, err := l.Classifier(context.Background(), string(summanryString))

	if err != nil {
		t.Fatal(err)
	}

	t.Log(*result)
}

func Test_addressPrompt(t *testing.T) {
	Loadenv()
	ctx := context.Background()

	l, err := NewLLMClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	result, err := l.AddressPrompting(ctx, "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(*result)
}
