# awx-end-to-end-apply-compatibility Specification

## Purpose
TBD - created by archiving change fix-awx-dev-terraform-apply. Update Purpose after archive.

## Requirements
### Requirement: AWX dev deployment apply compatibility
The provider and supported Terraform configuration SHALL allow a complete AWX deployment apply workflow to run successfully for the dev stack at `/Users/damien/git/mockshop-in-cloud-2/terraform/infrastructure/awx/dev` when valid AWX connectivity and credentials are supplied.

#### Scenario: Successful end-to-end apply with supported configuration
- **WHEN** an operator runs `terraform apply` in the AWX dev stack using this provider and a provider-compatible configuration
- **THEN** Terraform SHALL complete the apply without provider-caused hard failures in managed AWX object or relationship operations

### Requirement: Root-cause-specific correction
The implementation SHALL correct the concrete failure mode identified during reproduction at the layer where the contract is violated (provider behavior, infrastructure usage, or both) without introducing unrelated behavioral changes.

#### Scenario: Provider contract violation is fixed at provider layer
- **WHEN** the reproduced failure is caused by provider request/state/diagnostic behavior that violates documented provider contracts
- **THEN** the change SHALL modify provider code to restore contract-compliant behavior and remove the failure

#### Scenario: Infrastructure contract violation is fixed at configuration layer
- **WHEN** the reproduced failure is caused by Terraform configuration that violates documented provider input or lifecycle expectations
- **THEN** the change SHALL update infrastructure configuration to conform to provider contracts and remove the failure

### Requirement: Regression protection for discovered failure
The change SHALL include automated regression coverage that fails if the discovered failure mode reappears.

#### Scenario: Regression coverage executes in CI-compatible test flows
- **WHEN** project tests for the touched area are executed
- **THEN** at least one automated test SHALL exercise the corrected failure path and detect reintroduction of the bug

### Requirement: Actionable failure diagnostics
When apply fails for unsupported or invalid inputs, diagnostics SHALL identify the failing resource operation and provide enough detail for operators to resolve the issue without inspecting provider internals.

#### Scenario: Diagnostic identifies failing operation context
- **WHEN** the provider returns an error during AWX operation execution in apply
- **THEN** the diagnostic SHALL include the resource context and operation stage associated with the failure
