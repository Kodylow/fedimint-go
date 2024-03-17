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

	fc := NewFedimintClient(baseURL, password, federationID)

	assert.NotNil(t, fc)

	assert.Equal(t, baseURL+"/fedimint/v2", fc.BaseURL)
	assert.Equal(t, password, fc.Password)
	assert.Equal(t, federationID, fc.FederationId)

	assert.Equal(t, fc, fc.Ln.Client)
	assert.Equal(t, fc, fc.Wallet.Client)
	assert.Equal(t, fc, fc.Mint.Client)
}
