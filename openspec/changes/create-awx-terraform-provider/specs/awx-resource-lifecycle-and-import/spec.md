## ADDED Requirements

### Requirement: Consistent CRUD lifecycle semantics
Each managed AWX object resource SHALL implement create, read, update, and delete behavior aligned with AWX API v2 semantics and Terraform lifecycle expectations.

#### Scenario: Read-after-create convergence
- **WHEN** a resource create operation succeeds
- **THEN** the provider reads the created object and stores normalized state attributes for a stable next plan

#### Scenario: Deleted remote object
- **WHEN** a read operation receives a not-found response for an existing Terraform state object
- **THEN** the provider removes the resource from state on refresh

### Requirement: Deterministic import behavior
Object resources SHALL support import using numeric AWX IDs, and relationship resources SHALL support import using composite identifiers.

#### Scenario: Import object resource by numeric ID
- **WHEN** a user imports a resource with a valid numeric AWX object ID
- **THEN** the provider reads the object and populates complete Terraform state

#### Scenario: Import relationship resource by composite ID
- **WHEN** a user imports a relationship resource with a valid `<parent_id>:<child_id>` identifier
- **THEN** the provider resolves both sides and populates relationship state deterministically
