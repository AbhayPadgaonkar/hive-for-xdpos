function Write-JsonFile {
    param([string]$Path, $InputObject, [int]$Depth=4)
    $json = $InputObject | ConvertTo-Json -Depth $Depth
    [System.IO.File]::WriteAllText($Path, $json, [System.Text.UTF8Encoding]::new($false))
}

function Invoke-HiveRun {
    param(
        [Parameter(Mandatory)]
        [string]$Client,
        [Parameter(Mandatory)]
        [string]$Simulator,
        [string]$HiveExe = ".\hive.exe",
        [string]$HiveRoot = "C:\BlocksScan\hive",
        [string]$OutputDir = "$HiveRoot\error_ledger",
        [string]$Limit = "",
        [switch]$PassThru
    )
    $ErrorActionPreference = "Stop"
    Push-Location $HiveRoot
    try {
        $simName = $Simulator -replace '/', '-'
        $dateStamp = Get-Date -Format 'yyyyMMdd-HHmmss'
        $outFile = Join-Path $OutputDir "$Client-$simName-$dateStamp.json"
        $logFile = Join-Path $env:TEMP "hive-$Client-$simName-$dateStamp.log"

        Write-Host "=== Running hive --sim $Simulator --client $Client ==="

        $argsList = @("--sim", $Simulator, "--client", $Client)
        if ($Limit) { $argsList += @("--sim.limit", $Limit) }

        $output = & $HiveExe $argsList 2>&1
        $output | Out-File -FilePath $logFile -Encoding utf8

        $result = [PSCustomObject]@{
            client      = $Client
            simulator   = $Simulator
            date        = (Get-Date).ToString('yyyy-MM-dd')
            timestamp   = (Get-Date).ToString('o')
            command     = "hive --sim $Simulator --client $Client"
            version     = ""
            suites      = 0
            total       = 0
            passed      = 0
            failed      = 0
            tests       = @()
            raw_log     = $logFile
        }

        $tests = @{}
        foreach ($line in $output) {
            $text = "$line"
            if ($text -match 'test started.*test=(\d+).*name="(.+?)"') {
                $tid = [int]$matches[1]
                $tname = $matches[2]
                $tests[$tid] = @{ name = $tname; pass = $null }
            }
            elseif ($text -match 'test ended.*test=(\d+).*pass=(true|false)') {
                $tid = [int]$matches[1]
                $pass = $matches[2] -eq 'true'
                if ($tests.ContainsKey($tid)) {
                    $tests[$tid].pass = $pass
                }
            }
            elseif ($text -match 'simulation .+ finished suites=(\d+) tests=(\d+) failed=(\d+)') {
                $result.suites = [int]$matches[1]
                $result.total = [int]$matches[2]
                $result.failed = [int]$matches[3]
                $result.passed = $result.total - $result.failed
            }
            elseif ($text -match 'client version.*result":"(.+?)"') {
                $result.version = $matches[1]
            }
        }

        $result.tests = $tests.Keys | Sort-Object | ForEach-Object {
            $t = $tests[$_]
            [PSCustomObject]@{ id = $_; name = $t.name; pass = $t.pass }
        }

        Write-JsonFile -Path $outFile -InputObject $result -Depth 4
        Write-Host "Saved: $outFile"
        Write-Host "Result: $($result.passed)/$($result.total) passed, $($result.failed) failed"

        if ($PassThru) { return $result }
    }
    finally { Pop-Location }
}

function Get-HiveProbeResult {
    param(
        [string]$LogsDir = ".\workspace\logs",
        [string]$SimPattern = "*simulator*",
        [string]$OutDir = ".\error_ledger",
        [string]$ClientLabel = ""
    )
    $simLog = Get-ChildItem -Path $LogsDir -Filter $SimPattern -Name | Sort-Object -Descending | Select-Object -First 1
    if (-not $simLog) { Write-Warning "No simulator logs found in $LogsDir"; return $null }
    $fullPath = Join-Path $LogsDir $simLog
    $content = Get-Content -Path $fullPath -Raw -Encoding UTF8 2>$null
    if (-not $content) { Write-Warning "Could not read $fullPath"; return $null }
    $idx = $content.IndexOf('PROBE_RESULT')
    if ($idx -lt 0) { Write-Warning "PROBE_RESULT not found in $simLog"; return $null }
    $jsonPart = $content.Substring($idx + 12)
    $endIdx = $jsonPart.LastIndexOf('}')
    if ($endIdx -lt 0) { Write-Warning "No closing brace in PROBE_RESULT"; return $null }
    $jsonPart = $jsonPart.Substring(0, $endIdx + 1)
    $result = $jsonPart | ConvertFrom-Json
    if ($OutDir) {
        if (-not (Test-Path $OutDir)) { New-Item -ItemType Directory -Path $OutDir -Force | Out-Null }
        $label = if ($ClientLabel) { $ClientLabel } else { $result.client }
        $stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
        $outFile = Join-Path $OutDir "probe-$label-$stamp.json"
        Write-JsonFile -Path $outFile -InputObject $result -Depth 3
        Write-Host "Saved probe result: $outFile"
    }
    return $result
}

Export-ModuleMember -Function Invoke-HiveRun, Get-HiveProbeResult
