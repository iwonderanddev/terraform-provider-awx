# Resource: awx_host_group_association

Manages `host_group_association` relationships between `hosts` and `groups` objects.

Breaking change: use `host_id` and `group_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_host_group_association" "example" {
  host_id = 12
  group_id  = 34
}
```

## Argument Reference

- `host_id` (Number, Required) Parent object numeric ID.
- `group_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `host_id` (Number) Parent object numeric ID.
- `group_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_host_group_association.example 12:34
```
