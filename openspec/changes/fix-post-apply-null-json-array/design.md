# Design: Optional JSON-encoded arrays and `local_path`

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

## `projects.local_path` — read-only (AWX-assigned)

AWX chooses the canonical directory name under `PROJECTS_ROOT`; user-supplied `local_path` values caused plan vs read-back mismatches.

**Contract**

- Manifest override: `computed: true`, `readOnly: true` for `projects.local_path`.
- `FieldSpec.ReadOnly` drives Terraform schema: **Optional: false**, **Computed: true** (not configurable).
- `payloadFromConfig` skips read-only fields so `local_path` is never sent on create/update.
- `setState` does not apply the “preserve null for computed optional string” shortcut to read-only strings, so the API value is always stored.

**Breaking change**

- Existing configurations that set `local_path` on `awx_project` must remove that argument; Terraform will error until they do.
