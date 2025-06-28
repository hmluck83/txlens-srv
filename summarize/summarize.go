package summarize

import "github.com/hmluck83/txlens-srv/schemas"

type summuryProfile struct {
	Sender        string `json:"sender"`
	Status        bool   `json:"status"`
	RevertMessage string `json:"revertMessage"`
}

type summury struct {
	SummuryProfile summuryProfile         `json:"summuryProfile"`
	FundFlows      []schemas.FundFlow     `json:"fundFlow"`
	AddressLabel   []schemas.AddressLabel `json:"addressLabel"`
}

func Summarizer(txprofile schemas.Profile, txlabel schemas.AddressLabels) summury {

	summuryProfileObj := summuryProfile{
		Sender:        txprofile.BasicInfo.Sender,
		Status:        txprofile.BasicInfo.Status,
		RevertMessage: txprofile.BasicInfo.RevertMessage,
	}

	var addresslabels []schemas.AddressLabel

	for _, val := range txlabel.Labels {
		if (val.Label != "Null Address") && (val.Label != "Precompiled") {
			addresslabels = append(addresslabels, val)
		}
	}
	summuryObj := summury{
		SummuryProfile: summuryProfileObj,
		FundFlows:      txprofile.FundFlows,
		AddressLabel:   addresslabels,
	}

	return summuryObj
}
