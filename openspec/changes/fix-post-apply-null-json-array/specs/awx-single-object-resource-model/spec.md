# awx-single-object-resource-model Delta

## ADDED Requirements

### Requirement: Optional JSON-encoded array fields do not drift on empty API arrays

For managed object resources, manifest `array` fields that are configured and stored as JSON-encoded Terraform strings (not native Terraform lists) SHALL normalize AWX empty array responses so that when the prior plan or state value for that attribute is **null** (attribute omitted in configuration) and the API returns an **empty array**, the provider SHALL write **null** to state for that attribute instead of the JSON string `"[]"`.

This SHALL preserve explicit configuration where the operator sets the attribute to an empty JSON array (for example `jsonencode([])`), in which case the prior value is non-null and state SHALL continue to reflect the encoded empty array as today.

#### Scenario: Omitted optional JSON array matches null after apply

- **WHEN** a resource omits an optional JSON-encoded array field in Terraform configuration
- **AND** AWX returns that field as an empty JSON array
- **THEN** the provider state for that attribute is null and `terraform apply` does not report a post-apply inconsistency for that attribute

#### Scenario: Explicit empty JSON array remains encoded

- **WHEN** configuration sets the JSON-encoded array field to an explicit empty JSON array
- **AND** AWX returns an empty JSON array for that field
- **THEN** the provider state retains the non-null encoded empty array value consistent with configuration
