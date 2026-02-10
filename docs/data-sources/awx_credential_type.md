# Data Source: awx_credential_type

Reads AWX `credential_types` objects.

## Example Usage

```hcl
data "awx_credential_type" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `description` (string)
- `injectors` (object)
- `inputs` (object)
- `kind` (string)
- `name` (string)
