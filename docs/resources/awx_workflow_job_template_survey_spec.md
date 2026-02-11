# Resource: awx_workflow_job_template_survey_spec

Manages `workflow_job_template_survey_spec` survey specification for `workflow_job_templates` objects.

Breaking change: use `workflow_job_template_id` instead of legacy `parent_id`.

## Example Usage

```hcl
resource "awx_workflow_job_template_survey_spec" "example" {
  workflow_job_template_id = 12
  spec = jsonencode({
    name        = "Example survey"
    description = "Managed by Terraform"
    spec        = []
  })
}
```

## Argument Reference

- `workflow_job_template_id` (Number, Required) Parent object numeric ID.
- `spec` (String, Optional) JSON-encoded survey specification payload.

## Attributes Reference

- `id` (String) Survey specification ID (same as `workflow_job_template_id`).
- `workflow_job_template_id` (Number) Parent object numeric ID.
- `spec` (String) JSON-encoded survey specification payload.

## Import

```bash
terraform import awx_workflow_job_template_survey_spec.example 12
```
