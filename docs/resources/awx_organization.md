# Resource: awx_organization

Manages AWX `organizations` objects.

## Example Usage

```hcl
resource "awx_organization" "example" {
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `default_environment` (Optional) The default execution environment for jobs run by this organization.
- `description` (Optional) Managed field from AWX OpenAPI schema.
- `max_hosts` (Optional, Computed) Maximum number of hosts allowed to be managed by this organization.
- `name` (Optional) Managed field from AWX OpenAPI schema.
- `opa_query_path` (Optional) The query path for the OPA policy to evaluate prior to job execution. The query path should be formatted as package/rule.

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_organization.example 42
```
