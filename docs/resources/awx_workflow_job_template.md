# Resource: awx_workflow_job_template

Manages AWX `workflow_job_templates` objects.

## Example Usage

```hcl
resource "awx_workflow_job_template" "example" {
  name = "example"
  extra_vars = { key = "value" }
}
```

## Argument Reference

Argument qualifiers used below:
- `Required`: Must be set in configuration.
- `Optional`: May be omitted.
- `Optional, Computed`: May be omitted; AWX can apply a server-side default and Terraform records the resulting value after apply.

- `allow_simultaneous` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_inventory_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_labels_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_limit_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_scm_branch_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_skip_tags_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_tags_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `ask_variables_on_launch` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `description` (Optional) Managed field from AWX OpenAPI schema.
- `extra_vars` (Optional) Structured extra variables as a Terraform object.
- `inventory_id` (Optional) Inventory applied as a prompt, assuming job template prompts for inventory
- `job_tags` (Optional) Managed field from AWX OpenAPI schema.
- `limit` (Optional) Managed field from AWX OpenAPI schema.
- `name` (Required) Managed field from AWX OpenAPI schema.
- `organization_id` (Optional) The organization used to determine access to this template.
- `scm_branch` (Optional) Managed field from AWX OpenAPI schema.
- `skip_tags` (Optional) Managed field from AWX OpenAPI schema.
- `survey_enabled` (Optional, Computed) Managed field from AWX OpenAPI schema.
- `webhook_credential_id` (Optional) Personal Access Token for posting back the status to the service API
- `webhook_service` (Optional) Service that webhook requests will be accepted from
  - `github` - GitHub
  - `gitlab` - GitLab
  - `bitbucket_dc` - BitBucket DataCenter

## Attributes Reference

- `id` (Number) Numeric AWX object identifier.

## Import

```bash
terraform import awx_workflow_job_template.example 42
```
