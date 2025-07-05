package llmclient

import (
	"context"
	_ "embed"
	"os"

	"google.golang.org/genai"
)

//go:embed prompts/templates/scaffold.tmpl
var scaffoldTemplate string

func GetScaffoldTemplate() string {
	return scaffoldTemplate
}

type LLMClient struct {
	APIkey    string
	ModelName string
}

// LLM 요청을 위한 client 이제 보니 설계를 조금 잘못한 듯
func NewLLMClient() *LLMClient {
	return &LLMClient{
		APIkey:    os.Getenv("GEMINIAPI"),
		ModelName: os.Getenv("gemini-2.5-flash"),
	}
}

// LLM 요청 전달
func (l *LLMClient) Request(ctx context.Context, instruct string, inquiry string) (*string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  l.APIkey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	zeroPointer := int32(0)

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(instruct, genai.RoleUser),
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &zeroPointer,
		},
	}
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite-preview-06-17",
		genai.Text(inquiry),
		config,
	)

	if err != nil {
		return nil, err
	}

	text := result.Text()

	return &text, nil
}
