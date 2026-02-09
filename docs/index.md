# Provider: awx

The `awx` provider manages AWX 24.6.1 objects via API v2.

## Example Usage

```hcl
provider "awx" {
  base_url  = var.awx_base_url
  username  = var.awx_username
  password  = var.awx_password
}
```

## Schema

### Required

- `base_url` (String) AWX base URL, for example `https://awx.example.com`.
- `username` (String) HTTP Basic username.
- `password` (String, Sensitive) HTTP Basic password.

### Optional

- `insecure_skip_tls_verify` (Boolean) Skip TLS verification.
- `ca_cert_pem` (String, Sensitive) PEM CA certificate bundle.
- `request_timeout_seconds` (Number) API request timeout.
- `retry_max_attempts` (Number) Retry attempts for retryable failures.
- `retry_backoff_millis` (Number) Initial retry backoff in milliseconds.

## Compatibility

This provider targets AWX 24.6.1 API v2 only. Runtime-only objects are excluded from managed resources.
