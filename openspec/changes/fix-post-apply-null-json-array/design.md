# Design: Optional JSON-encoded arrays and `local_path` (Option C)

## JSON-encoded array normalization

Manifest fields with `FieldTypeArray` that are **not** native `list(string)` (today only `role_definitions.permissions` uses native list) are stored as Terraform `string` values containing JSON (`json.Marshal` on read).

When configuration **omits** such an attribute:

- `payloadFromConfig` does not send the field.
- AWX may respond with `policy_instance_list: []`.
- Previously, `setState` encoded `[]` as `types.StringValue("[]")`, while the plan had `null` → Terraform post-apply consistency error.

**Behavior after this change**

- Collect prior `types.String` for these array fields from plan/state (`jsonEncodedArrayStringValuesFromSource`).
- In `setState`, before `toTerraformValue`, if the field is optional, non-computed, JSON-encoded array transport, prior is **null**, and the API value is an **empty slice**, write **`types.StringNull()`** instead of `"[]"`.
- If the user explicitly sets `policy_instance_list = jsonencode([])`, prior is the non-null string `"[]"`; normalization does **not** apply, and state remains `"[]"`.

## `projects.local_path` — Option C (out of scope for runtime in this change)

AWX may return a canonical `local_path` that differs from the configured value (for example prefixed with `_{pk}__`), which can trigger the same class of post-apply inconsistency for projects.

**This change does not** add `setState` preservation or mark `local_path` computed-only.

**Follow-up** is owned by the native survey / spec-driven program: [`native-survey-spec-and-role-permissions`](../native-survey-spec-and-role-permissions/) and/or a future OpenSpec change that defines how `awx_project` exposes server-normalized paths (schema, semantic equality, or documentation).

**Interim mitigation for operators**: where AWX rewrites `local_path`, use:

```hcl
lifecycle {
  ignore_changes = [local_path]
}
```

until a first-class contract is implemented.
