# Hive Log Output Guide: Geth + XDC Gateway + XDC Core Node Testing

This document explains the logs produced when running Hive tests against:

- `go-ethereum`
- `xdc-gateway` (new RPC-proxy client)
- `xdpos` (new XDC Core Node client)

Actual runs:

```powershell
# Standard Ethereum genesis smoke test
.\hive --% -sim smoke/genesis -client go-ethereum,xdpos,xdc-gateway

# Post-merge RPC compatibility
.\hive --% -sim ethereum/rpc-compat -client go-ethereum,xdc-gateway

# XDC-specific smoke test
.\hive --% -sim smoke/xdc -client xdpos
```

Actual results:

```text
INF simulation smoke/genesis finished suites=1 tests=18 failed=6
INF simulation ethereum/rpc-compat finished suites=1 tests=454 failed=2
INF simulation smoke/xdc finished suites=1 tests=1 failed=0
```

---

## 1. Where Logs Come From

When you run Hive, three things produce logs:

| Source | Where it appears | What it tells you |
|--------|------------------|-------------------|
| **Hive controller** | Terminal / console output | What Hive is doing: building images, starting containers, test status |
| **Simulator container** | `workspace/logs/<timestamp>-<hash>-simulator.log` | Test logic output: what the simulator requested, what it checked |
| **Client container(s)** | `workspace/logs/<client>/client-<full-container-id>.log` | Node output: Geth, XDC, or gateway proxy startup, RPC responses, errors |
| **Result JSON** | `workspace/logs/<timestamp>-<hash>.json` | Structured pass/fail data with timings |

---

## 2. Console Output Explained

### 2.1 Startup phase

```text
INF building image image=hive/hiveproxy
INF building 3 clients...
INF building image image=hive/clients/go-ethereum:latest dir=clients\go-ethereum
INF building image image=hive/clients/xdpos:latest dir=clients\xdpos
INF building image image=hive/clients/xdc-gateway:latest dir=clients\xdc-gateway
INF building 1 simulators...
INF building image image=hive/simulators/smoke/genesis:latest dir=simulators\smoke\genesis
```

| Line | Meaning |
|------|---------|
| `building image image=hive/hiveproxy` | Hive builds a small proxy container that helps containers talk to each other. |
| `building 3 clients...` | It is building three client images: `go-ethereum`, `xdpos`, and `xdc-gateway`. |
| `building image image=hive/clients/...` | Docker is building each client wrapper image. |
| `building 1 simulators...` | Building the simulator image (`smoke/genesis`). |
| `building image image=hive/simulators/smoke/genesis:latest` | Docker build for the test program. |

### 2.2 Simulation phase

```text
INF running simulation: smoke/genesis
INF hiveproxy started container=895c6ef3a09d addr=172.17.0.2:8081
INF API: suite started suite=0 name=genesis
INF API: test started suite=0 test=1 name="empty genesis (go-ethereum)"
INF API: client go-ethereum started suite=0 test=1 container=0e501046
INF API: test ended suite=0 test=1 pass=true
INF API: test started suite=0 test=2 name="empty genesis (xdpos)"
INF API: client xdpos started suite=0 test=2 container=8462d2f6
INF API: test ended suite=0 test=2 pass=false
INF API: test started suite=0 test=3 name="empty genesis (xdc-gateway)"
INF API: client xdc-gateway started suite=0 test=3 container=17b5a62b
INF API: test ended suite=0 test=3 pass=true
...
INF API: suite ended suite=0
INF simulation smoke/genesis finished suites=1 tests=18 failed=6
```

| Line | Meaning |
|------|---------|
| `running simulation: smoke/genesis` | The simulator is now executing. |
| `hiveproxy started container=895c6ef3a09d` | The proxy container is running at IP `172.17.0.2:8081`. |
| `API: suite started suite=0 name=genesis` | A test suite called "genesis" begins. Simulators can have multiple suites. |
| `API: test started suite=0 test=1 name="empty genesis (go-ethereum)"` | Test case #1 starts: launching Geth with an empty genesis. |
| `API: client go-ethereum started ... container=0e501046` | Hive successfully launched the client container and it opened RPC port 8545. |
| `API: test ended suite=0 test=1 pass=true` | The simulator reported that the test passed. |
| `API: test ended suite=0 test=2 pass=false` | The XDC node started but the genesis hash did not match. |
| `API: suite ended suite=0` | All tests in suite 0 are done. |
| `simulation ... finished suites=1 tests=18 failed=6` | **Final result.** 1 suite, 18 tests (6 per client), 6 failures (all in `xdpos`). |

