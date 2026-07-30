param(
    [ValidateSet("probe", "test", "matrix", "diff", "help")]
    [string]$Command = "help"
)

$HiveRoot = $PSScriptRoot

switch ($Command) {
    "probe" {
        & "$HiveRoot\scripts\New-FeatureGapMatrix.ps1" -RunProbe -ClientNames "xdc-geth-audit","go-ethereum"
    }
    "test" {
        & "$HiveRoot\scripts\Compare-HiveResults.ps1" -Clients "xdpos","xdc-geth-audit" -Simulator "xdc/rpc-compat"
    }
    "matrix" {
        & "$HiveRoot\scripts\New-FeatureGapMatrix.ps1" -CompareOnly -ClientNames "xdc-geth-audit","go-ethereum"
    }
    "diff" {
        & "$HiveRoot\scripts\Compare-HiveResults.ps1" -CompareOnly -Clients "xdc-geth-audit"
    }
    default {
        Write-Host @"
Hive Comparison Tool
Usage: .\compare.ps1 <command>

Commands:
  probe       Run RPC probes for xdc-geth-audit and go-ethereum, then show gap matrix
  test        Run xdc/rpc-compat tests (requires Docker, takes ~5 min)
  matrix      Show latest feature gap matrix from stored probe data
  diff        Show latest test comparison from stored results
  help        Show this help
"@
    }
}
