## ADDED Requirements

### Requirement: Unit test coverage for provider core behaviors
The provider SHALL include unit tests for API client logic, schema mapping, lifecycle scaffolding, and import/state normalization.

#### Scenario: Unit test execution
- **WHEN** unit tests are run in CI
- **THEN** tests validate provider core logic without requiring a live AWX instance

### Requirement: Local opt-in acceptance and e2e tests
The provider SHALL provide acceptance and e2e test suites that run only when required AWX environment variables and credentials are explicitly supplied.

#### Scenario: Acceptance tests without environment
- **WHEN** acceptance tests are executed without required AWX local environment configuration
- **THEN** tests are skipped with actionable guidance instead of failing as hard errors

#### Scenario: Acceptance tests with environment
- **WHEN** acceptance tests are executed with valid AWX credentials and endpoint configuration
- **THEN** tests verify CRUD, import, and relationship behavior against AWX 24.6.1
