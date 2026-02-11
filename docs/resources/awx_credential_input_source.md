# Resource: awx_credential_input_source

Manages AWX `credential_input_sources` objects.

## Example Usage

```hcl
resource "awx_credential_input_source" "example" {
  input_field_name = "example"
  source_credential_id = 1
  target_credential_id = 1
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `description` (Optional) Managed field from AWX OpenAPI schema.
- `input_field_name` (Required) Managed field from AWX OpenAPI schema.
- `metadata` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `source_credential_id` (Required) Managed field from AWX OpenAPI schema.
- `target_credential_id` (Required) Managed field from AWX OpenAPI schema.

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_credential_input_source.example 42
```
