param(
    [string]$Direction = "up",
    [int]$Steps = 0
)

$ScriptDir  = Split-Path -Parent $MyInvocation.MyCommand.Path
$BackendDir = Join-Path $ScriptDir "..\backend"

if (-not $env:DATABASE_URL) {
    $env:DATABASE_URL = "postgres://messenger:messenger_secret@localhost:5432/messenger?sslmode=disable"
}
if (-not $env:MIGRATIONS_PATH) {
    $env:MIGRATIONS_PATH = "migrations"
}

Write-Host "Database : $($env:DATABASE_URL)"
Write-Host "Direction: $Direction"
if ($Direction -eq "down" -and $Steps -gt 0) {
    Write-Host "Steps    : $Steps"
}
Write-Host ""

Push-Location $BackendDir
try {
    go run ./cmd/migrate -direction="$Direction" -steps="$Steps"
} finally {
    Pop-Location
}
