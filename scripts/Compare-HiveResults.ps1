param(
    [string[]]$Clients = @("xdc-geth-audit"),
    [string]$Simulator = "xdc/rpc-compat",
    [string]$BaselineClient = "xdpos",
    [string]$BaselineDir = "$PSScriptRoot\..\error_ledger",
    [string]$HiveExe = ".\hive.exe",
    [switch]$Run,
    [switch]$CompareOnly,
    [string]$Limit = ""
)

$ErrorActionPreference = "Stop"

function Import-HiveResult {
    param([string]$Path)
    $json = Get-Content -Path $Path -Raw -Encoding utf8 | ConvertFrom-Json
    return $json
}

function Compare-Clients {
    param(
        [Parameter(Mandatory)] $Results,
        [string]$LabelA = "client_a",
        [string]$LabelB = "client_b"
    )
    if ($Results.Count -lt 2) {
        Write-Warning "Need at least 2 results to compare"
        return
    }
    $a = $Results[0]
    $b = $Results[1]

    $aTests = @{}; $a.tests | ForEach-Object { $aTests[$_.name] = $_ }
    $bTests = @{}; $b.tests | ForEach-Object { $bTests[$_.name] = $_ }

    $allNames = ($aTests.Keys + $bTests.Keys) | Sort-Object -Unique
    $comparison = @()
    $bothPass = 0; $aOnly = 0; $bOnly = 0; $aFailBPass = 0; $bFailAPass = 0; $bothFail = 0

    foreach ($name in $allNames) {
        $ta = $aTests[$name]
        $tb = $bTests[$name]
        $aPass = if ($ta) { $ta.pass } else { $null }
        $bPass = if ($tb) { $tb.pass } else { $null }
        if ($ta -and $tb) {
            if ($aPass -and $bPass) { $status = "BOTH_PASS"; $bothPass++ }
            elseif (-not $aPass -and -not $bPass) { $status = "BOTH_FAIL"; $bothFail++ }
            elseif ($aPass -and -not $bPass) { $status = "A_ONLY"; $aOnly++ }
            else { $status = "B_ONLY"; $bOnly++ }
        }
        elseif ($ta) {
            $status = if ($aPass) { "MISSING_IN_B" } else { "FAIL_IN_A_MISSING_B" }
        }
        else {
            $status = if ($bPass) { "MISSING_IN_A" } else { "FAIL_IN_B_MISSING_A" }
        }
        $comparison += [PSCustomObject]@{
            name   = $name
            a_pass = $aPass
            b_pass = $bPass
            status = $status
        }
    }

    $summary = [PSCustomObject]@{
        comparison     = "$($a.client) vs $($b.client)"
        simulator      = $a.simulator
        date           = (Get-Date).ToString('yyyy-MM-dd')
        both_pass      = $bothPass
        a_only         = $aOnly
        b_only         = $bOnly
        both_fail      = $bothFail
        total_matched  = $allNames.Count
        client_a       = [PSCustomObject]@{ name = $a.client; total = $a.total; passed = $a.passed; failed = $a.failed }
        client_b       = [PSCustomObject]@{ name = $b.client; total = $b.total; passed = $b.passed; failed = $b.failed }
        tests          = $comparison
    }
    return $summary
}

function Show-ComparisonReport {
    param([Parameter(Mandatory)] $Summary)
    $a = $Summary.client_a
    $b = $Summary.client_b
    Write-Host "`n========================================"
    Write-Host "COMPARISON: $($a.name) vs $($b.name)"
    Write-Host "Simulator: $($Summary.simulator)"
    Write-Host "========================================"
    Write-Host "  $($a.name): $($a.passed)/$($a.total) passed ($($a.failed) failed)"
    Write-Host "  $($b.name): $($b.passed)/$($b.total) passed ($($b.failed) failed)"
    Write-Host "----------------------------------------"
    Write-Host "  Both pass : $($Summary.both_pass)"
    Write-Host "  A only    : $($Summary.a_only)  (passes on $($a.name) but fails on $($b.name))"
    Write-Host "  B only    : $($Summary.b_only)  (passes on $($b.name) but fails on $($a.name))"
    Write-Host "  Both fail : $($Summary.both_fail)"
    Write-Host "----------------------------------------"

    $regressions = $Summary.tests | Where-Object { $_.status -eq "A_ONLY" -or $_.status -eq "FAIL_IN_A_MISSING_B" }
    if ($regressions) {
        Write-Host "`nREGRESSIONS (pass on $($a.name), fail on $($b.name)):"
        $regressions | ForEach-Object { Write-Host "  FAIL  $($_.name)" }
    }

    $improvements = $Summary.tests | Where-Object { $_.status -eq "B_ONLY" -or $_.status -eq "FAIL_IN_B_MISSING_A" }
    if ($improvements) {
        Write-Host "`nIMPROVEMENTS (fail on $($a.name), pass on $($b.name)):"
        $improvements | ForEach-Object { Write-Host "  PASS  $($_.name)" }
    }

    $bothFailList = $Summary.tests | Where-Object { $_.status -eq "BOTH_FAIL" }
    if ($bothFailList) {
        Write-Host "`nBOTH FAIL:"
        $bothFailList | ForEach-Object { Write-Host "  FAIL  $($_.name)" }
    }
}

function Write-JsonFile {
    param([string]$Path, $InputObject, [int]$Depth=4)
    $json = $InputObject | ConvertTo-Json -Depth $Depth
    [System.IO.File]::WriteAllText($Path, $json, [System.Text.UTF8Encoding]::new($false))
}

function Save-ComparisonReport {
    param(
        [Parameter(Mandatory)] $Summary,
        [string]$OutputDir = "$PSScriptRoot\..\error_ledger"
    )
    $a = $Summary.client_a
    $b = $Summary.client_b
    $dateStamp = Get-Date -Format 'yyyyMMdd'
    $fileName = "compare-$($a.name)-vs-$($b.name)-$dateStamp.json"
    $outPath = Join-Path $OutputDir $fileName
    Write-JsonFile -Path $outPath -InputObject $Summary -Depth 5
    Write-Host "`nComparison saved: $outPath"
}

# --- Main Execution ---

if (-not $CompareOnly) {
    $results = @()
    foreach ($client in $Clients) {
        $limitArg = @{}
        if ($Limit) { $limitArg.Limit = $Limit }
        $r = Invoke-HiveRun -Client $client -Simulator $Simulator -HiveExe $HiveExe @limitArg -PassThru
        $results += $r
    }
    if ($results.Count -ge 2) {
        $summary = Compare-Clients -Results $results -LabelA $results[0].client -LabelB $results[1].client
        Show-ComparisonReport -Summary $summary
        Save-ComparisonReport -Summary $summary
    }
}
else {
    $compareFiles = @()
    foreach ($client in $Clients) {
        $simName = $Simulator -replace '/', '-'
        $pattern = "$client-$simName-*.json"
        $found = Get-ChildItem -Path $BaselineDir -Filter $pattern | Sort-Object LastWriteTime -Descending | Select-Object -First 1
        if ($found) {
            $compareFiles += Import-HiveResult -Path $found.FullName
            Write-Host "Loaded: $($found.FullName)"
        }
        else {
            Write-Warning "No result file found matching '$pattern' in $BaselineDir"
        }
    }
    if ($compareFiles.Count -ge 2) {
        $summary = Compare-Clients -Results $compareFiles
        Show-ComparisonReport -Summary $summary
        Save-ComparisonReport -Summary $summary
    }
}
