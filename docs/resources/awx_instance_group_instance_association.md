# Resource: awx_instance_group_instance_association

Manages `instance_group_instance_association` relationships between `instance_groups` and `instances` objects.

## Example Usage

```hcl
resource "awx_instance_group_instance_association" "example" {
  parent_id = 12
  child_id  = 34
}
```

## Argument Reference

- `parent_id` (Number, Required) Parent object numeric ID.
- `child_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.

## Import

```bash
terraform import awx_instance_group_instance_association.example 12:34
```
