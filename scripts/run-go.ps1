param(
  [Parameter(Mandatory = $true)]
  [string]$ProjectRoot,

  [Parameter(Mandatory = $true)]
  [string]$DatabaseUrl,

  [Parameter(Mandatory = $true)]
  [string]$JwtSecret
)

$ErrorActionPreference = "Stop"

Set-Location (Join-Path $ProjectRoot "backend\go_post_service")
$env:DATABASE_URL = $DatabaseUrl
$env:READ_DATABASE_URL = $DatabaseUrl
$env:JWT_SECRET_KEY = $JwtSecret
$env:GO_POST_SERVICE_ADDR = ":8080"

go run ./cmd/server
