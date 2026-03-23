# Tasks: GitLab CI mirror to GitHub via GitHub App

## 1. GitHub App and GitLab configuration (operators)

- [ ] 1.1 Create a GitHub App with **Contents: Read and write**, generate a private key, install the app on the target owner/org, and record **App ID** and **Installation ID**
- [ ] 1.2 Create or designate a **GitHub mirror repository** (dedicated repo recommended) and set branch protection consistent with `git push --mirror` (or relax rules for that repo)
- [ ] 1.3 In the GitLab project, add CI/CD variables: `GITHUB_APP_CLIENT_ID`, `GITHUB_APP_INSTALLATION_ID`, `GITHUB_APP_PRIVATE_KEY` (File or masked), `GITHUB_MIRROR_REPOSITORY` (`owner/repo`); protect variables to match the default branch as needed

## 2. Implementation in this repository

- [x] 2.1 Add `scripts/ci/github-installation-token.sh` that builds a GitHub App JWT (RS256), exchanges it for an installation token via `POST /app/installations/{id}/access_tokens`, and prints the token to stdout only (no secrets in logs on success)
- [x] 2.2 Add `.gitlab-ci.yml` with `GIT_DEPTH: "0"`, a `mirror_to_github` (or similarly named) job, `rules` limiting to `$CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH`, Alpine or equivalent image with `git`, `curl`, `jq`, `openssl`, `bash`, and steps to add the GitHub remote and `git push --mirror`
- [x] 2.3 Document operator steps in [`AGENTS.md`](../../../AGENTS.md) (short subsection: variables, GitHub App permissions, risks of `--mirror`)

## 3. Verification

- [x] 3.1 Run `openspec validate --changes --strict` (or validate this change) and fix any reported issues
- [x] 3.2 After merge to the default branch on GitLab, confirm the pipeline succeeds and the GitHub mirror receives expected refs; confirm logs do not contain the installation token
