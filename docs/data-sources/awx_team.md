# Data Source: awx_team

Reads AWX `teams` objects.

## Example Usage

```hcl
data "awx_team" "example" {
  id = 1
}
```

## Schema

### Optional

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

### Read-Only

- `id` (Number, Read-Only) Numeric AWX object ID.
- `description` (String, Read-Only) Value for `description`.
- `name` (String, Read-Only) Value for `name`.
- `organization_id` (Number, Read-Only) Numeric ID of the related AWX organization object.

## Further Reading

- [AWX Teams](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/teams.html)
