package llmclient

import (
	"bytes"
	"context"
	"embed"
	_ "embed"
	"os"

	"google.golang.org/genai"
)

//go:embed prompts/templates
var templates embed.FS

var summarySwapPrompt string
var summaryDepositCex string
var summaryWithdrawCex string
var summarySimpleTrans string
var classifierPrompt string

var ClassificationEnum []string

func init() {
	ClassificationEnum = []string{
		"Withdraw_From_CentralizedExchange",
		"Deposit_To_CentralizedExchange",
		"Simple_Transfer",
		"Swap_Transaction",
	}

	// build template strings
	personaSummary, _ := templates.ReadFile("prompts/templates/persona_summary.tmpl")
	personaClassification, _ := templates.ReadFile("prompts/templates/persona_classification.tmpl")

	instructionSummary, _ := templates.ReadFile("prompts/templates/instruction_summary.tmpl")
	instructionClassification, _ := templates.ReadFile("prompts/templates/instruction_classification.tmpl")

	transactionTemplate, _ := templates.ReadFile("prompts/templates/transaction_template.tmpl")

	infoSwap, _ := templates.ReadFile("prompts/templates/info_swap.tmpl")
	infoDepositCex, _ := templates.ReadFile("prompts/templates/info_depositCEX.tmpl")
	infoWithdrawCex, _ := templates.ReadFile("prompts/templates/info_withdrawCEX.tmpl")

	// TODO: 분명히 개선될거 같은데
	cp := bytes.Buffer{}
	cp.Write(personaClassification)
	cp.Write(transactionTemplate)
	cp.Write(instructionClassification)

	classifierPrompt = cp.String()

	swapPrompt := bytes.Buffer{}
	swapPrompt.Write(personaSummary)
	swapPrompt.Write(transactionTemplate)
	swapPrompt.Write(instructionSummary)
	swapPrompt.Write(infoSwap)

	summarySwapPrompt = swapPrompt.String()

	depositCex := bytes.Buffer{}
	depositCex.Write(personaSummary)
	depositCex.Write(transactionTemplate)
	depositCex.Write(instructionSummary)
	depositCex.Write(infoDepositCex)

	summaryDepositCex = depositCex.String()

	withdrawCex := bytes.Buffer{}
	withdrawCex.Write(personaSummary)
	withdrawCex.Write(transactionTemplate)
	withdrawCex.Write(instructionSummary)
	withdrawCex.Write(infoWithdrawCex)

	summaryWithdrawCex = withdrawCex.String()

	simpleTrans := bytes.Buffer{}
	simpleTrans.Write(personaSummary)
	simpleTrans.Write(transactionTemplate)
	simpleTrans.Write(instructionSummary)

	summarySimpleTrans = simpleTrans.String()
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
func (l *LLMClient) Classifier(ctx context.Context, inquiry string) (*string, error) {

	zeroPointer := int32(0)

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(classifierPrompt, genai.RoleUser),
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &zeroPointer,
		},
		ResponseMIMEType: "text/x.enum",
		ResponseSchema: &genai.Schema{
			Type:   "STRING",
			Format: "enum",
			Enum:   ClassificationEnum,
		},
		// Google Search Option하고 Structured Output은 동시에 쓸 수 없는 듯 한번 더 갔다오???
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
	}

	return "" // Never

}
