package genapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInterface for testing
type TestInterface interface {
	Interface
	TestMethod() string
}

func TestNew_ErrorPropagation(t *testing.T) {
	t.Run("ReturnsErrorForUnregisteredInterface", func(t *testing.T) {
		// TestInterface is not registered in the registry
		client, err := New[TestInterface]()
		
		// Should return an error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no registration found for interface")
		
		// Client should be nil/zero value
		assert.Nil(t, client)
	})
}

// TestValidRegistration ensures that registered interfaces work correctly
type ValidTestInterface interface {
	Interface
	ValidMethod() string
}

type validTestImpl struct {
	client HttpClient
}

func (v *validTestImpl) SetHttpClient(client HttpClient) {
	v.client = client
}

func (v *validTestImpl) ValidMethod() string {
	return "valid"
}

func TestNew_ValidRegistration(t *testing.T) {
	// Register a valid interface-implementation pair
	Register[ValidTestInterface, *validTestImpl](Config{})
	
	client, err := New[ValidTestInterface]()
	
	// Should succeed
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, "valid", client.ValidMethod())
}