### 2.3 Failure example (xdpos)

```text
INF API: client xdpos started suite=0 test=2 container=8462d2f6
INF API: test ended suite=0 test=2 pass=false
```

This means the XDC container **started and opened RPC**, but the simulator's check failed (in this case, the genesis hash was wrong).

To find the reason, look at:

1. `workspace/logs/details/<timestamp>-<simulator-id>.log`
2. `workspace/logs/xdpos/client-<id>.log`

---

## 3. Result JSON File Explained

After the run, Hive writes a JSON file like:

```text
workspace/logs/1781546433-f5313cf4804cc717976cc76b09482bad.json
```

### 3.1 Top-level fields

```json
{
  "id": 0,
  "name": "genesis",
  "description": "This test suite checks client initialization with genesis blocks.",
  "clientVersions": {
    "go-ethereum": "Geth/v1.17.4-unstable-e2164cc7-20260615/linux-amd64/go1.26.4",
    "xdpos": "unknown",
    "xdc-gateway": "xdc-gateway-hive-client"
  },
  "runMetadata": {
    "hiveCommand": [
      "C:\\BlocksScan\\hive\\hive.exe",
      "-sim", "smoke/genesis",
      "-client", "go-ethereum",
      "-loglevel", "5"
    ],
    "hiveVersion": {
      "commit": "4db8f994...",
      "commitDate": "2026-06-15T02:27:19Z",
      "branch": "master",
      "dirty": true
    }
  },
  "testCases": { ... }
}
```

| Field | Meaning |
|-------|---------|
| `id` | Suite number (0 for the first suite). |
| `name` | Name of the test suite. |
| `description` | What the suite checks. |
| `clientVersions` | Version strings reported by each client (read from `/version.txt` inside the container). |
| `runMetadata.hiveCommand` | The exact command you ran. Useful for reproducing. |
| `runMetadata.hiveVersion` | Git commit/branch of the Hive binary. `dirty=true` means local modifications exist. |
| `testCases` | Map of test numbers to test results. |

### 3.2 testCases fields

```json
{
  "1": {
    "name": "empty genesis (go-ethereum)",
    "description": "This imports an empty genesis block with no environment variables.",
    "start": "2026-06-15T23:30:16.4123375+05:30",
    "end": "2026-06-15T23:30:19.5769737+05:30",
    "summaryResult": {
      "pass": true,
      "details": ""
    },
    "clientInfo": {
      "b5482d41": {
        "id": "b5482d41",
        "ip": "172.17.0.4",
        "name": "go-ethereum",
        "instantiatedAt": "2026-06-15T23:30:19.5509195+05:30",
        "logFile": "go-ethereum/client-b5482d41....log",
        "logOffsets": { "begin": 41, "end": 41 }
      }
    }
  }
}
```

| Field | Meaning |
|-------|---------|
| `name` | Human-readable test name, usually includes the client name. |
| `description` | What the test does. |
| `start` / `end` | ISO timestamps for when the test started and finished. |
| `summaryResult.pass` | `true` = passed, `false` = failed. |
| `summaryResult.details` | Optional failure message or extra info. |
| `clientInfo` | Map of container short IDs to client details. |
| `clientInfo.<id>.ip` | Container IP address. Empty if the container never got an IP. |
| `clientInfo.<id>.name` | Client name: `go-ethereum` or `xdpos`. |
| `clientInfo.<id>.instantiatedAt` | When the client container was created. |
| `clientInfo.<id>.logFile` | Path to the client detail log. |
| `clientInfo.<id>.logOffsets` | Byte offsets in the simulator log that correspond to this test. |

---

## 4. Simulator Log Explained

File pattern:

```text
workspace/logs/<timestamp>-<hash>-simulator-<container-id>.log
```

Example content (passing `smoke/genesis` for go-ethereum / xdc-gateway):

```text
-- empty genesis (go-ethereum)
genesis hash 0x433d0b859a77a29753d2a6df477c971dcc6300af33f9d64d821a1d490b4148b1

-- empty genesis (xdc-gateway)
genesis hash 0x433d0b859a77a29753d2a6df477c971dcc6300af33f9d64d821a1d490b4148b1
```

| Line | Meaning |
|------|---------|
| `-- empty genesis (go-ethereum)` | Separator showing which test this section belongs to. |
| `genesis hash 0x433d...` | The simulator computed the genesis hash from the running client. It matches the expected hash, so the test passes. |

Example content (failing `smoke/genesis` for xdpos):

