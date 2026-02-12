# Resource: awx_instance_group_instance_association

Manages `instance_group_instance_association` relationships between `instance_groups`
and `instances` objects.

## Example Usage

```hcl
resource "awx_instance_group_instance_association" "example" {
  instance_group_id = 12
  instance_id  = 34
}
```

## Argument Reference

- `instance_group_id` (Number, Required) Parent object numeric ID.
- `instance_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<primary_id>:<related_id>` format.
- `instance_group_id` (Number) Parent object numeric ID.
- `instance_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_instance_group_instance_association.example \
  12:34
```
