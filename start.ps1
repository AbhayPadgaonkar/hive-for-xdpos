param(
    [switch]$BackendOnly,
    [switch]$FrontendOnly,
    [switch]$Import,
    [switch]$Build,
    [string]$ApiPort = "8080",
    [string]$FrontendPort = "3000"
)

$Root = $PSScriptRoot

function Start-Backend {
    Write-Host "Starting API server on port $ApiPort..."
    $api = Start-Process -FilePath "$Root\backend\hive-api.exe" -WorkingDirectory "$Root\backend" -PassThru
    Start-Sleep -Seconds 2
    $ok = $false
    try {
        $r = Invoke-WebRequest -Uri "http://localhost:$ApiPort/api/health" -UseBasicParsing -ErrorAction Stop
        $ok = $true
    } catch {}
    if (-not $ok) { Write-Warning "Backend may not have started correctly" }
    else { Write-Host "Backend OK (PID: $($api.Id))" }
    return $api
}

function Start-Frontend {
    Write-Host "Starting frontend on port $FrontendPort..."
    if ($Build -or -not (Test-Path "$Root\frontend\.next\BUILD_ID")) {
        Write-Host "Building frontend..." 
        Set-Location "$Root\frontend"
        npx next build --webpack 2>&1 | Out-Null
        Set-Location $Root
    }
    $env:NEXT_PUBLIC_API_URL = "http://localhost:$ApiPort"
    $fe = Start-Process -FilePath "node" -WorkingDirectory "$Root\frontend" -ArgumentList "node_modules/next/dist/bin/next", "start", "-p", $FrontendPort -PassThru
    Start-Sleep -Seconds 8
    $ok = $false
    try {
        $r = Invoke-WebRequest -Uri "http://localhost:$FrontendPort" -UseBasicParsing -ErrorAction Stop
        $ok = $true
    } catch {}
    if (-not $ok) { Write-Warning "Frontend may not have started correctly" }
    else { Write-Host "Frontend OK (PID: $($fe.Id))" }
    return $fe
}

if ($FrontendOnly) {
    $fe = Start-Frontend
} elseif ($BackendOnly) {
    $be = Start-Backend
    if ($Import) {
        Start-Sleep -Seconds 1
        Invoke-RestMethod -Uri "http://localhost:$ApiPort/api/import" -Method Post -ContentType "application/json" -Body '{}' | Out-Null
    }
} else {
    $be = Start-Backend
    if ($Import) {
        Start-Sleep -Seconds 1
        Invoke-RestMethod -Uri "http://localhost:$ApiPort/api/import" -Method Post -ContentType "application/json" -Body '{}' | Out-Null
    }
    $fe = Start-Frontend
}

Write-Host @"

=====================================
  Hive Dashboard
=====================================
  API:  http://localhost:$ApiPort
  UI:   http://localhost:$FrontendPort

  Quick links:
    http://localhost:$FrontendPort          — Dashboard (runs)
    http://localhost:$FrontendPort/gap-matrix  — Feature gap matrix
    http://localhost:$FrontendPort/comparisons — Comparisons

  Press any key to stop all services.
=====================================
"@

if (-not $BackendOnly -and -not $FrontendOnly) {
    $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
    Stop-Process -Name "hive-api" -Force -ErrorAction SilentlyContinue
    Stop-Process -Name "node" -Force -ErrorAction SilentlyContinue
}
