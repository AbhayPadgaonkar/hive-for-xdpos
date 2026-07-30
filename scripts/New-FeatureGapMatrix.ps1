param(
    [string[]]$ClientNames = @("xdc-geth-audit"),
    [string]$OutputDir = "$PSScriptRoot\..\error_ledger",
    [string]$HiveExe = ".\hive.exe",
    [switch]$RunProbe,
    [switch]$CompareOnly
)

$ErrorActionPreference = "Stop"

function Write-JsonFile {
    param([string]$Path, $InputObject, [int]$Depth=4)
    $json = $InputObject | ConvertTo-Json -Depth $Depth
    [System.IO.File]::WriteAllText($Path, $json, [System.Text.UTF8Encoding]::new($false))
}

function Invoke-Probe {
    param([string]$Client)
    Write-Host "=== Probing $Client ==="
    $logFile = Join-Path $env:TEMP "probe-$Client-$(Get-Date -Format 'yyyyMMdd-HHmmss').log"
    $output = & $HiveExe "--sim", "probe/rpc-methods", "--client", $Client 2>&1
    $output | Out-File -FilePath $logFile -Encoding utf8

    $simDir = Join-Path $PSScriptRoot "..\workspace\logs"
    $simLog = Get-ChildItem -Path $simDir -Filter "*simulator*" -Name | Sort-Object -Descending | Select-Object -First 1
    if (-not $simLog) {
        Write-Warning "No simulator log found for $Client"
        return $null
    }
    $simLogPath = Join-Path $simDir $simLog
    $content = Get-Content -Path $simLogPath -Raw -Encoding UTF8 2>$null
    if (-not $content) { return $null }

    $probeIdx = $content.IndexOf('PROBE_RESULT')
    if ($probeIdx -lt 0) {
        Write-Warning "No PROBE_RESULT found in simulator log for $Client"
        return $null
    }
    $jsonPart = $content.Substring($probeIdx + 12)
    $endIdx = $jsonPart.LastIndexOf('}')
    if ($endIdx -lt 0) {
        Write-Warning "No closing brace in PROBE_RESULT for $Client"
        return $null
    }
    $probeJson = $jsonPart.Substring(0, $endIdx + 1)
    $result = $probeJson | ConvertFrom-Json
    $resultPath = Join-Path $OutputDir "probe-$Client-$(Get-Date -Format 'yyyyMMdd-HHmmss').json"
    Write-JsonFile -Path $resultPath -InputObject $result -Depth 4
    Write-Host "Probe result saved: $resultPath"
    return $result
}

function Get-StoredProbe {
    param([string]$Client, [string]$Dir)
    $pattern = "probe-$Client-*.json"
    $found = Get-ChildItem -Path $Dir -Filter $pattern | Sort-Object LastWriteTime -Descending | Select-Object -First 1
    if (-not $found) {
        Write-Warning "No probe result found for $Client in $Dir"
        return $null
    }
    Write-Host "Loaded probe: $($found.Name)"
    return Get-Content -Path $found.FullName -Raw -Encoding utf8 | ConvertFrom-Json
}

function New-Matrix {
    param($ProbeA, $ProbeB)
    if (-not $ProbeA -or -not $ProbeB) {
        Write-Error "Need two probe results to compare"
        return
    }

    $methodsA = @{}; $ProbeA.methods | ForEach-Object { $methodsA[$_.method] = $_ }
    $methodsB = @{}; $ProbeB.methods | ForEach-Object { $methodsB[$_.method] = $_ }

    $allMethods = ($methodsA.Keys + $methodsB.Keys) | Sort-Object -Unique

    $rows = @()
    $inANotB = @()
    $inBNotA = @()
    $bothSupported = 0
    $bothUnsupported = 0

    foreach ($m in $allMethods) {
        $ma = $methodsA[$m]
        $mb = $methodsB[$m]
        $aSupported = $ma -and $ma.supported
        $bSupported = $mb -and $mb.supported

        $row = [PSCustomObject]@{
            method = $m
            a_supported = if ($aSupported) { $true } else { $false }
            b_supported = if ($bSupported) { $true } else { $false }
            a_error = if ($ma -and $ma.error) { $ma.error } else { "" }
            b_error = if ($mb -and $mb.error) { $mb.error } else { "" }
        }
        $rows += $row

        if ($aSupported -and -not $bSupported) { $inANotB += $m }
        if (-not $aSupported -and $bSupported) { $inBNotA += $m }
        if ($aSupported -and $bSupported) { $bothSupported++ }
        if (-not $aSupported -and -not $bSupported) { $bothUnsupported++ }
    }

    $summary = [PSCustomObject]@{
        date           = (Get-Date).ToString('yyyy-MM-dd')
        client_a       = $ProbeA.client
        client_b       = $ProbeB.client
        version_a      = $ProbeA.version
        version_b      = $ProbeB.version
        modules_a      = $ProbeA.modules
        modules_b      = $ProbeB.modules
        total_methods  = $allMethods.Count
        both_supported = $bothSupported
        both_unsupported = $bothUnsupported
        in_a_not_b     = $inANotB.Count
        in_b_not_a     = $inBNotA.Count
        gaps_a_not_b   = $inANotB
        gaps_b_not_a   = $inBNotA
        matrix         = $rows
    }

    return $summary
}

function Show-Matrix {
    param($Matrix)
    Write-Host "`n============================================"
    Write-Host "FEATURE GAP MATRIX"
    Write-Host "  A: $($Matrix.client_a) v$($Matrix.version_a)"
    Write-Host "  B: $($Matrix.client_b) v$($Matrix.version_b)"
    Write-Host "============================================"
    Write-Host "  Total methods checked : $($Matrix.total_methods)"
    Write-Host "  Both supported        : $($Matrix.both_supported)"
    Write-Host "  Both unsupported      : $($Matrix.both_unsupported)"
    Write-Host "  In A not B (gaps)     : $($Matrix.in_a_not_b)"
    Write-Host "  In B not A (extras)   : $($Matrix.in_b_not_a)"

    Write-Host "`n--- Features in A (gaps to port to B) ---"
    $Matrix.gaps_a_not_b | ForEach-Object { Write-Host "  GAP: $_" }

    Write-Host "`n--- Features in B (extras vs A) ---"
    $Matrix.gaps_b_not_a | ForEach-Object { Write-Host "  EXTRA: $_" }
}

function Save-Matrix {
    param($Matrix, [string]$Dir)
    $aName = $Matrix.client_a
    $bName = $Matrix.client_b
    $fileName = "gapmatrix-$aName-vs-$bName-$(Get-Date -Format 'yyyyMMdd').json"
    $path = Join-Path $Dir $fileName
    Write-JsonFile -Path $path -InputObject $Matrix -Depth 5
    Write-Host "Gap matrix saved: $path"
}

# --- Main ---

$probes = @()
foreach ($client in $ClientNames) {
    if ($RunProbe) {
        $p = Invoke-Probe -Client $client
        if ($p) { $probes += $p }
    }
    else {
        $p = Get-StoredProbe -Client $client -Dir $OutputDir
        if ($p) { $probes += $p }
    }
}

if ($probes.Count -ge 2) {
    $matrix = New-Matrix -ProbeA $probes[0] -ProbeB $probes[1]
    Show-Matrix -Matrix $matrix
    Save-Matrix -Matrix $matrix -Dir $OutputDir
}
else {
    Write-Warning "Need probe results for at least 2 clients. Use -RunProbe or ensure probe JSONs exist in $OutputDir"
}
