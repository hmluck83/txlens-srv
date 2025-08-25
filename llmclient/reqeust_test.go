package llmclient

import (
	"context"
	"os"
	"testing"
)

func Test_classification(t *testing.T) {
	// Loadenv()

	l, err := NewLLMClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	inquiry, _ := os.ReadFile("test/zkBridgeTransaction.txt")

	result, err := l.Classifier(context.Background(), classifierPrompt, string(inquiry))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(*result)
}

func Test_addressPrompt(t *testing.T) {
	// Loadenv()
	ctx := context.Background()

	l, err := NewLLMClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	result, err := l.AddressPrompting(ctx, "0xf63D29B67AAbbaFa772C51e29DC7A89D391cFa7E")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(*result)
}
