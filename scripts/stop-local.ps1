$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$projectRoot = Split-Path $PSScriptRoot -Parent
$runDir = Join-Path $projectRoot ".run"

function Stop-TrackedProcess {
  param([string]$Name)

  $pidFile = Join-Path $runDir "$Name.pid"
  if (-not (Test-Path $pidFile)) {
    return
  }

  $rawPid = (Get-Content $pidFile | Select-Object -First 1).Trim()
  if ($rawPid -match "^\d+$") {
    $process = Get-Process -Id ([int]$rawPid) -ErrorAction SilentlyContinue
    if ($process) {
      & taskkill.exe /PID $process.Id /T /F | Out-Null
    }
  }

  Remove-Item $pidFile -Force
}

foreach ($name in "vite", "go", "flask") {
  Stop-TrackedProcess -Name $name
}

Start-Sleep -Seconds 2

Write-Host "Stopped tracked local services." -ForegroundColor Green
$remaining = netstat -ano | Select-String "LISTENING" | Select-String ":5000|:8080|:5173"
if ($remaining) {
  Write-Host "Ports still in use:" -ForegroundColor Yellow
  $remaining
}
