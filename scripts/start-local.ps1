param(
  [switch]$OpenBrowser
)

$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$projectRoot = Split-Path $PSScriptRoot -Parent
$backendRoot = Join-Path $projectRoot "backend"
$runDir = Join-Path $projectRoot ".run"
$logDir = Join-Path $projectRoot ".runlogs"

$powershellExe = "C:\WINDOWS\System32\WindowsPowerShell\v1.0\powershell.exe"
$mysqlFallback = "C:\Program Files\MySQL\MySQL Server 8.0\bin\mysql.exe"

function Write-Step {
  param([string]$Message)
  Write-Host "==> $Message" -ForegroundColor Cyan
}

function Get-ToolPath {
  param(
    [string]$Name,
    [string[]]$Fallbacks = @()
  )

  $command = Get-Command $Name -ErrorAction SilentlyContinue
  if ($command) {
    return $command.Source
  }

  foreach ($fallback in $Fallbacks) {
    if (Test-Path $fallback) {
      return $fallback
    }
  }

  throw "Missing required tool: $Name"
}

function Test-PortListening {
  param([int]$Port)

  $matches = netstat -ano | Select-String ":$Port"
  return [bool]($matches | Where-Object { $_.Line -match "LISTENING" })
}

function Test-ServiceReady {
  param(
    [int]$Port,
    [string]$Name,
    [string]$Url
  )

  if (-not (Test-PortListening -Port $Port)) {
    return $false
  }

  try {
    $response = Invoke-WebRequest -UseBasicParsing -Uri $Url -TimeoutSec 3
    if ($response.StatusCode -ge 200 -and $response.StatusCode -lt 500) {
      Write-Step "$Name is already running on port $Port"
      return $true
    }
  } catch {
  }

  throw "$Name port $Port is already in use, but the health check failed. Stop the existing process before retrying."
}

function Wait-HttpReady {
  param(
    [string]$Name,
    [string]$Url,
    [int]$Attempts = 30
  )

  for ($i = 0; $i -lt $Attempts; $i++) {
    try {
      $response = Invoke-WebRequest -UseBasicParsing -Uri $Url -TimeoutSec 3
      if ($response.StatusCode -ge 200 -and $response.StatusCode -lt 500) {
        return
      }
    } catch {
    }

    Start-Sleep -Seconds 1
  }

  throw "$Name failed to become ready at $Url"
}

function Ensure-NodeDependencies {
  $nodeModules = Join-Path $projectRoot "node_modules"
  if (-not (Test-Path $nodeModules)) {
    Write-Step "Installing frontend dependencies"
    & npm.cmd install
    if ($LASTEXITCODE -ne 0) {
      throw "npm install failed"
    }
  }
}

function Ensure-PythonDependencies {
  & python -c "import flask, flask_sqlalchemy, flask_jwt_extended, pymysql" 2>$null
  if ($LASTEXITCODE -ne 0) {
    Write-Step "Installing Python dependencies"
    & python -m pip install --user -r (Join-Path $backendRoot "requirements.txt")
    if ($LASTEXITCODE -ne 0) {
      throw "Python dependency installation failed"
    }
  }
}

function Ensure-Database {
  param(
    [string]$MySqlExe,
    [string]$DatabaseUser,
    [string]$DatabasePassword,
    [string]$DatabaseHost,
    [string]$DatabasePort,
    [string]$DatabaseName
  )

  Write-Step "Ensuring MySQL database exists"
  & $MySqlExe "-u$DatabaseUser" "-p$DatabasePassword" "-h$DatabaseHost" "-P$DatabasePort" -e "CREATE DATABASE IF NOT EXISTS $DatabaseName CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
  if ($LASTEXITCODE -ne 0) {
    throw "Failed to create database $DatabaseName"
  }

  Write-Step "Initializing Flask tables and seed data"
  $env:DATABASE_URL = $script:flaskDatabaseUrl
  $env:JWT_SECRET_KEY = $script:jwtSecret
  & python (Join-Path $backendRoot "init_db.py")
  if ($LASTEXITCODE -ne 0) {
    throw "backend/init_db.py failed"
  }

  Write-Step "Ensuring community tables exist"
  Get-Content -Raw (Join-Path $PSScriptRoot "setup-community-schema.sql") |
    & $MySqlExe "-u$DatabaseUser" "-p$DatabasePassword" "-h$DatabaseHost" "-P$DatabasePort" $DatabaseName
  if ($LASTEXITCODE -ne 0) {
    throw "Community schema bootstrap failed"
  }
}

