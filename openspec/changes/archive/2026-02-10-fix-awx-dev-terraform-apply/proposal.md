## Why

The full AWX deployment at `/Users/damien/git/mockshop-in-cloud-2/terraform/infrastructure/awx/dev` does not currently complete with `terraform apply` when using this provider. This blocks validation of the provider in a production-like usage path and prevents reliable rollout.

## What Changes

- Reproduce the failing `terraform apply` in the AWX dev infrastructure stack and capture the concrete failure mode.
- Determine whether the root cause is in provider behavior, infrastructure usage, or both.
- Implement the minimal corrective changes needed in the provider and/or the AWX dev infrastructure configuration to make apply succeed.
- Add regression coverage so the same failure mode is detected automatically in future changes.
- Update affected documentation when behavior or expected configuration changes.

## Capabilities

### New Capabilities

- `awx-end-to-end-apply-compatibility`: Ensure a complete AWX deployment workflow can run through Terraform plan/apply successfully with deterministic behavior and actionable error handling.

### Modified Capabilities

- None identified at proposal time; if investigation requires requirement-level changes to existing capabilities, those deltas will be captured in this change's specs.

## Impact

- Provider runtime behavior in `/Users/damien/git/terraform-provider-awx-iwd/internal/provider` and/or `/Users/damien/git/terraform-provider-awx-iwd/internal/client`.
- Potential updates to infrastructure usage in `/Users/damien/git/mockshop-in-cloud-2/terraform/infrastructure/awx/dev`.
- Test coverage under `/Users/damien/git/terraform-provider-awx-iwd/internal/provider` and related acceptance paths.
- Documentation updates in `/Users/damien/git/terraform-provider-awx-iwd/docs` if user-facing behavior or configuration expectations change.
