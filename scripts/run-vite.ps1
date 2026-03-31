param(
  [Parameter(Mandatory = $true)]
  [string]$ProjectRoot
)

$ErrorActionPreference = "Stop"

Set-Location $ProjectRoot
$env:VITE_FLASK_API_URL = "http://127.0.0.1:5000"
$env:VITE_COMMUNITY_API_URL = "http://127.0.0.1:8080"

npm.cmd run dev -- --host 127.0.0.1 --port 5173
