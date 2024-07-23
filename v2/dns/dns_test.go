package dns

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tailscale/tailscale-client-go/v2/internal/testsupport"
)

func TestClient_DNSNameservers(t *testing.T) {
	t.Parallel()

	client, server := testsupport.NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	expectedNameservers := map[string][]string{
		"dns": {"127.0.0.1"},
	}

	server.ResponseBody = expectedNameservers
	nameservers, err := With(client).DNSNameservers(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/dns/nameservers", server.Path)
	assert.Equal(t, expectedNameservers["dns"], nameservers)
}

func TestClient_DNSPreferences(t *testing.T) {
	t.Parallel()

	client, server := testsupport.NewTestHarness(t)
	server.ResponseCode = http.StatusOK
	server.ResponseBody = &Preferences{
		MagicDNS: true,
	}

	preferences, err := With(client).DNSPreferences(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/dns/preferences", server.Path)
	assert.Equal(t, server.ResponseBody, preferences)
}

func TestClient_DNSSearchPaths(t *testing.T) {
	t.Parallel()

	client, server := testsupport.NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	expectedPaths := map[string][]string{
		"searchPaths": {"test"},
	}

	server.ResponseBody = expectedPaths

	paths, err := With(client).DNSSearchPaths(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/dns/searchpaths", server.Path)
	assert.Equal(t, expectedPaths["searchPaths"], paths)
}

func TestClient_SplitDNS(t *testing.T) {
	t.Parallel()

	client, server := testsupport.NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	expectedNameservers := SplitDnsResponse{
		"example.com": {"1.1.1.1", "1.2.3.4"},
	}

	server.ResponseBody = expectedNameservers
	nameservers, err := With(client).SplitDNS(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/dns/split-dns", server.Path)
	assert.Equal(t, expectedNameservers, nameservers)
}

func TestClient_SetDNSNameservers(t *testing.T) {
	t.Parallel()

	client, server := testsupport.NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	nameservers := []string{"127.0.0.1"}

	assert.NoError(t, With(client).SetDNSNameservers(context.Background(), nameservers))
	assert.Equal(t, http.MethodPost, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/dns/nameservers", server.Path)

	body := make(map[string][]string)
	assert.NoError(t, json.Unmarshal(server.Body.Bytes(), &body))
	assert.EqualValues(t, nameservers, body["dns"])
}

func TestClient_SetDNSPreferences(t *testing.T) {
	t.Parallel()

	client, server := testsupport.NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	preferences := Preferences{
		MagicDNS: true,
	}

	assert.NoError(t, With(client).SetDNSPreferences(context.Background(), preferences))
	assert.Equal(t, http.MethodPost, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/dns/preferences", server.Path)

	var body Preferences
	assert.NoError(t, json.Unmarshal(server.Body.Bytes(), &body))
	assert.EqualValues(t, preferences, body)
}

func TestClient_SetDNSSearchPaths(t *testing.T) {
	t.Parallel()

	client, server := testsupport.NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	paths := []string{"test"}

	assert.NoError(t, With(client).SetDNSSearchPaths(context.Background(), paths))
	assert.Equal(t, http.MethodPost, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/dns/searchpaths", server.Path)

	body := make(map[string][]string)
	assert.NoError(t, json.Unmarshal(server.Body.Bytes(), &body))
	assert.EqualValues(t, paths, body["searchPaths"])
}

func TestClient_UpdateSplitDNS(t *testing.T) {
	t.Parallel()

	client, server := testsupport.NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	nameservers := []string{"1.1.2.1", "3.3.3.4"}
	request := SplitDnsRequest{
		"example.com": nameservers,
	}

	expectedNameservers := SplitDnsResponse{
		"example.com": nameservers,
	}
	server.ResponseBody = expectedNameservers

	resp, err := With(client).UpdateSplitDNS(context.Background(), request)
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPatch, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/dns/split-dns", server.Path)

	body := make(SplitDnsResponse)
	assert.NoError(t, json.Unmarshal(server.Body.Bytes(), &body))
	assert.EqualValues(t, nameservers, body["example.com"])
	assert.Equal(t, expectedNameservers, resp)
}

func TestClient_SetSplitDNS(t *testing.T) {
	t.Parallel()

	client, server := testsupport.NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	nameservers := []string{"1.1.2.1", "3.3.3.4"}
	request := SplitDnsRequest{
		"example.com": nameservers,
	}

	assert.NoError(t, With(client).SetSplitDNS(context.Background(), request))
	assert.Equal(t, http.MethodPut, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/dns/split-dns", server.Path)

	body := make(SplitDnsResponse)
	assert.NoError(t, json.Unmarshal(server.Body.Bytes(), &body))
	assert.EqualValues(t, nameservers, body["example.com"])
}
