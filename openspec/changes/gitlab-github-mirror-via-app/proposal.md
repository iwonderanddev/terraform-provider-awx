# Proposal: GitLab CI mirror to GitHub via GitHub App

## Why

The project source of truth may live on GitLab while a public or secondary copy on GitHub is desired. Automating a **full ref mirror** on each default-branch push avoids manual sync and avoids long-lived personal access tokens by using a **GitHub App** installation token.

## What Changes

- Add a **GitLab CI/CD pipeline** (`.gitlab-ci.yml`) with a job that runs **only on pushes to the default branch**.
- Use **GitHub App** authentication: short-lived **JWT** (RS256, app private key) exchanged for an **installation access token**, then `git push --mirror` to a configured GitHub repository over HTTPS (`x-access-token`).
- Add a **maintained script** (e.g. under `scripts/ci/`) that performs JWT construction and token exchange (shell + OpenSSL + curl + jq, or equivalent), with **no token leakage** in logs.
- Document **required GitLab CI/CD variables** and **GitHub App permissions** (Contents read/write; installation on target owner/repo). JWT `iss` uses the app **Client ID** (GitHub recommendation); the installation token API still requires the **Installation ID** (see GitHub App authentication docs).
- **Note**: `git push --mirror` can overwrite refs on the GitHub remote; operators should use a **dedicated mirror repository** and align branch protection with automation needs.

## Capabilities

### New Capabilities

- `gitlab-ci-github-mirror`: Repository SHALL provide automated mirroring from GitLab CI to GitHub using GitHub App installation tokens, default-branch-only triggers, full-history checkout, and a documented secret contract.

### Modified Capabilities

- None (no AWX provider runtime or OpenAPI-derived behavior changes).

## Impact

- **New files**: `.gitlab-ci.yml`, `scripts/ci/` helper for GitHub App token exchange.
- **Operational**: GitLab project **CI/CD variables** (protected as appropriate); GitHub **App** registration and installation; optional updates to [`AGENTS.md`](../../../AGENTS.md) or contributor docs for maintainers running the mirror.
- **No impact** on Go provider code, manifests, or generated Terraform schemas unless documentation is explicitly updated.
