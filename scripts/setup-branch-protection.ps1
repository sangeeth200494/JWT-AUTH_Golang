<#
.SYNOPSIS
    Configures GitHub Branch Protection on the 'main' branch of the repository.

.DESCRIPTION
    Enforces quality gates before any code can be merged into 'main':
    - Requires Pull Request reviews before merging.
    - Requires passing CI status checks ("Quality Gate & Validation").
    - Requires branch to be up-to-date before merging.
    - Enforces rules on administrators.
    - Prevents force pushes and branch deletion.

.PARAMETER Token
    GitHub Personal Access Token with repo administration permissions.
    Can also be provided via the GITHUB_TOKEN environment variable.

.PARAMETER Owner
    GitHub repository owner (default: sangeeth200494).

.PARAMETER Repo
    GitHub repository name (default: JWT-AUTH_Golang).
#>

param (
    [Parameter(Mandatory = $false)]
    [string]$Token = $env:GITHUB_TOKEN,

    [Parameter(Mandatory = $false)]
    [string]$Owner = "sangeeth200494",

    [Parameter(Mandatory = $false)]
    [string]$Repo = "JWT-AUTH_Golang"
)

if (-not $Token) {
    Write-Error "Error: GitHub Personal Access Token is required. Set `$env:GITHUB_TOKEN or pass -Token <PAT>"
    exit 1
}

$headers = @{
    "Authorization" = "Bearer $Token"
    "Accept"        = "application/vnd.github+json"
    "X-GitHub-Api-Version" = "2022-11-28"
}

$body = @{
    required_status_checks = @{
        strict = $true
        contexts = @("Quality Gate & Validation")
    }
    enforce_admins = $true
    required_pull_request_reviews = @{
        dismiss_stale_reviews = $true
        require_code_owner_reviews = $false
        required_approving_review_count = 1
        require_last_push_approval = $true
    }
    restrictions = $null
    required_linear_history = $false
    allow_force_pushes = $false
    allow_deletions = $false
    block_creations = $false
    required_conversation_resolution = $true
} | ConvertTo-Json -Depth 5

$url = "https://api.github.com/repos/$Owner/$Repo/branches/main/protection"

Write-Host "Configuring branch protection on $Owner/$Repo (branch: main)..."

try {
    $response = Invoke-RestMethod -Uri $url -Method Put -Headers $headers -Body $body -ContentType "application/json"
    Write-Host "Success! Branch protection rules successfully applied to 'main':" -ForegroundColor Green
    Write-Host "  - Required status check: 'Quality Gate & Validation'"
    Write-Host "  - Require branches to be up-to-date: YES"
    Write-Host "  - Required approving reviews: 1"
    Write-Host "  - Dismiss stale reviews on new commits: YES"
    Write-Host "  - Require conversation resolution: YES"
    Write-Host "  - Enforce rules on administrators: YES"
    Write-Host "  - Prevent force pushes: YES"
    Write-Host "  - Prevent branch deletion: YES"
} catch {
    Write-Error "Failed to apply branch protection: $_"
    exit 1
}
