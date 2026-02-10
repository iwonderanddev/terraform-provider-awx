## 1. Reproduce and Diagnose

- [x] 1.1 Build or select the provider binary under test and ensure `/Users/damien/git/mockshop-in-cloud-2/terraform/infrastructure/awx/dev` is using that exact build.
- [x] 1.2 Run `terraform init` and `terraform apply` in the AWX dev stack, then capture the first failing resource operation and full diagnostic output.
- [x] 1.3 Classify the root cause as provider contract violation, infrastructure contract violation, or mixed cause.

## 2. Implement Corrective Changes

- [x] 2.1 Implement the minimal provider code fix for any provider-side contract violation in the failing path.
- [x] 2.2 Implement the minimal infrastructure configuration fix in the AWX dev stack for any usage-side contract violation.
- [x] 2.3 Verify touched resources preserve key invariants (lifecycle behavior, import ID contracts, sensitive field handling).

## 3. Add Regression Coverage

- [x] 3.1 Add or update automated tests that reproduce the discovered failure mode and assert corrected behavior.
- [x] 3.2 Run relevant provider tests (targeted and broader suite as appropriate) and confirm they pass.

## 4. Validate End-to-End and Document

- [x] 4.1 Re-run `terraform apply` in `/Users/damien/git/mockshop-in-cloud-2/terraform/infrastructure/awx/dev` and confirm successful completion.
- [x] 4.2 Update provider and/or infrastructure documentation for any changed configuration expectations or diagnostics.
- [x] 4.3 Record root cause, implemented fix, and validation evidence in this change for implementation traceability.
