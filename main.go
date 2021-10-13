package main

import (
	jdbsdk "bank-apis/banksdk"
	"bank-apis/core/transfer"
	"context"
	"fmt"
	"time"
)

func main() {
	cfg := jdbsdk.Config{
		BaseURL:   "https://services.jdbbank.com.la:12009/mbapi-uat",
		UserID:    "jdbmbapiuat",
		SecretKey: "secdQM6CLQD4SBLM1W9I7XYZ5OE9DRS9DB9KNSR9G3RNNH57ZXOZYXCX6D7I5EOFIRCXW3BPRAOOY1S1UH0AY87SW8TKO03UP6PWWAYGt",
		HMacKey:   []byte("M0UTLC2QFU9HX4JN9LXFCLMNA49F11CW075ONJVJP1HCLH3AQZ25A775EQSSXYK9YLCLJDM7WSLTDNRP651YDILFDJQAL58NGXFO"),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	sdk, err := jdbsdk.New(ctx, &cfg)
	if err != nil {
		panic(err)
	}

	bankList := transfer.NewBankListCore(sdk)
	// destAcc := transfer.NewdestAccCore(sdk)

	// reqBody := transfer.ReqDest{
	// 	RequestID: "20210922120212456",
	// 	Type: "P2P",           
	// 	SourceCustNo: "001063976",    
	// 	SourceAccountNo: "00120010110008245", 
	// 	SourceCurrency: "USD",
	// 	DestBankCode: "JDB",   
	// 	DestAccountNo: "001082973", 
	// 	ExReferenceNo: "W92JG1BG4I1J1TH3S2",
	// }
	data, err := bankList.GetBankList(ctx)
	// data, err := destAcc.GetDestInfo(ctx, reqBody)
	if err != nil {
		println(err)
	}

	fmt.Printf("%+v", data)
}
