param(
    [switch]$Import,
    [switch]$Rebuild,
    [string]$Port = "8080",
    [string]$DbPath = "$PSScriptRoot\hive.db"
)

$BackendDir = Join-Path $PSScriptRoot "backend"
$ApiExe = Join-Path $BackendDir "hive-api.exe"

if ($Rebuild -or -not (Test-Path $ApiExe)) {
    Write-Host "Building API server..."
    Push-Location $BackendDir
    try {
        go build -o hive-api.exe . 2>&1
        if (-not $?) { Write-Host "Build failed"; exit 1 }
    } finally { Pop-Location }
}

$env:HIVE_DB_PATH = $DbPath
$env:HIVE_API_ADDR = ":$Port"

Write-Host "Starting API server on port $Port (DB: $DbPath)"

$proc = Start-Process -FilePath $ApiExe -WorkingDirectory $BackendDir -NoNewWindow -PassThru
Start-Sleep -Seconds 2

if ($Import) {
    Write-Host "Importing data from error_ledger..."
    $result = Invoke-RestMethod -Uri "http://localhost:$Port/api/import" -Method Post -ContentType "application/json" -Body '{}'
    Write-Host "Imported $($result.imported) files"
    if ($result.errors) { $result.errors | ForEach-Object { Write-Warning $_ } }
}

Write-Host @"

API server running (PID: $($proc.Id))

Endpoints:
  GET  /api/health          — Health check
  GET  /api/runs            — List test runs
  POST /api/runs            — Create a run
  GET  /api/runs/{id}       — Get run details
  GET  /api/probes          — List RPC probes
  POST /api/probes          — Submit probe result
  GET  /api/gap-matrices    — List gap matrices
  POST /api/gap-matrices    — Create gap matrix
  GET  /api/comparisons     — List comparisons
  POST /api/comparisons     — Create comparison
  GET  /api/stats           — Summary statistics
  POST /api/import          — Bulk import from error_ledger/

Press Ctrl+C to stop.
"@

$proc.WaitForExit()
