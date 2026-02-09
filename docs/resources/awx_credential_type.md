# Resource: awx_credential_type

Manages AWX `credential_types` objects.

## Example Usage

```hcl
resource "awx_credential_type" "example" {
  kind = "example"
  name = "example"
}
```

## Argument Reference

- `description` (Optional) Managed field from AWX OpenAPI schema.
- `injectors` (Optional) Enter injectors using either JSON or YAML syntax. Refer to the documentation for example syntax.
- `inputs` (Optional) Enter inputs using either JSON or YAML syntax. Refer to the documentation for example syntax.
- `kind` (Required) * `cloud` - Cloud\n* `net` - Network
- `name` (Required) Managed field from AWX OpenAPI schema.

## Attributes Reference

- `id` (String) Numeric AWX object identifier.

## Import

```bash
terraform import awx_credential_type.example 42
```
