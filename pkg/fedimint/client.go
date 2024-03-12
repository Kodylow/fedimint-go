package fedimint

import (
	"bytes"
	"encoding/json"
	"fedimint-go-client/pkg/fedimint/types"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FedimintClient struct {
	BaseURL      string
	Password     string
	FederationId string
	Ln           LnModule
	Wallet       WalletModule
	Mint         MintModule
}

type LnModule struct {
	Client *FedimintClient
}

type MintModule struct {
	Client *FedimintClient
}

type WalletModule struct {
	Client *FedimintClient
}

func NewFedimintClient(baseURL, password string, federationId string) *FedimintClient {
	fc := &FedimintClient{
		BaseURL:      baseURL + "/fedimint/v2",
		Password:     password,
		FederationId: federationId,
	}
	fc.Ln.Client = fc
	fc.Wallet.Client = fc
	fc.Mint.Client = fc

	return fc
}

func (fc *FedimintClient) fetchWithAuth(endpoint string, method string, body []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, fc.BaseURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+fc.Password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error! status: %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func (fc *FedimintClient) get(endpoint string) ([]byte, error) {
	return fc.fetchWithAuth(endpoint, "GET", nil)
}

func (fc *FedimintClient) post(endpoint string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	fmt.Println("jsonBody: ", string(jsonBody))
	if err != nil {
		return nil, err
	}
	return fc.fetchWithAuth(endpoint, "POST", jsonBody)
}

func (fc *FedimintClient) Info() (*types.InfoResponse, error) {
	resp, err := fc.get("/admin/info")
	if err != nil {
		return nil, err
	}
	var infoResp types.InfoResponse
	err = json.Unmarshal(resp, &infoResp)
	if err != nil {
		return nil, err
	}
	return &infoResp, nil
}

func (fc *FedimintClient) Backup(metadata *types.BackupRequest, federationId *string) error {
	_, err := fc.post("/admin/backup", metadata)
	return err
}

func (fc *FedimintClient) DiscoverVersion() (*types.FedimintResponse, error) {
	resp, err := fc.get("/admin/discover-version")
	if err != nil {
		return nil, err
	}
	var versionResp types.FedimintResponse
	err = json.Unmarshal(resp, &versionResp)
	if err != nil {
		return nil, err
	}
	return &versionResp, nil
}

func (fc *FedimintClient) ListOperations(request *types.ListOperationsRequest, federationId *string) (*types.OperationOutput, error) {
	resp, err := fc.post("/admin/list-operations", request)
	if err != nil {
		return nil, err
	}
	var operationsResp types.OperationOutput
	err = json.Unmarshal(resp, &operationsResp)
	if err != nil {
		return nil, err
	}
	return &operationsResp, nil
}

func (fc *FedimintClient) Config() (*types.FedimintResponse, error) {
	resp, err := fc.get("/admin/config")
	if err != nil {
		return nil, err
	}
	var configResp types.FedimintResponse
	err = json.Unmarshal(resp, &configResp)
	if err != nil {
		return nil, err
	}
	return &configResp, nil
}

////////////
// Wallet //
////////////

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

//////////
// mint //
//////////

func (mint *MintModule) Reissue(request ReissueRequest, federationId *string) (*ReissueResponse, error) {
	resp, err := mint.Client.post("/mint/reissue", request)
	if err != nil {
		return nil, err
	}
	var reissueResp ReissueResponse
	err = json.Unmarshal(resp, &reissueResp)
	if err != nil {
		return nil, err
	}
	return &reissueResp, nil
}

func (mint *MintModule) Spend(request SpendRequest, federationId *string) (*SpendResponse, error) {
	resp, err := mint.Client.post("/mint/spend", request)
	if err != nil {
		return nil, err
	}
	var spendResp SpendResponse
	err = json.Unmarshal(resp, &spendResp)
	if err != nil {
		return nil, err
	}
	return &spendResp, nil
}

func (mint *MintModule) Validate(request ValidateRequest, federationId *string) (*ValidateResponse, error) {
	resp, err := mint.Client.post("/mint/validate", request)
	if err != nil {
		return nil, err
	}
	var validateResp ValidateResponse
	err = json.Unmarshal(resp, &validateResp)
	if err != nil {
		return nil, err
	}
	return &validateResp, nil
}

func (mint *MintModule) Split(request SplitRequest) (*SplitResponse, error) {
	resp, err := mint.Client.post("/mint/split", request)
	if err != nil {
		return nil, err
	}
	var splitResp SplitResponse
	err = json.Unmarshal(resp, &splitResp)
	if err != nil {
		return nil, err
	}
	return &splitResp, nil
}

func (mint *MintModule) Combine(request CombineRequest) (*CombineResponse, error) {
	resp, err := mint.Client.post("/mint/combine", request)
	if err != nil {
		return nil, err
	}
	var combineResp CombineResponse
	err = json.Unmarshal(resp, &combineResp)
	if err != nil {
		return nil, err
	}
	return &combineResp, nil
}

////////
// ln //
////////

func (ln *LnModule) CreateInvoice(request LnInvoiceRequest, federationId *string) (*LnInvoiceResponse, error) {
	fmt.Println("request: ", request)
	resp, err := ln.Client.post("/ln/invoice", request)
	if err != nil {
		return nil, err
	}
	var invoiceResp LnInvoiceResponse
	err = json.Unmarshal(resp, &invoiceResp)
	if err != nil {
		return nil, err
	}
	return &invoiceResp, nil
}

func (ln *LnModule) AwaitInvoice(request AwaitInvoiceRequest, federationId *string) (*types.InfoResponse, error) {
	resp, err := ln.Client.post("/ln/await-invoice", request)
	if err != nil {
		return nil, err
	}
	var infoResp types.InfoResponse
	err = json.Unmarshal(resp, &infoResp)
	if err != nil {
		return nil, err
	}
	return &infoResp, nil
}

func (ln *LnModule) Pay(request LnPayRequest, federationId *string) (*LnPayResponse, error) {
	resp, err := ln.Client.post("/ln/pay", request)
	if err != nil {
		return nil, err
	}
	var payResp LnPayResponse
	err = json.Unmarshal(resp, &payResp)
	if err != nil {
		return nil, err
	}
	return &payResp, nil
}

func (ln *LnModule) AwaitPay(request AwaitLnPayRequest, federationId *string) (*LnPayResponse, error) {
	resp, err := ln.Client.post("/ln/await-pay", request)
	if err != nil {
		return nil, err
	}
	var payResp LnPayResponse
	err = json.Unmarshal(resp, &payResp)
	if err != nil {
		return nil, err
	}
	return &payResp, nil
}

func (ln *LnModule) ListGateways() ([]Gateway, error) {
	resp, err := ln.Client.get("/ln/list-gateways")
	if err != nil {
		return nil, err
	}
	var gateways []Gateway
	err = json.Unmarshal(resp, &gateways)
	if err != nil {
		return nil, err
	}
	return gateways, nil
}

func (ln *LnModule) SwitchGateway(request SwitchGatewayRequest, federationId *string) (*Gateway, error) {
	resp, err := ln.Client.post("/ln/switch-gateway", request)
	if err != nil {
		return nil, err
	}
	var gateway Gateway
	err = json.Unmarshal(resp, &gateway)
	if err != nil {
		return nil, err
	}
	return &gateway, nil
}
