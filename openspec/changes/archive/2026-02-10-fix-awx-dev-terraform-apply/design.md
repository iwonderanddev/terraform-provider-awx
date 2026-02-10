## Context

The AWX development infrastructure stack at `/Users/damien/git/mockshop-in-cloud-2/terraform/infrastructure/awx/dev` currently fails during `terraform apply` when using this provider build. The provider is intended to support complete AWX API v2 lifecycle management, so this failure indicates a contract or behavior gap either in provider runtime logic, the calling infrastructure configuration, or both.

Constraints:
- Preserve existing provider import ID and lifecycle contracts unless a requirement-level change is explicitly documented.
- Keep fixes minimal and targeted to the discovered failure path.
- Maintain compatibility with AWX `24.6.1` and current auth model.

## Goals / Non-Goals

**Goals:**
- Reproduce the apply failure deterministically from the AWX dev stack.
- Isolate root cause to provider implementation, infrastructure usage, or both.
- Implement a durable fix that allows successful apply without manual intervention.
- Add regression coverage so the failure mode is automatically detected in future changes.

**Non-Goals:**
- Broad redesign of resource model generation or manifest architecture.
- Expansion of provider feature surface beyond what is needed to fix this apply path.
- General performance tuning unrelated to the failing behavior.

## Decisions

### Decision: Reproduce against the real failing stack before changing code
Use the reported stack (`mockshop-in-cloud-2/.../awx/dev`) as the primary reproduction harness to avoid fixing an abstract symptom that does not reflect real usage.

Alternatives considered:
- Reproduce only with synthetic acceptance fixtures. Rejected because current issue is already known to occur in a full deployment path and may involve composition effects not present in isolated fixtures.

### Decision: Prefer provider-side correction when provider contract is violated
If the failure is caused by serialization, state handling, request construction, relationship semantics, or diagnostics that violate documented provider behavior, fix it in the provider.

Alternatives considered:
- Work around provider defect only in Terraform configuration. Rejected when provider contract is clearly broken, because this would propagate fragile consumer-side workarounds.

### Decision: Allow infrastructure-side updates when configuration is invalid for current contracts
If the failure is caused by infrastructure configuration that violates documented resource requirements, update the infra code to conform and document the expectation.

Alternatives considered:
- Relax provider contracts to accept invalid or ambiguous config. Rejected because it can weaken determinism and mask user errors.

### Decision: Capture regression at the smallest reliable test layer
Add automated coverage at the narrowest layer that reproduces the bug reliably (unit/provider tests first, acceptance when needed) to keep feedback fast while preventing regressions.

Alternatives considered:
- Rely solely on manual `terraform apply` verification. Rejected because manual validation is insufficient for long-term regression protection.

## Risks / Trade-offs

- [Risk] Failure is environment-specific and hard to reproduce consistently.
  -> Mitigation: capture exact command, provider binary source, and relevant inputs during first reproduction; codify minimal reproducible case in tests.

- [Risk] Fix may require touching both repositories, increasing coordination overhead.
  -> Mitigation: keep change scope explicit; separate provider contract fix from infra conformance updates in commits/tasks.

- [Risk] Regression test may be flaky if it depends on live AWX behavior.
  -> Mitigation: prioritize deterministic provider-level tests; gate live acceptance tests behind existing acceptance controls.

## Migration Plan

1. Reproduce and document the current failure signature from the AWX dev stack.
2. Implement targeted fix in provider and/or infra based on confirmed root cause.
3. Add or update automated tests covering the failing path.
4. Re-run relevant provider test suites and rerun `terraform apply` in the AWX dev stack.
5. Update documentation if configuration expectations or behavior changed.

Rollback:
- Revert the specific provider/infra commits for this change and restore previous provider binary in the dev stack.

## Open Questions

- What exact resource or relationship operation fails first in the AWX dev apply path?
- Does the failure reproduce with the current main branch provider, or only with in-progress local changes?
- Is any infra change required beyond provider correction once root cause is fixed?
