# Resource: awx_job_template_label_association

Manages `job_template_label_association` relationships between `job_templates` and `labels` objects.

Breaking change: use `job_template_id` and `label_id` instead of legacy `parent_id` and `child_id`.

## Example Usage

```hcl
resource "awx_job_template_label_association" "example" {
  job_template_id = 12
  label_id  = 34
}
```

## Argument Reference

- `job_template_id` (Number, Required) Parent object numeric ID.
- `label_id` (Number, Required) Child object numeric ID.

## Attributes Reference

- `id` (String) Composite ID in `<parent_id>:<child_id>` format.
- `job_template_id` (Number) Parent object numeric ID.
- `label_id` (Number) Child object numeric ID.

## Import

```bash
terraform import awx_job_template_label_association.example 12:34
```
