# Resource: awx_workflow_job_template_survey_spec

Manages the AWX survey specification for `workflow_job_templates` objects.

## Example Usage

```hcl
resource "awx_workflow_job_template_survey_spec" "example" {
  workflow_job_template_id = 12
  spec = {
    name        = "Example survey"
    description = "Managed by Terraform"
    spec        = []
  }
}
```

## Schema

### Required

- `workflow_job_template_id` (Number, Required) Parent object numeric ID.

### Optional

- `spec` (Object, Optional, Computed) Survey specification payload as a Terraform object (same logical content as the AWX API JSON body).

### Read-Only

- `id` (String, Read-Only) Survey specification ID (same as `workflow_job_template_id`).
- `workflow_job_template_id` (Number, Read-Only) Parent object numeric ID.
## Import

```bash
terraform import awx_workflow_job_template_survey_spec.example <resource_id>
```

## Further Reading

- [AWX Workflow Job Templates](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/workflow_templates.html)
