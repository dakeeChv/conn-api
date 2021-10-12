package transfer

import (
	jdbsdk "bank-apis/banksdk"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type bankListCore struct {
	sdk *jdbsdk.SDK
}

func NewBankListCore(sdk *jdbsdk.SDK) BankList {
	return bankListCore{sdk: sdk}
}

func (s bankListCore) GetBankList(ctx context.Context) ([]Bank, error) {
	body, _ := json.Marshal(struct {
		RequestID string `json:"requestId"`
	}{
		RequestID: strconv.FormatInt(time.Now().UnixNano(), 10),
	})
	req := s.sdk.NewHTTPReq(ctx, http.MethodPost, s.sdk.Cfg.BaseURL+"/txn/getInfo/bankList", body)
	res, err := s.sdk.Hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("listAccountByCIF: http.Do: %v", err)
	}
	defer res.Body.Close()

	info, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("listAccountByCIF: read response body: %v", err)
	}

	var reply struct {
		Success  bool   `json:"success"`
		Message  string `json:"message"`
		BankList []Bank `json:"data"`
	}

	if err := json.Unmarshal(info, &reply); err != nil {
		return nil, fmt.Errorf("listAccountByCIF: unmarshal json: %s", info)
	}
	if !reply.Success || len(reply.BankList) == 0 {
		return nil, fmt.Errorf("listAccountByCIF: %s", reply.Message)
	}

	bankList := reply.BankList
	fmt.Println(reply.Success)
	fmt.Println(reply.Message)

	return bankList, nil
}
