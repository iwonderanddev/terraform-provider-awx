# Data Source: awx_project

Reads AWX `projects` objects.

## Example Usage

```hcl
data "awx_project" "example" {
  id = 1
}
```

## Argument Reference

- `id` (String, Optional) Numeric AWX object ID.
- `name` (String, Optional) Deterministic exact-name lookup if `id` is omitted.

## Attributes Reference

- `id` (String) Numeric AWX object ID.
- `allow_override` (boolean)
- `credential` (integer)
- `default_environment` (integer)
- `description` (string)
- `local_path` (string)
- `name` (string)
- `organization` (integer)
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
