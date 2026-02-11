# Data Source: awx_project

Reads AWX `projects` objects.

## Example Usage

```hcl
data "awx_project" "example" {
  id = 1
}
```

## Argument Reference

- `id` (Number, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (Number) Numeric AWX object ID.
- `allow_override` (boolean)
- `credential_id` (integer)
- `default_environment` (integer)
- `description` (string)
- `local_path` (string)
- `name` (string)
- `organization_id` (integer)
- `scm_branch` (string)
- `scm_clean` (boolean)
- `scm_delete_on_update` (boolean)
- `scm_refspec` (string)
- `scm_track_submodules` (boolean)
- `scm_type` (string)
- `scm_update_cache_timeout` (integer)
- `scm_update_on_launch` (boolean)
- `scm_url` (string)
- `signature_validation_credential` (integer)
- `timeout` (integer)
