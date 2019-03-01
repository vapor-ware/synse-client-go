package synse

import (
	// "crypto/tls"
	// "net/http"
	// "net/http/httptest"
	"testing"
	"time"

	// "github.com/vapor-ware/synse-client-go/internal/test"
	// "github.com/vapor-ware/synse-client-go/synse/scheme"

	// "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestNewWebSocketClientV3_NilConfig(t *testing.T) {
	client, err := NewWebSocketClientV3(nil)
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestNewWebSocketClientV3_NoAddress(t *testing.T) {
	client, err := NewWebSocketClientV3(&Options{
		Address: "",
	})
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestNewWebSocketClientV3_NoTLSCertificates(t *testing.T) {
	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
		TLS: TLSOptions{
			// Enable TLS but not provide the certificates.
			Enabled: true,
		},
	})
	assert.Nil(t, client)
	assert.Error(t, err)
}

func TestNewWebSocketClientV3_defaults(t *testing.T) {
	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)

	assert.Equal(t, "localhost:5000", client.GetOptions().Address)
	assert.Equal(t, 45*time.Second, client.GetOptions().WebSocket.HandshakeTimeout)
	assert.Empty(t, client.GetOptions().TLS.CertFile)
	assert.Empty(t, client.GetOptions().TLS.KeyFile)
	assert.False(t, client.GetOptions().TLS.SkipVerify)
}

func TestNewWebSocketClientV3_ValidAddress(t *testing.T) {
	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

func TestNewWebSocketClientV3_ValidAddressAndTimeout(t *testing.T) {
	client, err := NewWebSocketClientV3(&Options{
		Address: "localhost:5000",
		WebSocket: WebSocketOptions{
			HandshakeTimeout: 46 * time.Second,
		},
	})
	assert.NotNil(t, client)
	assert.NoError(t, err)
}

// func TestWebSocketClientV3_200(t *testing.T) {
// 	// tests := []struct {
// 	// 	event    string
// 	// 	in       string
// 	// 	expected interface{}
// 	// }{
// 	// 	{
// 	// 		"request/version",
// 	// 		`
// 	// {
// 	// "version":"3.0.0",
// 	// "api_version":"v3"
// 	// }`,
// 	// 		&scheme.Version{
// 	// 			Version:    "3.0.0",
// 	// 			APIVersion: "v3",
// 	// 		},
// 	// 	},
// 	// }

// 	// var upgrader = websocket.Upgrader{}

// 	client, err := NewWebSocketClientV3(&Options{
// 		Address: "TODO",
// 	})
// 	assert.NotNil(t, client)
// 	assert.NoError(t, err)

// 	err = client.Open()
// 	assert.NoError(t, err)
// }
