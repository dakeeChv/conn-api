package transfer

import "context"

//requset info
type ReqDest struct {
	RequestID       string `json:"requestId"`
	Type            string `json:"txnType"`
	SourceCustNo    string `json:"sourceCustNo"`
	SourceAccountNo string `json:"sourceAccountNo"`
	SourceCurrency  string `json:"sourceCurrency"`
	DestBankCode    string `json:"destBankCode"`
	DestAccountNo   string `json:"destAccountNo"`
	ExReferenceNo   string `json:"exReferenceNo"`
}

//track data
type InitResp struct {
	DestAccountInfo []struct {
		AccountNo   string `json:"accountNo"`
		AccountName string `json:"accountName"`
		AccountType string `json:"accountType"`
		Currency    string `json:"currency"`
	}
	SourceAccountInfo []struct {
		ExReferenceNo string  `json:"exReferenceNo"`
		SourceCustNo  string  `json:"sourceCustNo"`
		LimitPerDay   float64 `json:"txnLimitPerDay"`
		ExchangeRate  []struct {
			SourceCurrency string `json:"sourceCurrency"`
			DestCurrency   string `json:"destCurrency"`
			Rate           int    `json:"rate"`
		}
		FeeList []struct {
			Currency string `json:"currency"`
			List     []struct {
				From      int64     `json:"from"`
				FeeAmount float64 `json:"feeAmount"`
			}
		}
	}
}

type DestAcc struct {
	AccountNo   string `json:"accountNo"`
	DisplayName string `json:"accountName"`
	AccountType string `json:"accountType"`
	Currency    string `json:"currency"`
}

type SourceAcc struct {
	ExReferenceNo string  `json:"exReferenceNo"`
	SourceCustNo  string  `json:"sourceCustNo"`
	LimiterPerDay float64 `json:"txnLimitPerDay"`
}

type ExchRate struct {
	SourceCurrency string `json:"sourceCurrency"`
	DestCurrency   string `json:"destCurrency"`
	Rate           int    `json:"rate"`
}

type Fee struct {
	From      int64 `json:"from"`
	FeeAmount float64 `json:"feeAmount"`
}

type RespDestAcc struct {
	DestAccount   []DestAcc   `json:"destAccount"`
	SourceAccount []SourceAcc `json:"sourceAccount"`
	ExchangeRate  []ExchRate  `json:"exchangeRate"`
}

type CallDest interface {
	GetDestInfo(context.Context, ReqDest) (*RespDestAcc, error)
}
