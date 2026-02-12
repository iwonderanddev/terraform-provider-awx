package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	defaultTimeout             = 30 * time.Second
	defaultRetryMaxAttempts    = 3
	defaultRetryInitialBackoff = 500 * time.Millisecond
)

var pathParameterPattern = regexp.MustCompile(`\{[^/}]+\}`)

// Config controls shared AWX API client behavior.
type Config struct {
	BaseURL               string
	Username              string
	Password              string
	InsecureSkipTLSVerify bool
	CACertPEM             string
	Timeout               time.Duration
	RetryMaxAttempts      int
	RetryInitialBackoff   time.Duration
	UserAgent             string
}

// Client provides AWX API helpers used by resources and data sources.
type Client struct {
	baseURL             *url.URL
	httpClient          *http.Client
	username            string
	password            string
	retryMaxAttempts    int
	retryInitialBackoff time.Duration
	userAgent           string
}

// New creates a new AWX API client.
func New(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, errors.New("base URL is required")
	}
	if strings.TrimSpace(cfg.Username) == "" {
		return nil, errors.New("username is required")
	}
	if strings.TrimSpace(cfg.Password) == "" {
		return nil, errors.New("password is required")
	}

	parsedBaseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}
	if parsedBaseURL.Scheme != "https" && parsedBaseURL.Scheme != "http" {
		return nil, fmt.Errorf("base URL must use http or https scheme: %q", parsedBaseURL.Scheme)
	}

	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
	if cfg.InsecureSkipTLSVerify {
		tlsConfig.InsecureSkipVerify = true //nolint:gosec
	}
	if strings.TrimSpace(cfg.CACertPEM) != "" {
		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM([]byte(cfg.CACertPEM)); !ok {
			return nil, errors.New("failed to parse ca_cert_pem")
		}
		tlsConfig.RootCAs = pool
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		TLSClientConfig:       tlsConfig,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	retryMaxAttempts := cfg.RetryMaxAttempts
	if retryMaxAttempts <= 0 {
		retryMaxAttempts = defaultRetryMaxAttempts
	}
	retryInitialBackoff := cfg.RetryInitialBackoff
	if retryInitialBackoff <= 0 {
		retryInitialBackoff = defaultRetryInitialBackoff
	}
	userAgent := strings.TrimSpace(cfg.UserAgent)
	if userAgent == "" {
		userAgent = "terraform-provider-awx-iwd/dev"
	}

	return &Client{
		baseURL: parsedBaseURL,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		username:            cfg.Username,
		password:            cfg.Password,
		retryMaxAttempts:    retryMaxAttempts,
		retryInitialBackoff: retryInitialBackoff,
		userAgent:           userAgent,
	}, nil
}

// Ping validates connectivity and authentication using the API v2 root endpoint.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.DoJSON(ctx, http.MethodGet, "/api/v2/", nil, nil)
	return err
}

