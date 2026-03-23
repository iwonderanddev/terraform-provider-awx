#!/usr/bin/env bash
# Exchange a GitHub App JWT for a short-lived installation access token (stdout only).
# Required env: GITHUB_APP_CLIENT_ID, GITHUB_APP_INSTALLATION_ID, GITHUB_APP_PRIVATE_KEY
# - GITHUB_APP_CLIENT_ID: JWT `iss` (GitHub recommends the app Client ID for `iss`; see GitHub docs).
# - GITHUB_APP_INSTALLATION_ID: required for POST /app/installations/{id}/access_tokens (not interchangeable with Client ID).
# GITHUB_APP_PRIVATE_KEY: PEM string, or path when the value points to an existing file (GitLab File variables).
set -euo pipefail

: "${GITHUB_APP_CLIENT_ID:?GITHUB_APP_CLIENT_ID is required (GitHub App Client ID for JWT iss)}"
: "${GITHUB_APP_INSTALLATION_ID:?GITHUB_APP_INSTALLATION_ID is required}"
: "${GITHUB_APP_PRIVATE_KEY:?GITHUB_APP_PRIVATE_KEY is required}"

key_file=$(mktemp)
cleanup() { rm -f "$key_file"; }
trap cleanup EXIT

if [ -f "${GITHUB_APP_PRIVATE_KEY}" ]; then
	cp "${GITHUB_APP_PRIVATE_KEY}" "$key_file"
else
	printf '%b' "${GITHUB_APP_PRIVATE_KEY}" > "$key_file"
fi
chmod 600 "$key_file"

b64url_encode() {
	openssl base64 -A 2>/dev/null | tr '+/' '-_' | tr -d '='
}

header=$(printf '%s' '{"alg":"RS256","typ":"JWT"}' | b64url_encode)

now=$(date +%s)
iat=$((now - 60))
exp=$((iat + 600))
payload=$(
	jq -cn --argjson iat "$iat" --argjson exp "$exp" --arg iss "$GITHUB_APP_CLIENT_ID" \
		'{iat: $iat, exp: $exp, iss: $iss}' | b64url_encode
)

sign_input="${header}.${payload}"
sig=$(printf '%s' "${sign_input}" | openssl dgst -sha256 -sign "${key_file}" | b64url_encode)

jwt="${sign_input}.${sig}"

# Do not log jwt or response body.
token_response="$(
	curl -sS -f -X POST \
		-H "Authorization: Bearer ${jwt}" \
		-H "Accept: application/vnd.github+json" \
		-H "X-GitHub-Api-Version: 2022-11-28" \
		-H "User-Agent: terraform-provider-awx-gitlab-ci" \
		"https://api.github.com/app/installations/${GITHUB_APP_INSTALLATION_ID}/access_tokens"
)"

echo "${token_response}" | jq -er '.token'
