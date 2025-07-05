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
	l := NewLLMClient()
	result, err := l.Request(context.Background(), tmpl, string(inquiry))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*result)

}
