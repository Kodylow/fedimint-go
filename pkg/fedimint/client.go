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
	BaseURL  string
	Password string
	Ln       LnModule
	Wallet   WalletModule
	Mint     MintModule
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

func NewFedimintClient(baseURL, password string) *FedimintClient {
	fc := &FedimintClient{
		BaseURL:  baseURL + "/fedimint/v2",
		Password: password,
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

func (fc *FedimintClient) Backup(metadata *types.BackupRequest) error {
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

func (fc *FedimintClient) ListOperations(request *types.ListOperationsRequest) (*types.OperationOutput, error) {
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

// LN Module

func (ln *LnModule) CreateInvoice(request LnInvoiceRequest) (*LnInvoiceResponse, error) {
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

func (ln *LnModule) AwaitInvoice(request AwaitInvoiceRequest) (*types.InfoResponse, error) {
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

func (ln *LnModule) Pay(request LnPayRequest) (*LnPayResponse, error) {
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

func (ln *LnModule) AwaitPay(request AwaitLnPayRequest) (*LnPayResponse, error) {
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

func (ln *LnModule) SwitchGateway(request SwitchGatewayRequest) (*Gateway, error) {
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

type LnInvoiceRequest struct {
	AmountMsat  int    `json:"amount_msat"`
	Description string `json:"description"`
	ExpiryTime  *int   `json:"expiry_time"`
}

type LnInvoiceResponse struct {
	OperationID string `json:"operation_id"`
	Invoice     string `json:"invoice"`
}

type AwaitInvoiceRequest struct {
	OperationID string `json:"operation_id"`
}

type LnPayRequest struct {
	payment_info         string  `json:"payment_info"`
	amount_msat          *int    `json:"amount_msat"`
	finish_in_background bool    `json:"finish_in_background"`
	lnurl_comment        *string `json:"lnurl_comment"`
}

type LnPayResponse struct {
	operation_id string `json:"operation_id"`
	payment_type string `json:"payment_type"`
	contract_id  string `json:"contract_id"`
	fee          int    `json:"fee"`
}

type AwaitLnPayRequest struct {
	operation_id string `json:"operation_id"`
}

type Gateway struct {
	node_pub_key string `json:"node_pub_key"`
	active       bool   `json:"active"`
}

type SwitchGatewayRequest struct {
	gateway_id string `json:"gateway_id"`
}

// Mint Module

func (mint *MintModule) Reissue(request ReissueRequest) (*ReissueResponse, error) {
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

func (mint *MintModule) Spend(request SpendRequest) (*SpendResponse, error) {
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

func (mint *MintModule) Validate(request ValidateRequest) (*ValidateResponse, error) {
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

type FederationIdPrefix struct {
	Zero, One, Two, Three uint8 `json:"zero"`
}

type TieredMulti struct {
	Amount []interface{} `json:"amount"`
}

type Signature struct {
	Zero G1Affine `json:"zero"`
}

type G1Affine struct {
	X, Y     Fp     `json:"x"`
	Infinity Choice `json:"infinity"`
}

type Fp struct {
	Zero []uint64 `json:"zero"`
}

type Choice struct {
	Zero uint8 `json:"zero"`
}

type KeyPair struct {
	Zero []uint8 `json:"zero"`
}

type OOBNotesData struct {
	Notes              *TieredMulti        `json:"notes"`
	FederationIdPrefix *FederationIdPrefix `json:"federation_id_prefix"`
	Default            struct {
		Variant uint64  `json:"variant"`
		Bytes   []uint8 `json:"bytes"`
	} `json:"default"`
}

type OOBNotes struct {
	Zero []OOBNotesData `json:"zero"`
}

type SpendableNote struct {
	Signature Signature `json:"signature"`
	SpendKey  KeyPair   `json:"spend_key"`
}

type ReissueRequest struct {
	Notes OOBNotes `json:"notes"`
}

type ReissueResponse struct {
	AmountMsat uint64 `json:"amount_msat"`
}

type SpendRequest struct {
	AmountMsat   uint64 `json:"amount_msat"`
	AllowOverpay bool   `json:"allow_overpay"`
	Timeout      uint64 `json:"timeout"`
}

type SpendResponse struct {
	Operation string   `json:"operation"`
	Notes     OOBNotes `json:"notes"`
}

type ValidateRequest struct {
	Notes OOBNotes `json:"notes"`
}

type ValidateResponse struct {
	AmountMsat uint64 `json:"amount_msat"`
}

type SplitRequest struct {
	Notes OOBNotes `json:"notes"`
}

type SplitResponse struct {
	Notes map[uint64]OOBNotes `json:"notes"`
}

type CombineRequest struct {
	Notes []OOBNotes `json:"notes"`
}

type CombineResponse struct {
	Notes OOBNotes `json:"notes"`
}

// Wallet Module

func (wallet *WalletModule) createDepositAddress(request DepositAddressRequest) (*DepositAddressResponse, error) {
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

func (wallet *WalletModule) awaitDeposit(request AwaitDepositRequest) (*AwaitDepositResponse, error) {
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

func (wallet *WalletModule) withdraw(request WithdrawRequest) (*WithdrawResponse, error) {
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
