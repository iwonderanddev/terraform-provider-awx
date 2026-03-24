# Changelog

## v0.2.2 (2026-03-24)

### Fixed

- Optional JSON-array object fields: keep Terraform state aligned when AWX returns an empty array (for example project local path lists).

## Unreleased

### Breaking changes

- **Survey spec resources** (`awx_job_template_survey_spec`, `awx_workflow_job_template_survey_spec`): `spec` is now a Terraform **object** (Plugin Framework dynamic/object), not a JSON string. Replace `spec = jsonencode({ ... })` with a map/object literal `spec = { ... }`.
- **`awx_role_definition` / `data.awx_role_definition`**: `permissions` is now **`list(string)`**, not a JSON-encoded string. Replace `permissions = jsonencode([...])` with `permissions = [...]`.
