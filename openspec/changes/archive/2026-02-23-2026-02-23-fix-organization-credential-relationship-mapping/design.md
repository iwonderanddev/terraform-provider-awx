# Design: Fix Organization Credential Relationship Mapping

## Context

The current derivation rule treats any `/api/v2/<parent>/{id}/<child>/` endpoint
with `GET+POST` as a relationship candidate and resolves the child collection
name to a managed object collection. For organizations and credentials, AWX
exposes both:

- `/api/v2/organizations/{id}/credentials/`
- `/api/v2/organizations/{id}/galaxy_credentials/`

Both can conceptually map to a credential relationship name, but only the
Galaxy path matches expected association behavior for this provider resource.

## Goals / Non-Goals

**Goals:**

- Ensure `organization_credential_association` resolves to the correct AWX path.
- Make path selection deterministic when aliases collide.
- Preserve resource names and import ID contracts.
- Add focused tests preventing recurrence.

**Non-Goals:**

- Global redesign of all relationship heuristics.
- Expanding provider surface with new resource names.
- Removing runtime safeguard in the same change.

## Decisions

### Decision: Add child-collection alias mapping in derivation

Treat `galaxy_credentials` as an alias of managed child object
`credentials` during relationship derivation.

Rationale:

- Keeps resource naming aligned with existing `*_credential_association`.
- Reflects AWX endpoint intent without introducing new Terraform resource names.

### Decision: Add deterministic collision preference for relationship path selection

When two endpoints resolve to the same relationship name, use explicit path
preference rules instead of first-seen lexicographic behavior. For this change,
prefer `/api/v2/organizations/{id}/galaxy_credentials/` over
`/api/v2/organizations/{id}/credentials/`.

Rationale:

- Eliminates accidental endpoint selection due to path ordering.
- Keeps derivation behavior explainable and testable.

### Decision: Keep runtime client normalization as defense-in-depth

Retain the existing client-side path normalization fallback during transition to
new generated metadata.

Rationale:

- Protects users running older manifests/provider artifacts.
- Low risk and narrow scope.

## Risks / Trade-offs

- [Risk] Alias/collision rules might unintentionally affect future endpoints.
  -> Mitigation: limit preference rules to explicit known collisions and cover
  with unit tests.

- [Risk] Generated manifest diff may touch docs output broadly.
  -> Mitigation: run full generate/validate/docs chain and review only
  relationship-path deltas as acceptance criteria.

## Migration Plan

1. Implement parser alias + collision preference.
2. Add parser tests for this collision and non-regression scenarios.
3. Keep/validate client fallback behavior.
4. Regenerate manifests/docs and run validation/test chain.
5. Confirm `organization_credential_association` path in generated
   `relationships.json` is `.../galaxy_credentials/`.

Rollback:

- Revert parser changes and regenerate manifests to previous state.
- Keep client fallback if desired for resilience.