```text
-- empty genesis (xdpos)
genesis hash 0x2e8b8aab256cb1904129fcc60d294e87ec26998234c59b06f70c382607667358
wrong genesis hash, want 0x433d0b859a77a29753d2a6df477c971dcc6300af33f9d64d821a1d490b4148b1
```

The XDC node returned a different genesis hash because the XDPoS consensus engine requires a different genesis format than standard Ethereum.

---

## 5. Client Detail Log Explained

File pattern:

```text
workspace/logs/go-ethereum/client-<full-container-id>.log
workspace/logs/xdpos/client-<full-container-id>.log
workspace/logs/xdc-gateway/client-<full-container-id>.log
```

This is the **actual stdout/stderr of the client container**.

### 5.1 Geth / xdc-gateway example (passing)

Both clients produce similar logs because `xdc-gateway` runs Geth as its upstream node.

```text
INFO [06-16|...] Initialised chain configuration
INFO [06-16|...] Initialising Ethereum protocol
INFO [06-16|...] Loaded most recent local header
INFO [06-16|...] HTTP server started
```

| Line | Meaning |
|------|---------|
| `Initialised chain configuration` | Chain config (forks, chain ID) loaded. |
| `Initialising Ethereum protocol` | Geth is initializing the Ethereum protocol. |
| `Loaded most recent local header` | Genesis block is now the current head. |
| `HTTP server started` | JSON-RPC is ready on port 8545. Hive detects this and marks the client online. |

For `xdc-gateway`, you will also see proxy lines:

```text
XDC Gateway proxy listening on 0.0.0.0:8545 -> http://127.0.0.1:8546
[gateway] POST /rpc
```

### 5.2 XDC example (container starts, test fails)

```text
Using XDC binary: /usr/bin/XDC
Converting genesis.json to XDC format...
Initializing XDC datadir with genesis...
INFO [06-16|...] Successfully wrote genesis state
Starting XDC with flags: ...
INFO [06-16|...] HTTP endpoint opened url=http://0.0.0.0:8545
INFO [06-16|...] WebSocket endpoint opened url=ws://[::]:8546
```

The container starts and opens RPC. The test fails later because the genesis hash does not match the Ethereum expected hash.

### 5.3 XDC startup failure examples

During development, the old XDC image rejected the Ethereum genesis for several XDPoS-specific reasons:

```text
Fatal: Failed to write genesis block: genesis has no chain configuration
Fatal: Can't verify masternode permission: etherbase must be explicitly specified
Fatal: Can't verify masternode permission: extra-data 32 byte vanity prefix missing
Fatal: Can't verify masternode permission: non-zero mix digest
Fatal: Failed to start staking: signer missing: unknown account
```

These were resolved by mapping the genesis to an ethash-compatible config, at the cost of changing the genesis hash.

---

## 6. Engine API / JWT Example (`xdc-gateway`)

When `ethereum/rpc-compat` launches a client, it first calls `engine_forkchoiceUpdatedV3` on port 8551 with a JWT signed by Hive's static secret. Geth's `--authrpc` rejects the call unless it shares the same secret.

### Before the fix

```text
-- client launch (xdc-gateway)
client rejected forkchoiceUpdated: 401 Unauthorized: signature is invalid
```

### After the fix

`xdc-gateway.sh` writes the Hive static secret as a 64-character hex file and starts upstream Geth with `--authrpc.jwtsecret /gateway-jwt-secret`:

```text
Using JWT secret: /gateway-jwt-secret
Starting upstream Geth...
```

The client launch succeeds:

```text
INF API: client xdc-gateway started suite=0 test=228 container=...
```

and all standard RPC tests run. The only failure is the same `eth_config/get-config` missing-method error seen on plain go-ethereum.

---

## 7. XDC-Specific Smoke Simulator (`smoke/xdc`)

Because `xdpos` cannot pass Ethereum genesis or Engine-API simulators, a dedicated simulator was added at `simulators/smoke/xdc/`.

Run:

```powershell
.\hive --% -sim smoke/xdc -client xdpos
```

Console output:

```text
INF API: test started suite=0 test=1 name="xdc rpc modules (xdpos)"
INF API: client xdpos started suite=0 test=1 container=...
...
INF API: test ended suite=0 test=11 pass=true
INF simulation smoke/xdc finished suites=1 tests=11 failed=0
```

Tests cover basic RPC, genesis block, and two-node peering.

---

## 8. XDC RPC-Compatibility Simulator (`xdc/rpc-compat`)

A separate simulator at `simulators/xdc/rpc-compat/` runs curated JSON-RPC tests against `xdpos` using the real XDC Apothem testnet genesis.

