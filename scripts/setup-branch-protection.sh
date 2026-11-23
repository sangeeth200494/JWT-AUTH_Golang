#!/usr/bin/env bash
#
# Configures GitHub Branch Protection on the 'main' branch of the repository.
# Requires GITHUB_TOKEN environment variable or $1 parameter with repo administration permissions.
#

set -e

TOKEN="${1:-${GITHUB_TOKEN}}"
OWNER="${2:-sangeeth200494}"
REPO="${3:-JWT-AUTH_Golang}"
BRANCH="main"

if [ -z "$TOKEN" ]; then
  echo "Error: GitHub Personal Access Token is required."
  echo "Usage: ./scripts/setup-branch-protection.sh <GITHUB_TOKEN> [OWNER] [REPO]"
  echo "Or set GITHUB_TOKEN environment variable."
  exit 1
fi

API_URL="https://api.github.com/repos/${OWNER}/${REPO}/branches/${BRANCH}/protection"

PAYLOAD='{
  "required_status_checks": {
    "strict": true,
    "contexts": [
      "Quality Gate & Validation"
    ]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "dismiss_stale_reviews": true,
    "require_code_owner_reviews": false,
    "required_approving_review_count": 1,
    "require_last_push_approval": true
  },
  "restrictions": null,
  "required_linear_history": false,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "block_creations": false,
  "required_conversation_resolution": true
}'

echo "Applying branch protection to ${OWNER}/${REPO} on branch '${BRANCH}'..."

curl -s -X PUT \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "${API_URL}" \
  -d "${PAYLOAD}" | grep -q '"url"' && echo "Branch protection successfully applied to 'main'!" || {
    echo "Failed to apply branch protection. Please check your token permissions."
    exit 1
  }