function Start-ServiceProcess {
  param(
    [string]$Name,
    [string]$ScriptName,
    [string[]]$Arguments,
    [string]$ReadyUrl
  )

  $stdoutLog = Join-Path $logDir "$Name.out.log"
  $stderrLog = Join-Path $logDir "$Name.err.log"
  $pidFile = Join-Path $runDir "$Name.pid"

  if (Test-Path $stdoutLog) {
    Remove-Item $stdoutLog -Force
  }
  if (Test-Path $stderrLog) {
    Remove-Item $stderrLog -Force
  }

  $argList = @(
    "-NoProfile",
    "-ExecutionPolicy",
    "Bypass",
    "-File",
    (Join-Path $PSScriptRoot $ScriptName)
  ) + $Arguments

  $process = Start-Process -FilePath $powershellExe `
    -ArgumentList $argList `
    -WorkingDirectory $projectRoot `
    -WindowStyle Minimized `
    -RedirectStandardOutput $stdoutLog `
    -RedirectStandardError $stderrLog `
    -PassThru

  Set-Content -Path $pidFile -Value $process.Id
  Wait-HttpReady -Name $Name -Url $ReadyUrl
}

$null = Get-ToolPath -Name "python.exe"
$null = Get-ToolPath -Name "go.exe"
$null = Get-ToolPath -Name "npm.cmd"
$mysqlExe = if ($env:MYSQL_EXE) { $env:MYSQL_EXE } else { Get-ToolPath -Name "mysql.exe" -Fallbacks @($mysqlFallback) }

$databaseUser = if ($env:HX_MYSQL_USER) { $env:HX_MYSQL_USER } else { "root" }
$databasePassword = if ($env:HX_MYSQL_PASSWORD) { $env:HX_MYSQL_PASSWORD } else { "hutbhutb0000" }
$databaseHost = if ($env:HX_MYSQL_HOST) { $env:HX_MYSQL_HOST } else { "127.0.0.1" }
$databasePort = if ($env:HX_MYSQL_PORT) { $env:HX_MYSQL_PORT } else { "3306" }
$databaseName = if ($env:HX_MYSQL_DATABASE) { $env:HX_MYSQL_DATABASE } else { "huxiang_culture" }
$script:jwtSecret = if ($env:JWT_SECRET_KEY) { $env:JWT_SECRET_KEY } else { "jwt-huxiang-secret-key-dev" }
$script:flaskDatabaseUrl = if ($env:DATABASE_URL) {
  $env:DATABASE_URL
} else {
  "mysql+pymysql://{0}:{1}@{2}:{3}/{4}?charset=utf8mb4" -f $databaseUser, $databasePassword, $databaseHost, $databasePort, $databaseName
}
$goDatabaseUrl = if ($env:GO_DATABASE_URL) {
  $env:GO_DATABASE_URL
} else {
  "{0}:{1}@tcp({2}:{3})/{4}?charset=utf8mb4" -f $databaseUser, $databasePassword, $databaseHost, $databasePort, $databaseName
}

New-Item -ItemType Directory -Force -Path $runDir, $logDir | Out-Null

Ensure-NodeDependencies
Ensure-PythonDependencies
Ensure-Database -MySqlExe $mysqlExe `
  -DatabaseUser $databaseUser `
  -DatabasePassword $databasePassword `
  -DatabaseHost $databaseHost `
  -DatabasePort $databasePort `
  -DatabaseName $databaseName

$flaskRunning = Test-ServiceReady -Port 5000 -Name "Flask" -Url "http://127.0.0.1:5000/health"
$goRunning = Test-ServiceReady -Port 8080 -Name "Go community service" -Url "http://127.0.0.1:8080/health"
$viteRunning = Test-ServiceReady -Port 5173 -Name "Vite" -Url "http://127.0.0.1:5173"

if (-not $flaskRunning) {
  Write-Step "Starting Flask API on 5000"
  Start-ServiceProcess -Name "flask" `
    -ScriptName "run-flask.ps1" `
    -Arguments @("-ProjectRoot", $projectRoot, "-DatabaseUrl", $script:flaskDatabaseUrl, "-JwtSecret", $script:jwtSecret) `
    -ReadyUrl "http://127.0.0.1:5000/health"
}

if (-not $goRunning) {
  Write-Step "Starting Go community API on 8080"
  Start-ServiceProcess -Name "go" `
    -ScriptName "run-go.ps1" `
    -Arguments @("-ProjectRoot", $projectRoot, "-DatabaseUrl", $goDatabaseUrl, "-JwtSecret", $script:jwtSecret) `
    -ReadyUrl "http://127.0.0.1:8080/health"
}

if (-not $viteRunning) {
  Write-Step "Starting Vite dev server on 5173"
  Start-ServiceProcess -Name "vite" `
    -ScriptName "run-vite.ps1" `
    -Arguments @("-ProjectRoot", $projectRoot) `
    -ReadyUrl "http://127.0.0.1:5173"
}

if ($OpenBrowser) {
  Start-Process "http://127.0.0.1:5173/"
  Start-Process "http://127.0.0.1:5173/login"
  Start-Process "http://127.0.0.1:5173/community"
}

Write-Host ""
Write-Host "Local stack is ready." -ForegroundColor Green
Write-Host "Frontend: http://127.0.0.1:5173"
Write-Host "Flask:    http://127.0.0.1:5000"
Write-Host "Go API:   http://127.0.0.1:8080"
Write-Host "Admin:    admin / admin123"
