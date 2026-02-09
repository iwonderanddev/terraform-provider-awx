# Data Source: awx_credential_input_source

Reads AWX `credential_input_sources` objects.

## Example Usage

```hcl
data "awx_credential_input_source" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `description` (string)
- `input_field_name` (string)
- `metadata` (string)
- `source_credential` (integer)
- `target_credential` (integer)
