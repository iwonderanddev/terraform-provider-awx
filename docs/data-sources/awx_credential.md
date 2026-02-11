# Data Source: awx_credential

Reads AWX `credentials` objects.

## Example Usage

```hcl
data "awx_credential" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `credential_type_id` (integer)
- `description` (string)
- `inputs` (object, Sensitive)
- `name` (string)
- `organization_id` (integer)
- `team_id` (integer, Sensitive)
- `user_id` (integer, Sensitive)
