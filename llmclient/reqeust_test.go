package llmclient

import (
	"context"
	"os"
	"testing"

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

	result, err := l.Request(context.Background(), tmpl, string(inquiry))
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
