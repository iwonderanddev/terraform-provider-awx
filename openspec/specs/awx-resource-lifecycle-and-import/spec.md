# awx-resource-lifecycle-and-import Specification

## Purpose

TBD - created by archiving change create-awx-terraform-provider. Update Purpose after archive.

## Requirements

### Requirement: Consistent CRUD lifecycle semantics

Each managed AWX object resource SHALL implement create, read, update, and delete behavior aligned with AWX API v2 semantics and Terraform lifecycle expectations.

#### Scenario: Read-after-create convergence

- **WHEN** a resource create operation succeeds
- **THEN** the provider reads the created object and stores normalized state attributes for a stable next plan

#### Scenario: Deleted remote object

- **WHEN** a read operation receives a not-found response for an existing Terraform state object
- **THEN** the provider removes the resource from state on refresh

### Requirement: Deterministic import behavior

Object resources SHALL support import using endpoint-aligned object identifiers, and relationship resources SHALL support import using endpoint-aligned relationship identifiers.

#### Scenario: Import object resource by numeric ID

- **WHEN** a user imports a resource with a valid numeric AWX object ID
- **THEN** the provider reads the object and populates complete Terraform state

#### Scenario: Import object resource by detail-key identifier

- **WHEN** a user imports a detail-key object resource with a valid endpoint identifier (for example `settings` category slug)
- **THEN** the provider reads the object and populates complete Terraform state

#### Scenario: Import relationship resource by composite ID

- **WHEN** a user imports a relationship resource with a valid `<primary_id>:<related_id>` identifier
- **THEN** the provider resolves both sides and populates relationship state deterministically

#### Scenario: Import singleton relationship resource by parent ID

- **WHEN** a user imports a parent-scoped singleton relationship resource with a valid `<resource_id>` identifier
- **THEN** the provider resolves the parent-scoped relationship and populates relationship state deterministically
