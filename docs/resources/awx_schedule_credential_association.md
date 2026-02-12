# Resource: awx_schedule_credential_association

Manages `schedule_credential_association` relationships between `schedules`
and `credentials` objects.

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

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `schedule_id` (Number) Parent object numeric ID.
- `credential_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_schedule_credential_association.example \
  12:34
```
