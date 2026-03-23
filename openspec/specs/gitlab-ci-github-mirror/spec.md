# gitlab-ci-github-mirror Specification

## Purpose

Define requirements for syncing this repository's default branch from GitLab CI
to GitHub using a GitHub App installation token.

## Requirements

### Requirement: Default-branch sync job

The repository SHALL include a GitLab CI job that pushes only the GitLab default
branch to a configured GitHub repository when the pipeline runs for a commit on
the GitLab **default branch** only.

#### Scenario: Feature branch does not mirror

- **WHEN** a pipeline runs for a branch other than the GitLab default branch
- **THEN** the mirror job SHALL NOT run (or SHALL be skipped by rules equivalent
  to default-branch-only execution)

#### Scenario: Default branch syncs only the target branch

- **WHEN** a pipeline runs for a push to the GitLab default branch
- **THEN** the sync job SHALL run and SHALL update
  `refs/heads/<default-branch>` on the configured GitHub remote to the current CI
  commit

#### Scenario: GitLab-specific refs are not retained

- **GIVEN** the configured GitHub remote contains stale `refs/merge-requests/*`
  or `refs/pipelines/*` refs from an earlier configuration
- **WHEN** the default-branch sync job runs
- **THEN** those GitLab-specific refs SHALL be deleted from the GitHub remote

### Requirement: Full history checkout for mirroring

The mirror job SHALL use a non-shallow git checkout so that mirrored refs and
history are complete relative to the GitLab repository state.

#### Scenario: Git depth

- **WHEN** the mirror job executes
- **THEN** CI configuration SHALL set `GIT_DEPTH` to `0` (or equivalent) so the
  clone is not shallow

### Requirement: GitHub App authentication

Mirroring SHALL NOT use long-lived personal access tokens as the primary
authentication mechanism. The job SHALL obtain a **GitHub App installation
access token** by signing a JWT with the app's private key and calling GitHub's
REST API to create an installation token, then use that token for HTTPS git
operations.

#### Scenario: Token exchange

- **WHEN** the mirror job needs to authenticate to GitHub
- **THEN** it SHALL derive an installation access token via the GitHub App JWT +
  installation token flow documented for GitHub Apps

#### Scenario: Git HTTPS

- **WHEN** pushing to GitHub
- **THEN** the job SHALL use HTTPS with credentials compatible with
  `x-access-token` usage and SHALL set `GIT_TERMINAL_PROMPT=0` to avoid
  interactive prompts

### Requirement: Secret handling

Tokens and private keys SHALL NOT be printed to job logs. Scripts SHALL avoid
exposing secrets under shell tracing (`set -x`).

#### Scenario: Log safety

- **WHEN** the mirror job runs
- **THEN** installation access tokens and PEM material SHALL NOT appear in
  stdout/stderr of successful runs

### Requirement: Configurable target and credentials

Behavior SHALL be driven by CI/CD variables (or file variables) including:
GitHub App **Client ID** (for JWT `iss`), **Installation ID** (for the
installation access token API path; distinct from Client ID), private key PEM,
and target repository `owner/repo` slug. Operators SHALL be able to mark
variables protected to restrict exposure to protected branches.

#### Scenario: Required configuration

- **WHEN** required variables are unset or invalid
- **THEN** the job SHALL fail with a clear error without partially pushing to
  GitHub
