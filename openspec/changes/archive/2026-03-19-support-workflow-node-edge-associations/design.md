# Design: Support Workflow Node Edge Associations

## Context

The provider derives relationship resources from OpenAPI paths that look like
`/api/v2/<parent>/{id}/<child>/` and only emits resources when the child
collection maps to a known managed object collection or a curated alias such as
`galaxy_credentials`. AWX already exposes the workflow template node edge
endpoints:

- `/api/v2/workflow_job_template_nodes/{id}/success_nodes/`
- `/api/v2/workflow_job_template_nodes/{id}/failure_nodes/`
- `/api/v2/workflow_job_template_nodes/{id}/always_nodes/`

Those endpoints are currently ignored because `success_nodes`,
`failure_nodes`, and `always_nodes` are not treated as relationship child
collections. A direct alias to `workflow_job_template_nodes` is not sufficient
because the relationship would then use the same Terraform attribute name for
both parent and child IDs.

The consumer migration in `/Users/damien/git/ansible` already removed inline
`success_nodes` and `failure_nodes` arguments from the workflow module to make
the current provider validate. The intended graph is still available from the
pre-migration configuration and should be restored without reintroducing inline
edge arguments.

## Goals / Non-Goals

**Goals:**

- Generate `awx_workflow_job_template_node_success_node_association`,
  `awx_workflow_job_template_node_failure_node_association`, and
  `awx_workflow_job_template_node_always_node_association` through the existing
  manifest pipeline.
- Use unambiguous canonical arguments:
  `workflow_job_template_node_id` plus one of `success_node_id`,
  `failure_node_id`, or `always_node_id`.
- Keep resource registration, docs generation, and import behavior aligned with
  the current relationship-resource architecture.
- Restore the workflow graph in
  `/Users/damien/git/ansible/terraform/modules/v5_migration_workflow/main.tf`
  with explicit relationship resources that match the latest pre-removal node
  semantics.
- Validate the provider build through
  `/Users/damien/git/terraform-provider-awx-iwd/dist` and Terraform
  `dev_overrides`.

**Non-Goals:**

- Reintroducing inline `success_nodes`, `failure_nodes`, or `always_nodes`
  arguments on `awx_workflow_job_template_node`.
- Adding runtime `workflow_job_node_*` relationship resources.
- Running Terraform apply, destroy, import, or state-manipulation commands.

## Decisions

### Decision: Extend generator alias handling instead of adding bespoke resources

Use the current OpenAPI-to-manifest derivation path and add a narrowly scoped
edge-alias rule for workflow template nodes.

Alternatives considered:

- Hand-author the new relationship entries in generated manifests.
  Rejected because generated files are not source-of-truth and the resources
  would disappear on the next `make generate`.
- Add handwritten provider resources outside the manifest system.
  Rejected because it bypasses the architecture used for all other AWX
  relationship resources and creates a permanent maintenance fork.

### Decision: Represent edge targets as the same child object with explicit child ID overrides

The child object for all three resources remains `workflow_job_template_nodes`,
but derivation sets edge-specific `ChildIDAttribute` values:

- `success_node_id`
- `failure_node_id`
- `always_node_id`

This preserves the real object type for docs and CRUD while avoiding
parent/child attribute collisions.

Alternative considered:

- Reuse `workflow_job_template_node_id` for both sides.
  Rejected because the schema would collapse the parent and child arguments.

### Decision: Scope edge aliasing to workflow template nodes only

The alias rule should only fire for parent collection
`workflow_job_template_nodes`. AWX exposes similar runtime edge endpoints under
`workflow_job_nodes`, but those are runtime records and must not become managed
relationship resources as part of this change.

Alternative considered:

- Generalize all `*_nodes` child collections to self-referential relationships.
  Rejected because it would unintentionally grow provider surface for runtime
  workflow-job endpoints.

### Decision: Restore the infra graph with explicit edge resources, not a generic edge matrix

In the Terraform consumer module, model each named workflow transition with an
association resource and use `count` or `for_each` only where the graph already
fans out or branches conditionally.

This keeps the execution graph readable in HCL and matches the existing module
style. The restored semantics should follow the latest committed edge-bearing
module version:

- `update_project -> each update_inventory` on success
- `each update_inventory -> enable_maintenance` on success
- `enable_maintenance -> mirror_instance` on success
- `enable_maintenance -> send_error_and_fail` on failure
- `mirror_instance -> v5_migration_3d_metadatas` on success
- `mirror_instance -> send_error_and_fail` on failure
- `v5_migration_3d_metadatas -> v5_migration_initial_checks` on success
- `v5_migration_3d_metadatas -> send_error_and_fail` on failure
- `v5_migration_initial_checks -> v5_migration_api` on success
- `v5_migration_initial_checks -> send_error_and_fail` on failure
- `v5_migration_api -> v5_migration_2d` and `v5_migration_3d` on success
- `v5_migration_api -> send_error_and_fail` on failure
- `v5_migration_2d -> end_workflow` on success in `prod`, otherwise
  `disable_maintenance`
- `v5_migration_2d -> send_error_and_fail` on failure
- `v5_migration_3d -> end_workflow` on success in `prod`, otherwise
  `disable_maintenance`
- `v5_migration_3d -> send_error_and_fail` on failure
- `disable_maintenance -> end_workflow` on success
- `disable_maintenance -> send_error_and_fail` on failure

Alternative considered:

- Build one generic local edge table and generate every association through a
  single `for_each`.
  Rejected because it obscures the workflow graph and makes the prod-only
  branch less readable.

## Risks / Trade-offs

- [Risk] Self-referential relationship naming may regress generic relationship
  attribute behavior.
  -> Mitigation: add parser and provider schema tests that assert exact
  resource names, paths, and attribute names.

- [Risk] Generator changes may unintentionally expose runtime workflow node
  edge resources.
  -> Mitigation: add explicit negative coverage for `workflow_job_nodes`
  endpoints in parser tests.

- [Risk] Infra edge rewiring may drift from the latest intended migration graph
  because the current file no longer contains inline edges.
  -> Mitigation: treat the latest committed edge-bearing module definition as
  the source of truth and mirror it one-for-one with association resources.

- [Risk] The ansible repo already has unrelated local modifications.
  -> Mitigation: limit edits to
  `terraform/modules/v5_migration_workflow/main.tf` and avoid touching other
  changed files.

## Migration Plan

1. Extend relationship derivation for workflow template node edge aliases and
   explicit child ID attributes.
2. Add parser, manifest helper, and provider relationship schema tests for the
   new self-referential resources.
3. Regenerate manifests/docs and run provider validation and test commands.
4. Build `dist/terraform-provider-awx`.
5. Point Terraform `dev_overrides` for `registry.terraform.io/iwd/awx` to
   `/Users/damien/git/terraform-provider-awx-iwd/dist`.
6. Update
   `/Users/damien/git/ansible/terraform/modules/v5_migration_workflow/main.tf`
   to restore the success/failure graph with explicit association resources.
7. Run `terraform fmt -recursive` and `terraform validate` from
   `/Users/damien/git/ansible/terraform/infrastructure`.

Rollback:

- Revert the provider parser/test/generated-output changes and rebuild the
  previous provider binary.
- Remove the new workflow edge association resources from the infra module and
  return to the edge-less validating state if necessary.

## Open Questions

- None. This change is decision-complete for implementation.
