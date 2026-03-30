package telegram

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHTTPClientEmptyProxy(t *testing.T) {
	client, err := NewHTTPClient("")

	require.NoError(t, err)
	require.NotNil(t, client)
	require.IsType(t, &http.Client{}, client)

	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, transport)

	require.NotNil(t, transport.Proxy)
	require.Nil(t, transport.DialContext)
}

func TestNewHTTPClientHTTPProxy(t *testing.T) {
	client, err := NewHTTPClient("http://127.0.0.1:8080")

	require.NoError(t, err)
	require.NotNil(t, client)

	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, transport)

	require.NotNil(t, transport.Proxy)
	require.Nil(t, transport.DialContext)
}

func TestNewHTTPClientHTTPSProxy(t *testing.T) {
	client, err := NewHTTPClient("https://127.0.0.1:8443")

	require.NoError(t, err)
	require.NotNil(t, client)

	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, transport)

	require.NotNil(t, transport.Proxy)
	require.Nil(t, transport.DialContext)
}

func TestNewHTTPClientSOCKS5Proxy(t *testing.T) {
	client, err := NewHTTPClient("socks5://127.0.0.1:1080")

	require.NoError(t, err)
	require.NotNil(t, client)

	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, transport)

	require.Nil(t, transport.Proxy)
	require.NotNil(t, transport.DialContext)
}

func TestNewHTTPClientSOCKS5HProxy(t *testing.T) {
	client, err := NewHTTPClient("socks5h://127.0.0.1:1080")

	require.NoError(t, err)
	require.NotNil(t, client)

	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, transport)

	require.Nil(t, transport.Proxy)
	require.NotNil(t, transport.DialContext)
}

func TestNewHTTPClientInvalidProxyURL(t *testing.T) {
	client, err := NewHTTPClient("://bad")

	require.Error(t, err)
	require.Nil(t, client)
	require.ErrorContains(t, err, "parse proxy url")
}

func TestNewHTTPClientUnsupportedProxyScheme(t *testing.T) {
	client, err := NewHTTPClient("ftp://127.0.0.1:21")

	require.Error(t, err)
	require.Nil(t, client)
	require.ErrorContains(t, err, "unsupported proxy scheme")
}
