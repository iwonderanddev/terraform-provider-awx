# Resource: awx_schedule_credential_association

Manages `schedule_credential_association` relationships between `schedules` and `credentials` objects.

Breaking change: use `schedule_id` and `credential_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_schedule_credential_association" "example" {
  schedule_id = 12
  credential_id  = 34
}
```

## Argument Reference

- `schedule_id` (Number, Required) Parent object numeric ID.
- `credential_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `schedule_id` (Number) Parent object numeric ID.
- `credential_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_schedule_credential_association.example 12:34
```
