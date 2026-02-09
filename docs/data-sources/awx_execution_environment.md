# Data Source: awx_execution_environment

Reads AWX `execution_environments` objects.

## Example Usage

```hcl
data "awx_execution_environment" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `credential` (integer)
- `description` (string)
- `image` (string)
- `name` (string)
- `organization` (integer)
- `pull` (string)