// GetObject retrieves an object by AWX identifier.
func (c *Client) GetObject(ctx context.Context, detailPath string, id string) (map[string]any, error) {
	relative := resolvePathParameter(detailPath, id)
	body, err := c.DoJSON(ctx, http.MethodGet, relative, nil, nil)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// CreateObject creates an object at an AWX collection endpoint.
func (c *Client) CreateObject(ctx context.Context, collectionPath string, payload map[string]any) (map[string]any, error) {
	body, err := c.DoJSON(ctx, http.MethodPost, collectionPath, nil, payload)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// UpdateObject applies partial updates to an AWX object.
func (c *Client) UpdateObject(ctx context.Context, detailPath string, id string, payload map[string]any) (map[string]any, error) {
	relative := resolvePathParameter(detailPath, id)
	body, err := c.DoJSON(ctx, http.MethodPatch, relative, nil, payload)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// DeleteObject deletes an object by AWX identifier.
func (c *Client) DeleteObject(ctx context.Context, detailPath string, id string) error {
	relative := resolvePathParameter(detailPath, id)
	_, err := c.DoJSON(ctx, http.MethodDelete, relative, nil, nil)
	if apiErr := asAPIError(err); apiErr != nil && apiErr.StatusCode == http.StatusNotFound {
		return nil
	}
	return err
}

// ListAll retrieves all paginated items from an AWX endpoint.
func (c *Client) ListAll(ctx context.Context, endpointPath string, query url.Values) ([]map[string]any, error) {
	results := make([]map[string]any, 0)
	nextPath := endpointPath
	currentQuery := cloneValues(query)

	for nextPath != "" {
		payload, err := c.DoJSON(ctx, http.MethodGet, nextPath, currentQuery, nil)
		if err != nil {
			return nil, err
		}

		items, ok := payload["results"].([]any)
		if !ok {
			// Non-paginated response, return it as a single-item list.
			results = append(results, payload)
			break
		}

		for _, item := range items {
			obj, ok := item.(map[string]any)
			if !ok {
				continue
			}
			results = append(results, obj)
		}

		nextRaw, _ := payload["next"].(string)
		if strings.TrimSpace(nextRaw) == "" {
			nextPath = ""
			continue
		}

		resolvedPath, resolvedQuery, err := resolvePaginationNext(nextPath, currentQuery, nextRaw)
		if err != nil {
			return nil, err
		}
		nextPath = resolvedPath
		currentQuery = resolvedQuery
	}

	return results, nil
}

func resolvePaginationNext(currentPath string, currentQuery url.Values, nextRaw string) (string, url.Values, error) {
	nextURL, err := url.Parse(nextRaw)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse pagination next URL %q: %w", nextRaw, err)
	}

	baseURL := &url.URL{
		Path:     currentPath,
		RawQuery: cloneValues(currentQuery).Encode(),
	}
	resolvedURL := baseURL.ResolveReference(nextURL)

	nextPath := resolvedURL.Path
	if strings.TrimSpace(nextPath) == "" {
		nextPath = currentPath
	}
	return nextPath, resolvedURL.Query(), nil
}

// FindByField returns deterministic single-result lookup results by exact field value.
func (c *Client) FindByField(ctx context.Context, endpointPath string, field string, target string) ([]map[string]any, error) {
	items, err := c.ListAll(ctx, endpointPath, nil)
	if err != nil {
		return nil, err
	}

	matches := make([]map[string]any, 0)
	for _, item := range items {
		value, ok := item[field]
		if !ok {
			continue
		}
		if fmt.Sprintf("%v", value) == target {
			matches = append(matches, item)
		}
	}
	return matches, nil
}

// Associate creates a relationship between a parent and child object using an AWX association endpoint.
func (c *Client) Associate(ctx context.Context, relationshipPath string, parentID, childID int64) error {
	resolvedPath := resolvePathParameter(relationshipPath, strconv.FormatInt(parentID, 10))

	payload := map[string]any{"id": childID}
	_, err := c.DoJSON(ctx, http.MethodPost, resolvedPath, nil, payload)
	if err == nil {
		return nil
	}

	// Some AWX association endpoints expect explicit disassociate=false.
	payload["disassociate"] = false
	_, err = c.DoJSON(ctx, http.MethodPost, resolvedPath, nil, payload)
	return err
}

// Disassociate removes a relationship between a parent and child object.
func (c *Client) Disassociate(ctx context.Context, relationshipPath string, parentID, childID int64) error {
	resolvedPath := resolvePathParameter(relationshipPath, strconv.FormatInt(parentID, 10))

	payload := map[string]any{
		"id":           childID,
		"disassociate": true,
	}
	_, err := c.DoJSON(ctx, http.MethodPost, resolvedPath, nil, payload)
	if err == nil {
		return nil
	}

	// Fallback for endpoints that support direct DELETE at /<relation>/<related_id>/.
	fallbackPath := path.Join(resolvedPath, strconv.FormatInt(childID, 10)) + "/"
	_, deleteErr := c.DoJSON(ctx, http.MethodDelete, fallbackPath, nil, nil)
	if apiErr := asAPIError(deleteErr); apiErr != nil && apiErr.StatusCode == http.StatusNotFound {
		return nil
	}
	if deleteErr != nil {
		return deleteErr
	}
	return nil
}

// RelationshipExists checks whether parent-child association exists by scanning the relationship list endpoint.
func (c *Client) RelationshipExists(ctx context.Context, relationshipPath string, parentID, childID int64) (bool, error) {
	resolvedPath := resolvePathParameter(relationshipPath, strconv.FormatInt(parentID, 10))
	items, err := c.ListAll(ctx, resolvedPath, nil)
	if err != nil {
		return false, err
	}
	needle := strconv.FormatInt(childID, 10)
	for _, item := range items {
		idAny, ok := item["id"]
		if !ok {
			continue
		}
		if fmt.Sprintf("%v", idAny) == needle {
			return true, nil
		}
	}
	return false, nil
}

func resolvePathParameter(pathTemplate string, identifier string) string {
	if strings.TrimSpace(pathTemplate) == "" {
		return pathTemplate
	}
	if !pathParameterPattern.MatchString(pathTemplate) {
		return pathTemplate
	}
	return pathParameterPattern.ReplaceAllString(pathTemplate, identifier)
}

// DoJSON executes an HTTP request with retry behavior and decodes the JSON response body.
func (c *Client) DoJSON(ctx context.Context, method string, endpointPath string, query url.Values, payload any) (map[string]any, error) {
	requestURL, err := c.buildURL(endpointPath, query)
	if err != nil {
		return nil, err
	}

	var bodyBytes []byte
	if payload != nil {
		bodyBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to encode payload: %w", err)
		}
	}

	var lastErr error
	for attempt := 1; attempt <= c.retryMaxAttempts; attempt++ {
		respBody, resp, reqErr := c.execute(ctx, method, requestURL, bodyBytes)
		if reqErr != nil {
			apiErr := &APIError{
				Method:    method,
				URL:       requestURL,
				Detail:    "request transport failed",
				Retryable: true,
				Err:       reqErr,
			}
			lastErr = apiErr
			if attempt < c.retryMaxAttempts {
				if sleepErr := c.sleepBackoff(ctx, attempt); sleepErr != nil {
					return nil, sleepErr
				}
				continue
			}
			return nil, apiErr
		}

		if isRetryableStatus(resp.StatusCode) {
			apiErr := &APIError{
				Method:     method,
				URL:        requestURL,
				StatusCode: resp.StatusCode,
				Detail:     "retryable response status",
				Body:       string(respBody),
				Retryable:  true,
			}
			lastErr = apiErr
			if attempt < c.retryMaxAttempts {
				if sleepErr := c.sleepBackoff(ctx, attempt); sleepErr != nil {
					return nil, sleepErr
				}
				continue
			}
			return nil, apiErr
		}

		if resp.StatusCode >= 400 {
			apiErr := &APIError{
				Method:     method,
				URL:        requestURL,
				StatusCode: resp.StatusCode,
				Detail:     extractErrorDetail(respBody),
				Body:       string(respBody),
				Retryable:  false,
			}
			return nil, apiErr
		}

		if len(respBody) == 0 || resp.StatusCode == http.StatusNoContent {
			return map[string]any{}, nil
		}

		decoded := make(map[string]any)
		if err := json.Unmarshal(respBody, &decoded); err != nil {
			return nil, &APIError{
				Method:     method,
				URL:        requestURL,
				StatusCode: resp.StatusCode,
				Detail:     "invalid JSON response",
				Body:       string(respBody),
				Retryable:  false,
				Err:        err,
			}
		}
		return decoded, nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, errors.New("request failed without a captured error")
}

func (c *Client) execute(ctx context.Context, method string, requestURL string, body []byte) ([]byte, *http.Response, error) {
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, bodyReader)
	if err != nil {
		return nil, nil, err
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	return respBody, resp, nil
}

func (c *Client) buildURL(endpointPath string, query url.Values) (string, error) {
	if strings.TrimSpace(endpointPath) == "" {
		return "", errors.New("endpoint path is required")
	}

	u := *c.baseURL
	if strings.HasPrefix(endpointPath, "http://") || strings.HasPrefix(endpointPath, "https://") {
		parsed, err := url.Parse(endpointPath)
		if err != nil {
			return "", fmt.Errorf("failed to parse absolute endpoint URL %q: %w", endpointPath, err)
		}
		u = *parsed
	} else {
		u.Path = path.Join(strings.TrimSuffix(c.baseURL.Path, "/"), endpointPath)
		if strings.HasSuffix(endpointPath, "/") && !strings.HasSuffix(u.Path, "/") {
			u.Path = u.Path + "/"
		}
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	return u.String(), nil
}

func (c *Client) sleepBackoff(ctx context.Context, attempt int) error {
	backoff := float64(c.retryInitialBackoff) * math.Pow(2, float64(attempt-1))
	timer := time.NewTimer(time.Duration(backoff))
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func isRetryableStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode >= 500
}

func extractErrorDetail(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return strings.TrimSpace(string(body))
	}
	if detail, ok := payload["detail"].(string); ok {
		return detail
	}
	if msg, ok := payload["error"].(string); ok {
		return msg
	}
	return strings.TrimSpace(string(body))
}

func asAPIError(err error) *APIError {
	if err == nil {
		return nil
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return nil
}

func cloneValues(in url.Values) url.Values {
	if in == nil {
		return nil
	}
	out := make(url.Values, len(in))
	for k, values := range in {
		cloned := make([]string, len(values))
		copy(cloned, values)
		out[k] = cloned
	}
	return out
}
