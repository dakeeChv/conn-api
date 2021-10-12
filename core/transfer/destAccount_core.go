package transfer

import (
	jdbsdk "bank-apis/banksdk"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type destAccountCore struct {
	sdk *jdbsdk.SDK
}

func NewdestAccCore(sdk *jdbsdk.SDK) CallDest {
	return destAccountCore{sdk: sdk}
}

func (d destAccountCore) GetDestInfo(ctx context.Context, r ReqDest) (*RespDestAcc, error) {
	body, _ := json.Marshal(r)
	req := d.sdk.NewHTTPReq(ctx, http.MethodPost, d.sdk.Cfg.BaseURL+"/txn/getInfo/destAccount", body)
	res, err := d.sdk.Hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}
	defer res.Body.Close()
	info, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var reply struct {
		Success bool          `json:"success"`
		Message string        `json:"message"`
		Data    []InitResp  `json:"data"`
	}
	err = json.Unmarshal(info, &reply)
	if err != nil {
		panic(err)
	}
	data := reply.Data[0]

	listDestAcc := data.DestAccountInfo
	destAccs := []DestAcc{}
	for _, v := range listDestAcc{
		destAcc := DestAcc{
			AccountNo: v.AccountNo,
			DisplayName: v.AccountName,
			AccountType: v.AccountType,
			Currency: v.Currency,
		}
		destAccs = append(destAccs, destAcc)
	}
	// fmt.Printf("%+v\n", destAccs)

	// sourceAcc := SourceAcc{
	// 	ExReferenceNo: data.SourceAccountInfo[0].ExReferenceNo,
	// 	SourceCustNo: data.SourceAccountInfo[0].SourceCustNo,
	// 	LimiterPerDay: data.SourceAccountInfo[0].LimitPerDay,
	// }
	// _ = sourceAcc
	listSource := data.SourceAccountInfo
	sourceAccs := []SourceAcc{}
	for _, v := range listSource{
		sourceAcc := SourceAcc{
			ExReferenceNo: v.ExReferenceNo,
			SourceCustNo: v.SourceCustNo,
			LimiterPerDay: v.LimitPerDay,
		}
		sourceAccs = append(sourceAccs, sourceAcc)
	}
	// fmt.Printf("%+v\n", sourceAccs)

	listExchangeRate := data.SourceAccountInfo[0].ExchangeRate
	exchangeRates := []ExchRate{}
	for _, v := range listExchangeRate{
		exchangeRate := ExchRate{
			SourceCurrency: v.SourceCurrency,
			DestCurrency: v.DestCurrency,
			Rate: v.Rate,
		}
		exchangeRates = append(exchangeRates, exchangeRate)
	}
	// fmt.Printf("%+v\n", exchangeRates)

	listFee := data.SourceAccountInfo[0].FeeList[0].List
	fees := []Fee{}
	for _, v := range listFee{
		fee := Fee{
			From: v.From,
			FeeAmount: v.FeeAmount,
		}
		fees = append(fees, fee)
	}
	// fmt.Printf("%+v", fees)

	response := RespDestAcc{
		DestAccount: destAccs,
		SourceAccount: sourceAccs,
		ExchangeRate: exchangeRates,
	}

	return &response, nil
}
