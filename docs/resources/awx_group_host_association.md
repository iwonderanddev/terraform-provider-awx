# Resource: awx_group_host_association

Manages `group_host_association` relationships between `groups`
and `hosts` objects.

## Example Usage

```hcl
resource "awx_group_host_association" "example" {
  group_id = 12
  host_id  = 34
}
```

## Argument Reference

- `group_id` (Number, Required) Parent object numeric ID.
- `host_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `group_id` (Number) Parent object numeric ID.
- `host_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_group_host_association.example \
  12:34
```
