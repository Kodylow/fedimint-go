// Client_tests
package fedimint

import (
	"fedimint-go-client/pkg/fedimint/types/modules"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateNewFedimintClient() *FedimintClient {
	// Define test data
	baseURL := "http://localhost:3333"
	password := "password"
	federationID := "federation123"

	fc := NewFedimintClient(baseURL, password, federationID)

	return fc
}

func TestNewFedimintClient(t *testing.T) {
	fc := CreateNewFedimintClient()
	assert.NotNil(t, fc)

	assert.Equal(t, fc.BaseURL, "http://localhost:3333/fedimint/v2")
	assert.Equal(t, fc.Password, "password")
	assert.Equal(t, fc.FederationId, "federation123")

	assert.Equal(t, fc, fc.Ln.Client)
	assert.Equal(t, fc, fc.Wallet.Client)
	assert.Equal(t, fc, fc.Mint.Client)
}

func TestGetActiveFederationId(t *testing.T) {
	fc := CreateNewFedimintClient()

	fedId := fc.getActiveFederationId()
	assert.Equal(t, fedId, "federation123")
}

func TestSetActiveFederationId(t *testing.T) {
	fc := CreateNewFedimintClient()
	new_fedId := "New_federation123"

	fedId_prev := fc.FederationId
	fc.setActiveFederationId(new_fedId)
	fedId_now := fc.FederationId
	assert.Equal(t, fedId_now, "New_federation123")
	assert.NotEqual(t, fedId_now, fedId_prev)
}

////////////
// Wallet //
////////////

func TestCreateDepositAddress(t *testing.T) {
	fc := CreateNewFedimintClient()

	depositAddressRequest := modules.DepositAddressRequest{
		Timeout: 3600,
	}

	depositResponse, err := fc.Wallet.createDepositAddress(depositAddressRequest, &fc.FederationId)
	if err != nil {
		assert.Equal(t, depositResponse, nil)
		assert.Equal(t, depositResponse.OperationID, nil)
		assert.Equal(t, depositResponse.Address, nil)
	} else {
		assert.Equal(t, err, nil)
		assert.NotEqual(t, depositResponse.OperationID, nil)
		assert.NotEqual(t, depositResponse.Address, nil)
	}

	awaitDepositRequest := modules.AwaitDepositRequest{
		OperationID: depositResponse.OperationID,
	}

	_, err = fc.Wallet.awaitDeposit(awaitDepositRequest, &fc.FederationId)
	if err != nil {
		fmt.Println("Error awaiting deposit: ", err)
		return
	}
}

func TestWithdraw(t *testing.T) {
	fc := CreateNewFedimintClient()

	withdrawRequest := modules.WithdrawRequest{
		Address:    "UNKNOWN",
		AmountMsat: "10000",
	}

	withdrawResponse, _ := fc.Wallet.withdraw(withdrawRequest, &fc.FederationId)

	assert.NotEqual(t, withdrawResponse, nil)

	// Intentionally make an error (like - wrong FederationId/request)
	wrong_fed_id := "12112"
	_, err := fc.Wallet.withdraw(withdrawRequest, &wrong_fed_id)
	assert.NotEqual(t, err, nil)
}

func TestAwaitWithdraw(t *testing.T) {
	fc := CreateNewFedimintClient()

	depositAddressRequest := modules.DepositAddressRequest{
		Timeout: 3600,
	}
	depositResponse, _ := fc.Wallet.createDepositAddress(depositAddressRequest, &fc.FederationId)

	awaitDepositRequest := modules.AwaitDepositRequest{
		OperationID: depositResponse.OperationID,
	}

	awaitDepositResponse, err := fc.Wallet.awaitDeposit(awaitDepositRequest, &fc.FederationId)
	if err != nil {
		assert.Equal(t, awaitDepositResponse, nil)
		assert.Equal(t, awaitDepositResponse.Status, nil)
	} else {
		assert.Equal(t, err, nil)
		// println(awaitDepositResponse.Status)
		assert.NotEqual(t, awaitDepositResponse.Status, nil)
	}

	// intentionally giving wrong parameters
	wrong_fed_id := "12112"
	_, err1 := fc.Wallet.awaitDeposit(awaitDepositRequest, &wrong_fed_id)
	assert.NotEqual(t, err1, nil)
}
