package fetcher

import (
	"fmt"
)

type summaryProfile struct {
	Sender        string `json:"sender"`
	Status        bool   `json:"status"`
	RevertMessage string `json:"revertMessage"`
}

// 향후 Summary Desc로 대체 됨
type summary struct {
	SummaryProfile summaryProfile `json:"summaryProfile"`
	FundFlows      []FundFlow     `json:"fundFlow"`
	// AddressLabel   []AddressLabel `json:"addressLabel"`
}

type summaryDesc struct {
	SummaryProfile summaryProfile `json:"summaryProfile"`
	FundDesc       []FundDesc     `json:"fundDesc"`
}

// blocksec으로 부터 Profile과 Address Label을 질의 형식에 맞게 정규화
func Summarizer(txprofile Profile, txlabel AddressLabels) summary {

	labelMap := make(map[string]string)

	for _, val := range txlabel.Labels {
		if (val.Label != "Null Address") && (val.Label != "Precompiled") {
			labelMap[val.Address] = val.Label
		}
	}

	summaryProfileObj := summaryProfile{
		Sender:        txprofile.BasicInfo.Sender,
		Status:        txprofile.BasicInfo.Status,
		RevertMessage: txprofile.BasicInfo.RevertMessage,
	}

	for idx := range txprofile.FundFlows {
		txprofile.FundFlows[idx].Token = labelMap[txprofile.FundFlows[idx].Token]

		if val, exist := labelMap[txprofile.FundFlows[idx].From]; exist {
			txprofile.FundFlows[idx].From = fmt.Sprintf("%s(%s)", val, shortenAddress(txprofile.FundFlows[idx].From))
		} else {
			txprofile.FundFlows[idx].From = shortenAddress(txprofile.FundFlows[idx].From)
		}

		if val, exist := labelMap[txprofile.FundFlows[idx].To]; exist {
			txprofile.FundFlows[idx].To = fmt.Sprintf("%s(%s)", val, shortenAddress(txprofile.FundFlows[idx].To))
		} else {
			txprofile.FundFlows[idx].To = shortenAddress(txprofile.FundFlows[idx].To)
		}
	}

	summaryObj := summary{
		SummaryProfile: summaryProfileObj,
		FundFlows:      txprofile.FundFlows,
	}

	return summaryObj
}

/*
blocksec으로 부터 Profile과 Address Label을 질의 형식에 맞게 정규화
Classification 용도에 맞도록 수정
*/
func SummarizerClassification(txprofile Profile, txlabel AddressLabels) summaryDesc {

	labelMap := make(map[string]string)

	for _, val := range txlabel.Labels {
		if (val.Label != "Null Address") && (val.Label != "Precompiled") {
			labelMap[val.Address] = val.Label
		}
	}

	summaryProfileObj := summaryProfile{
		Sender:        txprofile.BasicInfo.Sender,
		Status:        txprofile.BasicInfo.Status,
		RevertMessage: txprofile.BasicInfo.RevertMessage,
	}

	// 변수 네이밍  일관성 없네 진짜
	fundDecsSlice := []FundDesc{}

	for _, ff := range txprofile.FundFlows {
		fundDecs := FundDesc{
			Amount:    ff.Amount,
			From:      ff.From,
			ID:        ff.ID,
			IsERC1155: ff.IsERC1155,
			IsERC721:  ff.IsERC721,
			Order:     ff.Order,
			To:        ff.To,
			Token:     labelMap[ff.Token],
		}

		if val, exist := labelMap[ff.From]; exist {
			fundDecs.FromLabel = fmt.Sprintf("%s(%s)", val, shortenAddress(ff.From))
		} else {
			fundDecs.FromLabel = shortenAddress(ff.From)
		}

		if val, exist := labelMap[ff.To]; exist {
			fundDecs.ToLabel = fmt.Sprintf("%s(%s)", val, shortenAddress(ff.To))
		} else {
			fundDecs.ToLabel = shortenAddress(ff.To)
		}

		fundDecsSlice = append(fundDecsSlice, fundDecs)
	}

	summaryObj := summaryDesc{
		SummaryProfile: summaryProfileObj,
		FundDesc:       fundDecsSlice,
	}

	return summaryObj
}

// TODO: Util package로 정리해서 api/summary.go와 정리

// Address 주소 줄이기
func shortenAddress(address string) string {
	if len(address) < 11 {
		return address
	} else {
		return fmt.Sprintf("%s...%s", address[0:7], address[len(address)-3:])
	}
}
