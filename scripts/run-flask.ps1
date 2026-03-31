param(
  [Parameter(Mandatory = $true)]
  [string]$ProjectRoot,

  [Parameter(Mandatory = $true)]
  [string]$DatabaseUrl,

  [Parameter(Mandatory = $true)]
  [string]$JwtSecret
)

$ErrorActionPreference = "Stop"

Set-Location (Join-Path $ProjectRoot "backend")
$env:DATABASE_URL = $DatabaseUrl
$env:JWT_SECRET_KEY = $JwtSecret

python app.py
