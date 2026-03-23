# Proposal: GitLab CI default-branch sync to GitHub via GitHub App

## Why

The project source of truth may live on GitLab while a public or secondary copy
of the **default branch** on GitHub is desired. Automating a default-branch sync
on each GitLab default-branch push avoids manual sync and avoids long-lived
personal access tokens by using a **GitHub App** installation token.

## What Changes

- Add a **GitLab CI/CD pipeline** (`.gitlab-ci.yml`) with a job that runs **only
  on pushes to the default branch**.
- Use **GitHub App** authentication: short-lived **JWT** (RS256, app private
  key) exchanged for an **installation access token**, then force-push only the
  default branch to a configured GitHub repository over HTTPS
  (`x-access-token`).
- Add a **maintained script** (e.g. under `scripts/ci/`) that performs JWT construction and token exchange (shell + OpenSSL + curl + jq, or equivalent), with **no token leakage** in logs.
- Document **required GitLab CI/CD variables** and **GitHub App permissions** (Contents read/write; installation on target owner/repo). JWT `iss` uses the app **Client ID** (GitHub recommendation); the installation token API still requires the **Installation ID** (see GitHub App authentication docs).
- **Note**: the sync force-updates the GitHub default branch to match GitLab, so
  branch protection must allow the automation to update that branch. The job also
  cleans up stale GitLab-only refs (`refs/merge-requests/*`,
  `refs/pipelines/*`) if an older mirror configuration created them.

## Capabilities

### New Capabilities

- `gitlab-ci-github-mirror`: Repository SHALL provide automated default-branch
  synchronization from GitLab CI to GitHub using GitHub App installation tokens,
  default-branch-only triggers, full-history checkout, and a documented secret
  contract.

### Modified Capabilities

- None (no AWX provider runtime or OpenAPI-derived behavior changes).

## Impact

- **New files**: `.gitlab-ci.yml`, `scripts/ci/` helper for GitHub App token exchange.
- **Operational**: GitLab project **CI/CD variables** (protected as
  appropriate); GitHub **App** registration and installation; optional updates to
  [`AGENTS.md`](../../../AGENTS.md) or contributor docs for maintainers running
  the branch sync.
- **No impact** on Go provider code, manifests, or generated Terraform schemas unless documentation is explicitly updated.
