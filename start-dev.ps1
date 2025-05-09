# start-dev.ps1
# Usage: .\start-dev.ps1
# Starts Go API and React dev server with one command

# Resolve script directory & go there
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition
Set-Location $scriptDir

# Define frontend dir
$frontendDir = Join-Path $scriptDir 'luno-trading-dashboard-pro-main'

# Load environment variables from .env
if (Test-Path .\.env) {
  Get-Content .\.env | ForEach-Object {
    if ($_ -match '^\s*([^#][^=]+)=([^#]+)') {
      $name = $Matches[1].Trim()
      $value = $Matches[2].Trim()
      [System.Environment]::SetEnvironmentVariable($name, $value, 'Process')
    }
  }
}

# Ensure Go dependencies
Write-Host "Running go mod tidy..."
go mod tidy

# Install frontend deps if missing
if (-not (Test-Path (Join-Path $frontendDir 'node_modules'))) {
  Write-Host "Installing frontend dependencies..."
  Push-Location $frontendDir
  npm install
  Pop-Location
}

# Start Go API in background
Write-Host "Starting Go API..."
Start-Process -NoNewWindow go -ArgumentList 'run cmd/bot/main.go --config=config/config.json'

# Start React dev server
Write-Host "Starting React dev server..."
Push-Location $frontendDir
npm run dev
