# Resource: awx_job_template_survey_spec

Manages the AWX survey specification for `job_templates` objects.

## Example Usage

```hcl
resource "awx_job_template_survey_spec" "example" {
  job_template_id = 12
  spec = {
    name        = "Example survey"
    description = "Managed by Terraform"
    spec        = []
  }
}
```

## Schema

### Required

- `job_template_id` (Number, Required) Parent object numeric ID.

### Optional

- `spec` (Object, Optional, Computed) Survey specification payload as a Terraform object (same logical content as the AWX API JSON body).

### Read-Only

- `id` (String, Read-Only) Survey specification ID (same as `job_template_id`).
- `job_template_id` (Number, Read-Only) Parent object numeric ID.
## Import

```bash
terraform import awx_job_template_survey_spec.example <resource_id>
```

## Further Reading

- [AWX Job Templates](https://docs.ansible.com/projects/awx/en/24.6.1/userguide/job_templates.html)
