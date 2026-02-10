# Resource: awx_credential_type

Manages AWX `credential_types` objects.

## Example Usage

```hcl
resource "awx_credential_type" "example" {
  kind = "example"
  name = "example"
  injectors = { key = "value" }
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `description` (Optional) Managed field from AWX OpenAPI schema.
- `injectors` (Optional, Computed) Terraform object defining credential type injectors. Refer to the documentation for expected structure.
- `inputs` (Optional, Computed) Terraform object defining credential type inputs. Refer to the documentation for expected structure.
- `kind` (Required) * `cloud` - Cloud\n* `net` - Network
- `name` (Required) Managed field from AWX OpenAPI schema.

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_credential_type.example 42
```
