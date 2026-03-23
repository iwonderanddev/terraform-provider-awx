# Proposal: Fix post-apply inconsistency for optional JSON-encoded array fields

## Why

Terraform `apply` fails with `Provider produced inconsistent result after apply` when optional manifest `array` fields that use JSON-string transport (for example `awx_instance_group.policy_instance_list`) are **omitted** in configuration. The plan value is null, but AWX returns an empty JSON array and the provider encodes it as the string `"[]"`, which does not match the planned null.

This blocks large stacks (for example Ansible/AWX Terraform modules) that create many instance groups without setting `policy_instance_list`.

Separately, `awx_project.local_path` can diverge after apply when AWX returns a server-canonical path (for example `_{id}__{slug}`) that does not match the configured value. That behavior is **not** addressed by runtime changes in this proposal; it is tracked under Option C for the native survey program (see `design.md`).

## What Changes

- Implement state normalization in `object_resource.go` so optional non-native JSON-encoded array fields map an API empty array to Terraform null when the prior plan/state value was null, matching omitted configuration.
- Add unit tests for the normalization and `setState` behavior.
- Document semantics in OpenSpec (this change) and link `projects.local_path` follow-up to [`native-survey-spec-and-role-permissions`](../native-survey-spec-and-role-permissions/).

## Capabilities

### New Capabilities

- None (behavioral fix to existing single-object resource model).

### Modified Capabilities

- `awx-single-object-resource-model`: Optional JSON-encoded array attributes SHALL NOT cause post-apply inconsistency when omitted and AWX returns an empty array.
- Backlog linkage: `awx-project-local-path-follow-up` (new delta in this change) defers `local_path` canonicalization to the native survey initiative.

## Impact

- Runtime: [`internal/provider/object_resource.go`](../../../internal/provider/object_resource.go), tests under [`internal/provider/object_resource_state_test.go`](../../../internal/provider/object_resource_state_test.go).
- No breaking change to user-facing attribute types in this change.
