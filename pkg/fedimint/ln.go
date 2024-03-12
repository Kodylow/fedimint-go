package fedimint

import (
	"encoding/json"
	"fedimint-go-client/pkg/fedimint/types"
	"fmt"
)

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
	Payment_info         string  `json:"payment_info"`
	Amount_msat          *int    `json:"amount_msat"`
	Finish_in_background bool    `json:"finish_in_background"`
	Lnurl_comment        *string `json:"lnurl_comment"`
}

type LnPayResponse struct {
	Pperation_id string `json:"operation_id"`
	Payment_type string `json:"payment_type"`
	Contract_id  string `json:"contract_id"`
	Fee          int    `json:"fee"`
}

type AwaitLnPayRequest struct {
	Operation_id string `json:"operation_id"`
}

type Gateway struct {
	Node_pub_key string `json:"node_pub_key"`
	Active       bool   `json:"active"`
}

// string::> FederationId
type ListGatewaysResponse map[string][]Gateway

type SwitchGatewayRequest struct {
	Gateway_id string `json:"gateway_id"`
}