Run:

```powershell
.\hive --% -sim xdc/rpc-compat -client xdpos
```

Console output:

```text
INF simulation xdc/rpc-compat finished suites=1 tests=160 failed=0
```

Tests cover:

- `eth_accounts`, `eth_blockNumber`
- `eth_getBalance` (all genesis accounts), `eth_getBlockByHash`, `eth_getBlockByNumber` (genesis, earliest, full transactions, not found), `eth_getCode`, `eth_getStorageAt` (validator slots 0x7-0xf), `eth_getTransactionCount`
- `eth_getTransactionByHash`, `eth_getTransactionReceipt` (not found)
- `eth_getRawTransactionByHash`, `eth_getRawTransactionByBlockHashAndIndex`, `eth_getRawTransactionByBlockNumberAndIndex`
- `eth_getBlockTransactionCountByNumber`, `eth_getUncleCountByBlockNumber`, `eth_getUncleByBlockNumberAndIndex`, `eth_getUncleByBlockHashAndIndex`
- `eth_getBlockSignersByHash`, `eth_getBlockSignersByNumber`, `eth_getBlockFinalityByHash`, `eth_getBlockFinalityByNumber`, `eth_getRewardByHash`, `eth_getCandidateStatus`
- `eth_newBlockFilter`, `eth_newFilter`, `eth_newPendingTransactionFilter`, `eth_getFilterChanges`, `eth_getFilterLogs`, `eth_uninstallFilter`
- `eth_pendingTransactions`
- `eth_getBlockReceipts`, `eth_getProof`, `eth_getCompensation`, `txpool_contentFrom`, `personal_listAccounts`, `eth_chainId`, `eth_feeHistory`, `eth_blobBaseFee`, `eth_maxPriorityFeePerGas` (unsupported on this XDC build)
- `eth_syncing`
- `net_listening`, `net_peerCount`, `net_version`
- `rpc_modules`
- `txpool_content`, `txpool_inspect`, `txpool_status`
- `web3_clientVersion`, `web3_sha3`
- `admin_peers`, `admin_datadir`, `admin_nodeInfo`, `admin_exportChain`, `admin_importChain`, `admin_addPeer`, `admin_removePeer`
- `admin_addTrustedPeer`, `admin_removeTrustedPeer`, `admin_peerEvents` (unsupported on this XDC build)
- `debug_dumpBlock`, `debug_getBadBlocks`, `debug_getBlockRlp`, `debug_printBlock`, `debug_chaindbProperty`, `debug_chaindbCompact`, `debug_preimage`, `debug_traceTransaction`
- `debug_getModifiedAccountsByHash`, `debug_getModifiedAccountsByNumber`, `debug_storageRangeAt`, `debug_traceBlock`, `debug_traceBlockByHash`, `debug_traceBlockByNumber`, `debug_traceBlockFromFile`
- `debug_accountRange` (unsupported on this XDC build)
- `miner_setEtherbase`, `miner_setExtra`, `miner_setGasPrice`, `miner_start` (no signer error), `miner_stop`
- `eth_signTransaction` (no account error)
- XDPoS-specific: `eth_getCandidateStatus`, `eth_getMasternodes`, `eth_getMasternodeInfo`, `eth_getVoters`, `eth_getRewards`, `eth_getBlockFinality`
- XDPoS namespace: `XDPoS_getSnapshot`, `XDPoS_getSnapshotAtHash`, `XDPoS_getSigners`, `XDPoS_getSignersAtHash`

---

## 9. What Changes When Adding `xdpos` and `xdc-gateway`

When you run:

```powershell
.\hive --% -sim smoke/genesis -client go-ethereum,xdpos,xdc-gateway
```

### Console output

```text
INF building 3 clients...
INF building image image=hive/clients/go-ethereum:latest
INF building image image=hive/clients/xdpos:latest
INF building image image=hive/clients/xdc-gateway:latest
...
INF API: test started suite=0 test=1 name="empty genesis (go-ethereum)"
INF API: test ended suite=0 test=1 pass=true
INF API: test started suite=0 test=2 name="empty genesis (xdpos)"
INF API: client xdpos started suite=0 test=2 container=...
INF API: test ended suite=0 test=2 pass=false
INF API: test started suite=0 test=3 name="empty genesis (xdc-gateway)"
INF API: client xdc-gateway started suite=0 test=3 container=...
INF API: test ended suite=0 test=3 pass=true
...
INF simulation smoke/genesis finished suites=1 tests=18 failed=6
```

