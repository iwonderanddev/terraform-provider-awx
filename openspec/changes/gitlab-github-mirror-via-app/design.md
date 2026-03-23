# Design: GitLab CI mirror to GitHub via GitHub App

## Context

The repository is a Terraform provider for AWX; primary development may occur on **GitLab** while a **GitHub** remote serves visibility, forks, or downstream automation. The mirror must run in **GitLab CI** without embedding long-lived PATs. **GitHub Apps** provide short-lived **installation access tokens** derived from a JWT signed with the app’s RSA private key.

Constraints:

- GitLab runner must perform a **full** clone for a faithful mirror (`GIT_DEPTH=0`).
- Secrets live in **GitLab CI/CD variables** (prefer **protected** + **masked** where applicable; private keys often use **File** type).
- The GitHub App needs **Contents: Read and write** on the target repository.

## Goals / Non-Goals

**Goals:**

- On every push to the **default branch**, run one job that **mirrors all refs** to GitHub (`git push --mirror`).
- Authenticate with **GitHub App** JWT → installation token exchange (HTTPS `git` push using `x-access-token`).
- Keep token material out of job logs; avoid `set -x` around secret use.

**Non-Goals:**

- Syncing GitLab merge request pipelines to GitHub Checks or PRs.
- Replacing GitLab with GitHub as the canonical forge.
- Supporting Git LFS in the initial deliverable (can be added later if `.gitattributes` adopts LFS).

## Decisions

1. **Token helper: shell + OpenSSL + curl + jq**  
   - *Rationale*: No Node/npm bootstrap in CI; Alpine-friendly; single executable script.  
   - *Alternative*: `npx @octokit/auth-app` — heavier cold start.

2. **Trigger: default branch only**  
   - *Rationale*: Matches agreed scope; reduces accidental pushes from feature branches.  
   - *Implementation*: `rules` with `$CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH`.

3. **Push mode: `--mirror`**  
   - *Rationale*: Matches “whole repository” mirroring (branches, tags, refs).  
   - *Trade-off*: Can **overwrite/delete** refs on GitHub; use a **dedicated** mirror repo.

4. **Remote URL**  
   - One-off `git push --mirror "https://x-access-token:${TOKEN}@github.com/${GITHUB_MIRROR_REPOSITORY}.git"` so the token is not stored in `.git/config`.

5. **Variables**  
   - `GITHUB_APP_CLIENT_ID` (JWT `iss`; GitHub recommends the app **Client ID** for `iss`), `GITHUB_APP_INSTALLATION_ID` (REST path for installation tokens — **not** interchangeable with Client ID), `GITHUB_APP_PRIVATE_KEY` (PEM), `GITHUB_MIRROR_REPOSITORY` (`owner/repo`).  
   - Script normalizes PEM from file or escaped newlines to a temp file for OpenSSL; JWT payload uses `jq` so `iss` is a proper JSON string.

## Risks / Trade-offs

| Risk | Mitigation |
| --- | --- |
| GitHub **branch protection** blocks force updates | Use a mirror-only repo with permissive rules, or allow the app to bypass where policy allows |
| **JWT / clock skew** | Use NTP-backed runners; JWT `iat` slightly in the past per GitHub docs |
| **Secret in logs** | Never echo token; mask variables; disable verbose tracing around git push |
| **Mirror deletes refs** operators did not expect | Document destructive nature; restrict job to protected default branch |

## Migration Plan

1. Create GitHub App; install on org/user; grant Contents read/write; record App ID and Installation ID.
2. Create empty (or disposable) GitHub repo for the mirror.
3. Configure GitLab CI/CD variables.
4. Merge `.gitlab-ci.yml` and script; verify first pipeline on default branch.

**Rollback**: Remove or disable the job in `.gitlab-ci.yml`; delete GitHub App installation if no longer needed.

## Open Questions

- None blocking implementation; org-specific **branch protection** policy is an operator choice outside the repo.
