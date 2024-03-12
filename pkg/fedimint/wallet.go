package fedimint

import "encoding/json"

func (wallet *WalletModule) createDepositAddress(request DepositAddressRequest, federationId *string) (*DepositAddressResponse, error) {
	resp, err := wallet.Client.post("/wallet/deposit-address", request)
	if err != nil {
		return nil, err
	}
	var depositAddressResp DepositAddressResponse
	err = json.Unmarshal(resp, &depositAddressResp)
	if err != nil {
		return nil, err
	}
	return &depositAddressResp, nil
}

func (wallet *WalletModule) awaitDeposit(request AwaitDepositRequest, federationId *string) (*AwaitDepositResponse, error) {
	resp, err := wallet.Client.post("/wallet/await-deposit", request)
	if err != nil {
		return nil, err
	}
	var depositResp AwaitDepositResponse
	err = json.Unmarshal(resp, &depositResp)
	if err != nil {
		return nil, err
	}
	return &depositResp, nil
}

func (wallet *WalletModule) withdraw(request WithdrawRequest, federationId *string) (*WithdrawResponse, error) {
	resp, err := wallet.Client.post("/wallet/withdraw", request)
	if err != nil {
		return nil, err
	}
	var withdrawResp WithdrawResponse
	err = json.Unmarshal(resp, &withdrawResp)
	if err != nil {
		return nil, err
	}
	return &withdrawResp, nil
}

type DepositAddressRequest struct {
	Timeout int `json:"timeout"`
}

type DepositAddressResponse struct {
	OperationID string `json:"operation_id"`
	Address     string `json:"address"`
}

type AwaitDepositRequest struct {
	OperationID string `json:"operation_id"`
}

type AwaitDepositResponse struct {
	Status string `json:"status"`
}

type WithdrawRequest struct {
	Address    string `json:"address"`
	AmountMsat string `json:"amount_msat"`
}

type WithdrawResponse struct {
	Txid    string `json:"txid"`
	FeesSat int    `json:"fees_sat"`
}
