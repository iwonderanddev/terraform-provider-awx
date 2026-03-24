# Proposal: Fix post-apply inconsistency for optional JSON-encoded array fields

## Why

Terraform `apply` fails with `Provider produced inconsistent result after apply` when optional manifest `array` fields that use JSON-string transport (for example `awx_instance_group.policy_instance_list`) are **omitted** in configuration. The plan value is null, but AWX returns an empty JSON array and the provider encodes it as the string `"[]"`, which does not match the planned null.

This blocks large stacks (for example Ansible/AWX Terraform modules) that create many instance groups without setting `policy_instance_list`.

Separately, `awx_project.local_path` can diverge after apply when AWX returns a server-canonical path (for example `_{id}__{slug}`) that does not match a user-configured value.

## What Changes

- Implement state normalization in `object_resource.go` so optional non-native JSON-encoded array fields map an API empty array to Terraform null when the prior plan/state value was null, matching omitted configuration.
- Add unit tests for the normalization and `setState` behavior.
- Make `awx_project.local_path` **read-only computed** in schema (`FieldSpec.readOnly` + field override) so it cannot be set in Terraform; document via OpenSpec and generated docs.

## Capabilities

### New Capabilities

- None (behavioral fix to existing single-object resource model).

### Modified Capabilities

- `awx-single-object-resource-model`: Optional JSON-encoded array attributes SHALL NOT cause post-apply inconsistency when omitted and AWX returns an empty array.
- `awx-project-local-path-follow-up`: `local_path` is server-assigned and exposed as read-only in Terraform.

## Impact

- Runtime: [`internal/provider/object_resource.go`](../../../internal/provider/object_resource.go), [`internal/manifest/manifest.go`](../../../internal/manifest/manifest.go), [`internal/openapi/overrides.go`](../../../internal/openapi/overrides.go), [`internal/manifest/field_overrides.json`](../../../internal/manifest/field_overrides.json), [`cmd/awxgen/main.go`](../../../cmd/awxgen/main.go), tests.
- **Breaking:** Configurations that set `local_path` on `awx_project` must remove that argument.
