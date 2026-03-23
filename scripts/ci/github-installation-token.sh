#!/usr/bin/env bash
# Exchange a GitHub App JWT for a short-lived installation access token (stdout only).
# Required env: GITHUB_APP_CLIENT_ID, GITHUB_APP_INSTALLATION_ID, GITHUB_APP_PRIVATE_KEY
# - GITHUB_APP_CLIENT_ID: JWT `iss` (GitHub recommends the app Client ID for `iss`; see GitHub docs).
# - GITHUB_APP_INSTALLATION_ID: required for POST /app/installations/{id}/access_tokens (not interchangeable with Client ID).
# GITHUB_APP_PRIVATE_KEY: PEM string, or path when the value points to an existing file (GitLab File variables).
# GitLab often stores multiline PEM as one line with literal \n, or with CRLF — normalize before signing.
set -euo pipefail

: "${GITHUB_APP_CLIENT_ID:?GITHUB_APP_CLIENT_ID is required (GitHub App Client ID for JWT iss)}"
: "${GITHUB_APP_INSTALLATION_ID:?GITHUB_APP_INSTALLATION_ID is required}"
: "${GITHUB_APP_PRIVATE_KEY:?GITHUB_APP_PRIVATE_KEY is required}"

key_file=$(mktemp)
cleanup() { rm -f "$key_file"; }
trap cleanup EXIT

pem_valid() {
	openssl pkey -in "$key_file" -noout 2>/dev/null || openssl rsa -in "$key_file" -noout 2>/dev/null
}

normalize_github_app_private_key() {
	local raw="${GITHUB_APP_PRIVATE_KEY}"

	if [ -f "$raw" ]; then
		tr -d '\r' <"$raw" >"$key_file"
	else
		# Multiline-in-one-line: backslash-n sequences (GitLab UI / env export)
		printf '%b' "$raw" | tr -d '\r' >"$key_file"
		if ! pem_valid; then
			printf '%s' "$raw" | tr -d '\r' | sed 's/\\n/\n/g' >"$key_file"
		fi
	fi

	chmod 600 "$key_file"

	if ! pem_valid; then
		echo "github-installation-token: private key is not valid PEM after normalization." >&2
		echo "Use a GitLab File variable for the .pem, or a single-line value with \\n between PEM lines; avoid CRLF from Windows." >&2
		if head -n1 "$key_file" | grep -q 'BEGIN'; then
			echo "First line looks like: $(head -n1 "$key_file")" >&2
		else
			echo "First 24 bytes (hex): $(head -c 24 "$key_file" | od -An -tx1 | tr -s ' ')" >&2
		fi
		exit 1
	fi
}

normalize_github_app_private_key

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