- 3 clients built.
- 6 tests per client = 18 total tests.
- Tests 1, 4, 7, 10, 13, 16 are Geth.
- Tests 2, 5, 8, 11, 14, 17 are XDC Core Node.
- Tests 3, 6, 9, 12, 15, 18 are XDC Gateway.
- 6 failures are all in the `xdpos` client due to genesis hash mismatch.

### Result JSON

```json
{
  "clientVersions": {
    "go-ethereum": "Geth/v1.17.4...",
    "xdpos": "unknown",
    "xdc-gateway": "xdc-gateway-hive-client"
  },
  "testCases": {
    "1": { "name": "empty genesis (go-ethereum)", "summaryResult": { "pass": true }, ... },
    "2": { "name": "empty genesis (xdpos)", "summaryResult": { "pass": false }, ... },
    "3": { "name": "empty genesis (xdc-gateway)", "summaryResult": { "pass": true }, ... }
  }
}
```

### File tree after run

```text
workspace/logs/
├── <timestamp>-<hash>.json
├── <timestamp>-simulator-<id>.log
├── hive.json
├── go-ethereum/
│   └── client-<id>.log
├── xdpos/
│   └── client-<id>.log
└── xdc-gateway/
    └── client-<id>.log
```

---

## 10. Reading a Failed Test

### Console

```text
INF API: client xdpos started suite=0 test=2 container=8462d2f6
INF API: test ended suite=0 test=2 pass=false
```

### Simulator detail log

```text
-- empty genesis (xdpos)
genesis hash 0x2e8b8aab256cb1904129fcc60d294e87ec26998234c59b06f70c382607667358
wrong genesis hash, want 0x433d0b859a77a29753d2a6df477c971dcc6300af33f9d64d821a1d490b4148b1
```

### JSON

```json
{
  "2": {
    "name": "empty genesis (xdpos)",
    "summaryResult": {
      "pass": false,
      "details": "wrong genesis hash"
    },
    "clientInfo": {
      "8462d2f6": {
        "name": "xdpos",
        "ip": "172.17.0.4",
        "logFile": "xdpos/client-8462d2f6....log"
      }
    }
  }
}
```

### What to check

1. `workspace/logs/xdpos/client-8462d2f6....log` — did XDC start? Look for `HTTP endpoint opened`.
2. `workspace/logs/<timestamp>-simulator-....log` — what hash did it expect vs get?
3. If XDC did not start, check the mapped genesis in the client log.

---

## 11. Summary Table

| Log/File | Contains | When to look at it |
|----------|----------|--------------------|
| Terminal output | High-level progress and final pass/fail | Quick status check |
| `workspace/logs/<ts>-<hash>.json` | Structured results, timings, client versions | Reporting to sir |
| `workspace/logs/<ts>-simulator-...log` | Simulator assertions and errors | Understanding why a test failed |
| `workspace/logs/<client>/client-...log` | Client node stdout/stderr | Debugging client startup or RPC issues |
| `workspace/logs/hive.json` | Hive binary version/commit | Reproducibility info |

---

## 12. Key Numbers to Report

For every run, report these:

```text
simulation <name> finished suites=<N> tests=<N> failed=<N>
```

Examples:

| Run | Expected output |
|-----|-----------------|
| Geth only (smoke/genesis) | `finished suites=1 tests=6 failed=0` |
| xdc-gateway only (smoke/genesis) | `finished suites=1 tests=6 failed=0` |
| Geth + xdc-gateway (smoke/genesis) | `finished suites=1 tests=12 failed=0` |
| Geth + xdpos + xdc-gateway (smoke/genesis) | `finished suites=1 tests=18 failed=6` |
| Geth + xdc-gateway (rpc-compat) | `finished suites=1 tests=454 failed=2` |
| xdpos only (smoke/xdc) | `finished suites=1 tests=11 failed=0` |
| xdpos only (xdc/rpc-compat) | `finished suites=1 tests=160 failed=0` |

For the actual runs documented here:

```text
simulation smoke/genesis finished suites=1 tests=18 failed=6
simulation ethereum/rpc-compat finished suites=1 tests=454 failed=2
simulation smoke/xdc finished suites=1 tests=11 failed=0
simulation xdc/rpc-compat finished suites=1 tests=160 failed=0
```

- `go-ethereum`: 6/6 smoke/genesis, 226/227 rpc-compat
- `xdc-gateway`: 6/6 smoke/genesis, 226/227 rpc-compat
- `xdpos`: 0/6 smoke/genesis (expected), 11/11 smoke/xdc, 160/160 xdc/rpc-compat
