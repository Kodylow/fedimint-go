// Client_tests
package fedimint

import (
	"testing"

	"github.com/stretchr/testify/assert" // Import testify for assertions
)

func TestNewFedimintClient(t *testing.T) {
	// Define test data
	baseURL := "http://example.com"
	password := "password"
	federationID := "federation123"

	// Call the function being tested
	fc := NewFedimintClient(baseURL, password, federationID)

	// Assert that the returned FedimintClient is not nil
	assert.NotNil(t, fc)

	// Assert that the fields of the returned FedimintClient are set correctly
	assert.Equal(t, baseURL+"/fedimint/v2", fc.BaseURL)
	assert.Equal(t, password, fc.Password)
	assert.Equal(t, federationID, fc.FederationId)

	// Assert that the Ln, Wallet, and Mint modules have their Client fields set to the returned FedimintClient
	assert.Equal(t, fc, fc.Ln.Client)
	assert.Equal(t, fc, fc.Wallet.Client)
	assert.Equal(t, fc, fc.Mint.Client)
}
