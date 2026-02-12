# Data Source: awx_label

Reads AWX `labels` objects.

## Example Usage

```hcl
data "awx_label" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `name` (String, Read-Only) Value for `name`.
- `organization_id` (Number, Read-Only) Organization this label belongs to.

## Further Reading

- [AWX Job Templates](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/job_templates.html)
