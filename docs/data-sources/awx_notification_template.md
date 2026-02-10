# Data Source: awx_notification_template

Reads AWX `notification_templates` objects.

## Example Usage

```hcl
data "awx_notification_template" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `description` (string)
- `messages` (object)
- `name` (string)
- `notification_configuration` (object, Sensitive)
- `notification_type` (string)
- `organization` (integer)
