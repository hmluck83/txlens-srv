package llmclient

import (
	"context"
	_ "embed"
	"os"

	"google.golang.org/genai"
)

//go:embed prompts/templates/scaffold.tmpl
var scaffoldTemplate string

//go:embed prompts/templates/classification.txt
var classifierTemplate string

var ClassificationEnum []string

func init() {
	ClassificationEnum = []string{
		"Withdraw_From_CentralizedExchange",
		"Deposit_To_CentralizedExchange",
		"Simple_Transfer",
		"Swap_Transaction",
	}
}

func GetScaffoldTemplate() string {
	return scaffoldTemplate
}

type LLMClient struct {
	APIkey       string
	ModelName    string
	geminiClient *genai.Client
}

// LLM 요청을 위한 client 이제 보니 설계를 조금 잘못한 듯
func NewLLMClient(ctx context.Context) (*LLMClient, error) {
	apiKey := os.Getenv("GEMINIAPI")
	modelName := os.Getenv("gemini-2.5-flash")

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		return nil, err
	}

	return &LLMClient{
		APIkey:       apiKey,
		ModelName:    modelName,
		geminiClient: client,
	}, nil
}

// Build Summary Request
func (l *LLMClient) Summary(ctx context.Context, instruct string, inquiry string) (*string, error) {

	zeroPointer := int32(0)

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(instruct, genai.RoleUser),
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &zeroPointer,
		},
	}

	result, err := l.geminiClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite-preview-06-17", // ch Model name
		genai.Text(inquiry),
		config,
	)

	if err != nil {
		return nil, err
	}

	text := result.Text()

	return &text, nil
}

// TODO !TODO
func (l *LLMClient) Classifier(ctx context.Context, instruct string, inquiry string, schemas []string) (*string, error) {

	zeroPointer := int32(0)

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(instruct, genai.RoleUser),
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &zeroPointer,
		},
		ResponseMIMEType: "text/x.enum",
		ResponseSchema: &genai.Schema{
			Type:   "STRING",
			Format: "enum",
			Enum:   schemas,
		},
		// Google Search Option의 반환 값은 항상 동일한듯?
		// Tools: []*genai.Tool{
		// 	{
		// 		GoogleSearch: &genai.GoogleSearch{},
		// 	},
		// },
	}

	result, err := l.geminiClient.Models.GenerateContent(
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
