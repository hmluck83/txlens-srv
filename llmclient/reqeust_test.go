package llmclient

import (
	"bytes"
	"context"
	"os"
	"testing"
	"text/template"

	"github.com/joho/godotenv"
)

func Loadenv() {
	_ = godotenv.Load("../.env")
}

func Test_LLMrequest(t *testing.T) {
	Loadenv()

	// Load Prompt Template
	tmpl, err := template.ParseFiles("prompts/templates/scaffold.tmpl")
	if err != nil {
		t.Fatal(err)
	}

	var bytebuf bytes.Buffer

	err = tmpl.Execute(&bytebuf, struct{}{})
	if err != nil {
		t.Fatal(err)
	}

	inquiry, err := os.ReadFile("test/zkBridgeTransaction.txt")

	// Run LLM
	l := NewLLMClient()
	result, err := l.Request(context.Background(), bytebuf.String(), string(inquiry))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*result)

}

func Test_LLMrequestWithTempalte(t *testing.T) {
	Loadenv()

	tmpl, err := template.ParseFiles("prompts/templates/scaffold.tmpl")
	if err != nil {
		t.Fatal(err)
	}

	var bytebuf bytes.Buffer

	err = tmpl.Execute(&bytebuf, struct{}{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(bytebuf.String())
}
