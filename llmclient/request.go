package llmclient

import (
	"context"
	_ "embed"
	"os"

	"google.golang.org/genai"
)

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
func (l *LLMClient) Summary(ctx context.Context, prompt string, inquiry string) (*string, error) {

	zeroPointer := int32(0)

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(prompt, genai.RoleUser),
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
func (l *LLMClient) Classifier(ctx context.Context, addrPrompt string, inquiry string) (*string, error) {

	zeroPointer := int32(0)

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(classifierPrompt+addrPrompt, genai.RoleUser),
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &zeroPointer,
		},
		ResponseMIMEType: "text/x.enum",
		ResponseSchema: &genai.Schema{
			Type:   "STRING",
			Format: "enum",
			Enum:   ClassificationEnum,
		},
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

// TODO !TODO
func (l *LLMClient) AddressPrompting(ctx context.Context, inquiry string) (*string, error) {

	// Config Thinking budget
	thinkingBudget := int32(0)
	temparature := float32(0.3)

	config := &genai.GenerateContentConfig{
		Temperature:       &temparature,
		SystemInstruction: genai.NewContentFromText(addressPrompt, genai.RoleUser),
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &thinkingBudget,
		},
		Tools: []*genai.Tool{
			{
				GoogleSearch: &genai.GoogleSearch{},
			},
		},
	}

	result, err := l.geminiClient.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(inquiry),
		config,
	)

	if err != nil {
		return nil, err
	}

	text := result.Text()
	return &text, nil
}

func (l *LLMClient) GetSummaryPrompt(txtype string) string {
	switch txtype {
	case "Withdraw_From_CentralizedExchange":
		return summaryWithdrawCex
	case "Deposit_To_CentralizedExchange":
		return summaryDepositCex
	case "Simple_Transfer":
		return summarySimpleTrans
	case "Swap_Transaction":
		return summarySwapPrompt
	case "Withdraw_From_Bridge":
		return summaryWithdrawBridge
	case "Deposit_To_Bridge":
		return summaryDepositBridge
	}

	return "" // Never

}
