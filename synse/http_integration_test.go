package synse

import (
	// "encoding/json"
	// "fmt"
	"testing"
	// "time"

	"github.com/stretchr/testify/assert"
	// "github.com/vapor-ware/synse-client-go/synse/scheme"
)

func TestIntegration_Status(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Status()
	assert.NoError(t, err)
	assert.Equal(t, "ok", resp.Status)
	assert.NotNil(t, resp.Timestamp)
}

func TestIntegration_Version(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client, err := NewHTTPClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	resp, err := client.Version()
	assert.NoError(t, err)
	assert.Equal(t, "v3", resp.APIVersion)
	assert.NotNil(t, resp.Version)
}
