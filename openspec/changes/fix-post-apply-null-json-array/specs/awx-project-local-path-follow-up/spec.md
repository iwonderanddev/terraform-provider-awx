# awx-project-local-path-follow-up Delta

## Purpose

Track server-canonical `local_path` behavior for `awx_project` without implementing runtime fixes in the [`fix-post-apply-null-json-array`](../../) change. Implementation is deferred to the **native survey / spec-driven** program ([`native-survey-spec-and-role-permissions`](../../../native-survey-spec-and-role-permissions/)) or a dedicated follow-up OpenSpec change.

## ADDED Requirements

### Requirement: `projects.local_path` canonicalization is specified in the native survey program

The provider SHALL define how `awx_project.local_path` interacts with AWX server normalization (for example directory names prefixed with project primary key) so that plan, apply, and refreshed state remain consistent, **or** SHALL document supported workarounds (such as `lifecycle` blocks) as part of the native survey initiative.

This requirement SHALL NOT be satisfied solely by the JSON-encoded array normalization change; it MAY be implemented together with other `awx_project` field modeling work under [`native-survey-spec-and-role-permissions`](../../../native-survey-spec-and-role-permissions/).

#### Scenario: Follow-up links native survey work

- **WHEN** maintainers close the `fix-post-apply-null-json-array` change
- **THEN** `projects.local_path` post-apply inconsistencies remain tracked under the native survey program until a follow-up change implements the agreed contract
