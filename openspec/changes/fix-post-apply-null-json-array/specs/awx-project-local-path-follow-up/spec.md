# awx-project-local-path-follow-up Delta

## Purpose

Address server-assigned `local_path` on `awx_project` so Terraform does not treat user-supplied values as authoritative. The chosen contract is **`local_path` as a read-only computed attribute** (not configurable in Terraform); AWX’s canonical value is always what appears in state after create/read.

## ADDED Requirements

### Requirement: `awx_project.local_path` is not user-configurable

The provider SHALL expose `local_path` on `awx_project` (and the corresponding data source) as **read-only** in the Terraform schema: operators MUST NOT be able to set `local_path` in configuration. The value SHALL come only from AWX on create and refresh.

#### Scenario: Configuration rejects `local_path`

- **WHEN** a user adds `local_path` to an `awx_project` resource
- **THEN** Terraform reports that the argument is not expected (read-only attribute)

#### Scenario: State matches AWX after apply

- **WHEN** a project is created without `local_path` in configuration
- **THEN** state records the `local_path` returned by AWX and `terraform apply` does not fail with post-apply inconsistency for this attribute
