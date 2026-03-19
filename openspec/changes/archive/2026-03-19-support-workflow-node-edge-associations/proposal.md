# Proposal: Support Workflow Node Edge Associations

## Why

The provider currently generates workflow node association resources for
credentials, instance groups, and labels, but it does not generate the AWX
workflow edge endpoints for `success_nodes`, `failure_nodes`, and
`always_nodes`. As a result, the infrastructure migration to `iwd/awx` had to
remove workflow node edge wiring from
`/Users/damien/git/ansible/terraform/modules/v5_migration_workflow/main.tf`,
and `terraform validate` only passes after dropping the intended execution
graph.

## What Changes

- Extend relationship derivation so
  `/api/v2/workflow_job_template_nodes/{id}/{success|failure|always}_nodes/`
  generate dedicated Terraform relationship resources.
- Model these self-referential relationships with explicit edge-specific child
  identifiers: `success_node_id`, `failure_node_id`, and `always_node_id`.
- Exclude runtime `workflow_job_nodes` edge endpoints from provider surface
  expansion.
- Add focused parser, manifest, schema, and relationship-resource test
  coverage.
- Regenerate manifests and documentation so the new resources are embedded,
  registered, and documented through the normal generator flow.
- Update
  `/Users/damien/git/ansible/terraform/modules/v5_migration_workflow/main.tf`
  to replace removed inline `success_nodes` and `failure_nodes` arguments with
  explicit association resources that preserve the current workflow graph.
- Build `dist/terraform-provider-awx` and validate the infrastructure repo via
  Terraform `dev_overrides`.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `awx-relationship-resources`: relationship derivation and canonical argument
  naming must cover workflow job template node edge associations.

## Impact

- Provider generator logic in
  `/Users/damien/git/terraform-provider-awx-iwd/internal/openapi`.
- Generated provider metadata in
  `/Users/damien/git/terraform-provider-awx-iwd/internal/manifest`.
- Generated docs in `/Users/damien/git/terraform-provider-awx-iwd/docs`.
- Provider relationship tests in
  `/Users/damien/git/terraform-provider-awx-iwd/internal/provider` and parser
  tests in `/Users/damien/git/terraform-provider-awx-iwd/internal/openapi`.
- Consumer Terraform wiring in
  `/Users/damien/git/ansible/terraform/modules/v5_migration_workflow/main.tf`.
