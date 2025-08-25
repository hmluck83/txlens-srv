package llmclient

import (
	"bytes"
	"embed"
	"strings"
)

//go:embed prompts/templates
var templates embed.FS

var summarySwapPrompt string
var summaryDepositCex string
var summaryWithdrawCex string
var summarySimpleTrans string
var summaryDepositBridge string
var summaryWithdrawBridge string
var classifierPrompt string
var addressPrompt string

var ClassificationEnum []string

func init() {
	ClassificationEnum = []string{
		"Withdraw_From_CentralizedExchange",
		"Deposit_To_CentralizedExchange",
		"Simple_Transfer",
		"Swap_Transaction",
		"Withdraw_From_Bridge",
		"Deposit_To_Bridge",
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
	infoDepsitBridge, _ := templates.ReadFile("prompts/templates/info_depositBridge.tmpl")
	infoWithdrawBridge, _ := templates.ReadFile("prompts/templates/info_withdrawBridge.tmpl")

	at, _ := templates.ReadFile("prompts/templates/address_prompt.tmpl")
	addressPrompt = string(at)

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

	depositBridge := bytes.Buffer{}
	depositBridge.Write(personaSummary)
	depositBridge.Write(transactionTemplate)
	depositBridge.Write(instructionSummary)
	depositBridge.Write(infoDepsitBridge)

	summaryDepositBridge = depositBridge.String()

	withdrawBridge := bytes.Buffer{}
	withdrawBridge.Write(personaSummary)
	withdrawBridge.Write(transactionTemplate)
	withdrawBridge.Write(instructionSummary)
	withdrawBridge.Write(infoWithdrawBridge)

	summaryWithdrawBridge = withdrawBridge.String()

	simpleTrans := bytes.Buffer{}
	simpleTrans.Write(personaSummary)
	simpleTrans.Write(transactionTemplate)
	simpleTrans.Write(instructionSummary)

	summarySimpleTrans = simpleTrans.String()
}

func buildClassifyPrompt(adddressPrompt string) string {
	var sb strings.Builder

	sb.WriteString(classifierPrompt)
	sb.WriteString("\n # 주소 안내")
	sb.WriteString(adddressPrompt)

	return sb.String()
}
