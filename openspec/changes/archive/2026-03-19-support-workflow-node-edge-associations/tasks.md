# Tasks: Support Workflow Node Edge Associations

## 1. Provider Relationship Derivation

- [x] 1.1 Extend
  `internal/openapi/parser.go` to derive workflow template node edge
  relationships for `success_nodes`, `failure_nodes`, and `always_nodes`.
- [x] 1.2 Ensure the generated relationships use
  `workflow_job_template_nodes` as the child object while overriding child
  argument names to `success_node_id`, `failure_node_id`, and
  `always_node_id`.
- [x] 1.3 Restrict the new alias behavior to
  `workflow_job_template_nodes` parents so runtime `workflow_job_nodes` edge
  endpoints do not become managed resources.

## 2. Provider Tests and Generated Outputs

- [x] 2.1 Add parser unit coverage for the three workflow template node edge
  endpoints and assert resource names, paths, and argument names.
- [x] 2.2 Add negative parser coverage asserting runtime
  `workflow_job_nodes` edge endpoints are not emitted as managed
  relationships.
- [x] 2.3 Add provider-side schema or relationship tests covering the new
  self-referential edge resources and their import/identity behavior.
- [x] 2.4 Run `make generate` and verify
  `internal/manifest/relationships.json` and generated docs include the new
  workflow node edge resources.
- [x] 2.5 Run `make validate-manifest`.
- [x] 2.6 Run `make docs`.
- [x] 2.7 Run `make docs-validate`.
- [x] 2.8 Run `make test`.
- [x] 2.9 Run `make build`.

## 3. Terraform Consumer Rewiring

- [x] 3.1 Update
  `/Users/damien/git/ansible/terraform/modules/v5_migration_workflow/main.tf`
  to replace removed inline `success_nodes` and `failure_nodes` wiring with
  explicit workflow node edge association resources.
- [x] 3.2 Preserve the latest committed workflow semantics, including the
  `mirror_instance -> v5_migration_3d_metadatas -> v5_migration_initial_checks
  -> v5_migration_api` sequence, the shared failure path to
  `send_error_and_fail`, and the `prod` bypass from `v5_migration_2d` and
  `v5_migration_3d` directly to `end_workflow`.
- [x] 3.3 Verify Terraform uses the rebuilt provider through
  `dev_overrides` pointing `registry.terraform.io/iwd/awx` at
  `/Users/damien/git/terraform-provider-awx-iwd/dist`.
- [x] 3.4 Run `terraform fmt -recursive` in
  `/Users/damien/git/ansible/terraform`.
- [x] 3.5 Run `terraform validate` in
  `/Users/damien/git/ansible/terraform/infrastructure`.

## 4. Implementation Evidence

- [x] 4.1 Capture the exact provider test/build command outputs.
- [x] 4.2 Capture the exact Terraform formatting and validation outputs from
  `/Users/damien/git/ansible/terraform/infrastructure`.
- [x] 4.3 Summarize provider files changed, infra files changed, and any
  behavior differences from the old inline-edge model.
