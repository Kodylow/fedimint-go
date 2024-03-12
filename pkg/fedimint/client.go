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
