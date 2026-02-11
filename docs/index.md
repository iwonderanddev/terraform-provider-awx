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

### Resource Argument Qualifiers

Generated resource docs under `docs/resources/*` use these qualifiers:

- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX may apply a server-side default and Terraform records the resulting value in state after apply.

## Breaking Changes

Reference fields that link one AWX object to another use an explicit `_id` suffix in Terraform.
If upgrading from older provider releases, rename unsuffixed link fields (for example, `organization` -> `organization_id`) in resources and data sources.

## Compatibility

This provider targets AWX 24.6.1 API v2 only. Runtime-only objects are excluded from managed resources.
