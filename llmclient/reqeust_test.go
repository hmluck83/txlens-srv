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

	// Load Prompt Template
	tmpl := GetScaffoldTemplate()

	inquiry, err := os.ReadFile("test/zkBridgeTransaction.txt")

	// Run LLM
	l, err := NewLLMClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	result, err := l.Summary(context.Background(), tmpl, string(inquiry))
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

	result, err := l.Classifier(context.Background(), classifierTemplate, string(inquiry), ClassificationEnum)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(*result)
}

func Test_classificationWithFetch(t *testing.T) {
	Loadenv()

	profile, addresslabes, err := fetcher.FetchTransaction("0x348b46ca8d967ce1d1e74e866911c47a99796cec1e97c92cd2a95faea15a0745")
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

	result, err := l.Classifier(context.Background(), classifierTemplate, string(summanryString), ClassificationEnum)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(*result)
}
