package transfer

import "context"

type Bank struct {
	BankCode     string `json:"bankCode"`
	BankFullName string `json:"bankFullName"`
	UrlLogo		 string	`json:"urlLogo"`
}

type RespBank struct {
	BankCode     string `json:"bankCode"`
	BankFullName string `json:"bankFullName"`
}

type BankList interface {
	GetBankList(context.Context) ([]Bank, error)
}
