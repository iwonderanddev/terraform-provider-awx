# Resource: awx_job_template_survey_spec

Manages `job_template_survey_spec` survey specification for `job_templates` objects.

## Example Usage

```hcl
resource "awx_job_template_survey_spec" "example" {
  parent_id = 12
  spec = jsonencode({
    name        = "Example survey"
    description = "Managed by Terraform"
    spec        = []
  })
}
```

## Argument Reference

- `parent_id` (Number, Required) Parent object numeric ID.
- `spec` (String, Optional) JSON-encoded survey specification payload.

## Attributes Reference

- `id` (String) Survey specification ID (same as `parent_id`).
- `parent_id` (Number) Parent object numeric ID.
- `spec` (String) JSON-encoded survey specification payload.

## Import

```bash
terraform import awx_job_template_survey_spec.example 12
```
