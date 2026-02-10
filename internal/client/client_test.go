package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestClientListAllPagination(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v2/items/" {
				return jsonResponse(http.StatusNotFound, map[string]any{"detail": "not found"}), nil
			}
			if req.URL.Query().Get("page") == "2" {
				return jsonResponse(http.StatusOK, map[string]any{
					"count":    2,
					"next":     nil,
					"previous": "https://awx.example.com/api/v2/items/",
					"results": []map[string]any{{
						"id":   2,
						"name": "item-2",
					}},
				}), nil
			}
			return jsonResponse(http.StatusOK, map[string]any{
				"count":    2,
				"next":     "https://awx.example.com/api/v2/items/?page=2",
				"previous": nil,
				"results": []map[string]any{{
					"id":   1,
					"name": "item-1",
				}},
			}), nil
		}),
		Timeout: 5 * time.Second,
	}

	results, err := client.ListAll(context.Background(), "/api/v2/items/", nil)
	if err != nil {
		t.Fatalf("ListAll returned error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestClientListAllPaginationWithQueryOnlyNext(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path != "/api/v2/role_definitions/" {
				return jsonResponse(http.StatusNotFound, map[string]any{"detail": "not found"}), nil
			}
			if req.URL.Query().Get("page") == "2" {
				return jsonResponse(http.StatusOK, map[string]any{
					"count":    2,
					"next":     nil,
					"previous": "?page=1",
					"results": []map[string]any{{
						"id":   2,
						"name": "role-2",
					}},
				}), nil
			}
			return jsonResponse(http.StatusOK, map[string]any{
				"count":    2,
				"next":     "?page=2",
				"previous": nil,
				"results": []map[string]any{{
					"id":   1,
					"name": "role-1",
				}},
			}), nil
		}),
		Timeout: 5 * time.Second,
	}

	results, err := client.ListAll(context.Background(), "/api/v2/role_definitions/", nil)
	if err != nil {
		t.Fatalf("ListAll returned error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestClientRetriesRetryableFailures(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	client.retryMaxAttempts = 3
	client.retryInitialBackoff = 1 * time.Millisecond

	var attempts atomic.Int32
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			current := attempts.Add(1)
			if current <= 2 {
				return jsonResponse(http.StatusBadGateway, map[string]any{"detail": "temporary upstream issue"}), nil
			}
			return jsonResponse(http.StatusOK, map[string]any{"ok": true}), nil
		}),
		Timeout: 5 * time.Second,
	}

	_, err := client.DoJSON(context.Background(), http.MethodGet, "/api/v2/ping/", nil, nil)
	if err != nil {
		t.Fatalf("expected request to succeed after retries, got: %v", err)
	}
	if attempts.Load() != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts.Load())
	}
}

func TestClientSetsBasicAuth(t *testing.T) {
	t.Parallel()

	client, err := New(Config{
		BaseURL:             "https://awx.example.com",
		Username:            "admin",
		Password:            "secret",
		RetryMaxAttempts:    1,
		RetryInitialBackoff: 1 * time.Millisecond,
		Timeout:             5 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			username, password, ok := req.BasicAuth()
			if !ok {
				t.Fatalf("expected basic auth header to be present")
			}
			if username != "admin" || password != "secret" {
				t.Fatalf("unexpected credentials: %s/%s", username, password)
			}
			return jsonResponse(http.StatusOK, map[string]any{"ok": true}), nil
		}),
		Timeout: 5 * time.Second,
	}

	if _, err := client.DoJSON(context.Background(), http.MethodGet, "/api/v2/ping/", nil, nil); err != nil {
		t.Fatalf("DoJSON failed: %v", err)
	}
}

func TestClientNormalizesAPIError(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	client.httpClient = &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusBadRequest, map[string]any{"detail": "bad request payload"}), nil
		}),
		Timeout: 5 * time.Second,
	}

	_, err := client.DoJSON(context.Background(), http.MethodGet, "/api/v2/items/", nil, nil)
	if err == nil {
		t.Fatalf("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", apiErr.StatusCode)
	}
	if apiErr.Detail != "bad request payload" {
		t.Fatalf("unexpected error detail: %q", apiErr.Detail)
	}
}

func TestClientBuildURLUsesBasePath(t *testing.T) {
	t.Parallel()

	client, err := New(Config{
		BaseURL:             "https://awx.example.com/root",
		Username:            "user",
		Password:            "pass",
		RetryMaxAttempts:    1,
		RetryInitialBackoff: 1 * time.Millisecond,
		Timeout:             5 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	query := url.Values{"name": []string{"demo"}}
	got, err := client.buildURL("/api/v2/projects/", query)
	if err != nil {
		t.Fatalf("buildURL returned error: %v", err)
	}

	want := "https://awx.example.com/root/api/v2/projects/?name=demo"
	if got != want {
		t.Fatalf("unexpected url: got=%q want=%q", got, want)
	}
}

func TestResolvePathParameter(t *testing.T) {
	t.Parallel()

	got := resolvePathParameter("/api/v2/settings/{category_slug}/", "system")
	if got != "/api/v2/settings/system/" {
		t.Fatalf("unexpected resolved path: got=%q want=%q", got, "/api/v2/settings/system/")
	}

	got = resolvePathParameter("/api/v2/teams/{id}/users/", "12")
	if got != "/api/v2/teams/12/users/" {
		t.Fatalf("unexpected resolved path: got=%q want=%q", got, "/api/v2/teams/12/users/")
	}
}

func newTestClient(t *testing.T) *Client {
	t.Helper()

	client, err := New(Config{
		BaseURL:             "https://awx.example.com",
		Username:            "user",
		Password:            "pass",
		RetryMaxAttempts:    3,
		RetryInitialBackoff: 1 * time.Millisecond,
		Timeout:             5 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(status int, payload any) *http.Response {
	raw, err := json.Marshal(payload)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal JSON payload: %v", err))
	}
	if !strings.HasSuffix(string(raw), "\n") {
		raw = bytes.TrimSpace(raw)
	}
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(raw)),
	}
}
