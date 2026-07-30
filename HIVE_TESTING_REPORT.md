# Hive Testing Report: Geth + XDC Gateway + XDC Core Node

**Author:** opencode (OpenCode AI)  
**Date:** 2026-06-17  
**Repository:** `C:\BlocksScan\hive`  
**Objective:** Set up Ethereum Hive test harness and add `xdpos` (XDC Core Node) and `xdc-gateway` (RPC proxy) clients for side-by-side testing with `go-ethereum`.

---

## 1. Executive Summary

| Client | Smoke/Genesis | RPC-compat | XDC RPC-compat | Smoke/XDC | Notes |
|--------|---------------|------------|----------------|-----------|-------|
| `go-ethereum` | **6/6 passed** | **226/227 passed** | N/A | N/A | Standard Ethereum client; baseline. |
| `xdc-gateway` | **6/6 passed** | **226/227 passed** | N/A | N/A | Geth upstream + JSON-RPC proxy. Engine API JWT auth fixed. |
| `xdpos` | **0/6 passed** | N/A | **160/160 passed** | **11/11 passed** | XDC Core Node uses XDPoS consensus; standard Ethereum genesis/simulators are incompatible. Passes custom XDC test suites. |
| `xdc-geth-audit` | N/A | N/A | **134/160 passed** (26 failed) | N/A | XDC go-ethereum audit fork; 26 regressions vs xdpos baseline. See §13. |

Final combined commands run:

```powershell
# Smoke/genesis (pre-merge, no Engine API required)
.\hive --% -sim smoke/genesis -client go-ethereum,xdpos,xdc-gateway

# RPC compatibility (post-merge, requires Engine API)
.\hive --% -sim ethereum/rpc-compat -client go-ethereum,xdc-gateway

# XDC-specific smoke test for XDC Core Node
.\hive --% -sim smoke/xdc -client xdpos

# XDC-specific RPC compatibility
.\hive --% -sim xdc/rpc-compat -client xdpos
```

---

## 2. Environment

| Component | Version / Details |
|-----------|-------------------|
| OS        | Windows 11 |
| Docker    | Docker Desktop 29.5.3, Linux containers |
| Go        | 1.24 (installed at `C:\Program Files\Go`) |
| Hive      | Built from source at `C:\BlocksScan\hive` |
| Clients   | `go-ethereum`, `xdpos` (XDC Core Node), `xdc-gateway` |

### 2.1 Build Hive

```powershell
& "C:\Program Files\Go\bin\go.exe" build -o hive.exe .
```

---

## 3. Issues Fixed During Setup

### 3.1 `go` not in PATH

```text
go : The term 'go' is not recognized...
```

**Fix:** Use full path.

```powershell
& "C:\Program Files\Go\bin\go.exe" build -o hive.exe .
```

### 3.2 CRLF line endings in shell scripts

```text
exec /geth.sh: no such file or directory
```

**Fix:** Convert all `.sh` files to LF.

```powershell
Get-ChildItem -Path . -Recurse -Filter "*.sh" -File | ForEach-Object {
    $content = [System.IO.File]::ReadAllText($_.FullName)
    $content = $content -replace "`r`n", "`n"
    [System.IO.File]::WriteAllText($_.FullName, $content)
}
```

### 3.3 Hive panic on container close

```text
panic: close of closed channel
```

**Fix:** Added `sync.Once` guard and panic recovery in `internal/libdocker/container.go`.

### 3.4 Container has no IP address on Docker Desktop

```text
container has no IP address (check Docker network settings)
```

**Fix:** Added retry loop with fallback to per-network endpoint IPs in `StartContainer`.

### 3.5 XDC Core Node genesis incompatibility

```text
Fatal: Failed to write genesis block: genesis has no chain configuration
Fatal: Only support XDPoS consensus
```

**Fix:** Bundled real XDC Apothem testnet genesis (`genesis-testnet.json`, chainId 51) and added fallback logic in `clients/xdpos/xdpos.sh` to use it when the simulator-provided Ethereum genesis is incompatible.

### 3.6 XDC Core Node masternode etherbase requirement

```text
Fatal: Can't verify masternode permission: etherbase must be explicitly specified
```

**Fix:** Added `--etherbase` flag derived from `HIVE_MINER` in `clients/xdpos/xdpos.sh`, preserving the `xdc...` address prefix.

### 3.7 `xdc-gateway` Engine API JWT authentication

```text
client rejected forkchoiceUpdated: 401 Unauthorized: signature is invalid
```

**Fix:** Wrote Hive's static secret as a 64-character hex JWT secret to `/gateway-jwt-secret` and pointed the upstream Geth `--authrpc.jwtsecret` at it. Removed proxy-side JWT generation so Geth validates the Hive token directly.

---

## 4. Client: go-ethereum

### 4.1 Smoke/Genesis

```powershell
.\hive --% -sim smoke/genesis -client go-ethereum
```

```text
INF simulation smoke/genesis finished suites=1 tests=6 failed=0
```

**Status:** PASS

### 4.2 RPC Compatibility

```powershell
.\hive --% -sim ethereum/rpc-compat -client go-ethereum
```

| Metric | Value |
|--------|-------|
| Total  | 227   |
| Passed | 226   |
| Failed | 1     |

**Failed test:** `eth_config/get-config`

Simulator detail log:

```text
-- eth_config/get-config (go-ethereum)
>>  {"jsonrpc":"2.0","id":1,"method":"eth_config","params":[]}
<<  {"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"the method eth_config does not exist/is not available"}}
```

The current Geth build does not expose `eth_config`, so this single test fails. All other 226 RPC tests pass.

**Status:** PASS (226/227)

---

## 5. Client: xdc-gateway (New)

### 5.1 What it is

A new Hive client that demonstrates the XDC Gateway RPC-proxy pattern. It runs a local Geth node as upstream and exposes it through a lightweight JSON-RPC proxy on port 8545.

### 5.2 Files added/modified

```text
clients/xdc-gateway/
├── Dockerfile       # node:20-alpine + Geth from hive/clients/go-ethereum
├── hive.yaml        # declares eth1 role
├── xdc-gateway.sh   # entry point (starts upstream Geth + proxy)
├── mapper.jq        # genesis mapper (from go-ethereum client)
├── enode.sh         # enode retriever
├── proxy.js         # JSON-RPC gateway proxy (HTTP + Engine API forwarding)
├── package.json     # http-proxy dependency
└── genesis.json     # default genesis template
```

Changes made during this session:

- Added upstream Geth Engine API (`--authrpc.addr 127.0.0.1 --authrpc.port 8651`).
- Added a second Node.js proxy listener on port `8551` that forwards Engine API requests to the upstream Geth.
- Added in-proxy JWT generation using the Hive static secret so the simulator's `client.EngineAPI()` calls are authenticated.

### 5.3 Smoke/Genesis

```powershell
.\hive --% -sim smoke/genesis -client xdc-gateway
```

```text
INF simulation smoke/genesis finished suites=1 tests=6 failed=0
```

**Status:** PASS

### 5.4 RPC Compatibility

```powershell
.\hive --% -sim ethereum/rpc-compat -client xdc-gateway
```

| Metric | Value |
|--------|-------|
| Total  | 227 |
| Passed | 226 |
| Failed | 1 |

**Failed test:** `eth_config/get-config` (same as upstream go-ethereum)

The gateway now forwards Engine API requests with the correct Hive static JWT secret, so the client launch succeeds and all standard RPC tests run. The only failure is the upstream Geth missing `eth_config`.

**Fix applied:**

- Added `--authrpc.jwtsecret /gateway-jwt-secret` to the upstream Geth flags in `xdc-gateway.sh`.
- Wrote Hive's static secret as a 64-character hex string (`7365637265747365637265747365637265747365637265747365637265747365`) to `/gateway-jwt-secret` before starting Geth.
- Removed proxy-side JWT generation; the Engine API listener forwards requests transparently and Geth validates the secret directly.

**Status:** PASS (226/227)

### 5.5 Why xdc-gateway matches go-ethereum

The gateway is a thin proxy in front of the same upstream Geth build. Once Engine API authentication aligned with Hive's static secret, the gateway inherited Geth's RPC conformance (226/227).

---

## 6. Client: xdpos (XDC Core Node) (New)

### 6.1 What it is

A new Hive client wrapping the official XDC Core Node Docker image (`xinfinorg/xinfin-testnet-node:apothem_network`).

### 6.2 Files added/modified

```text
clients/xdpos/
├── Dockerfile           # wraps xinfinorg/xinfin-testnet-node:apothem_network
├── Dockerfile.git       # builds XDC from source (not used in final runs)
├── hive.yaml            # declares eth1 role
├── xdpos.sh             # entry point (uses real XDC testnet genesis)
├── mapper.jq            # converts Geth genesis to XDC-compatible genesis
├── enode.sh             # enode retriever
├── genesis.json         # default genesis template
├── genesis-mainnet.json # real XDC mainnet genesis (chainId 50)
└── genesis-testnet.json # real XDC Apothem testnet genesis (chainId 51)
```

### 6.3 Smoke/Genesis

```powershell
.\hive --% -sim smoke/genesis -client xdpos
```

| Metric | Value |
|--------|-------|
| Total  | 6 |
| Passed | 0 |
| Failed | 6 |

**Status:** FAIL

**Reason:** The simulator supplies a standard Ethereum genesis and expects a specific genesis hash. The XDC Core Node uses **XDPoS consensus**, which requires a different genesis format (real XDC testnet genesis bundled). Even with the real XDC genesis, the resulting hash differs from the Ethereum expected hash, so every `wrong genesis hash` test fails.

Simulator detail log example:

```text
-- empty genesis (xdpos)
genesis hash 0xbdea51...
wrong genesis hash, want 0x433d0b85...
```

XDC client log shows the node starts successfully:

```text
Using genesis: /genesis-testnet.json
Successfully wrote genesis state ... hash=bdea51...
HTTP endpoint opened url=http://0.0.0.0:8545
WebSocket endpoint opened url=ws://[::]:8546
```

### 6.4 RPC Compatibility

```powershell
.\hive --% -sim ethereum/rpc-compat -client xdpos
```

| Metric | Value |
|--------|-------|
| Total  | 1 (client launch) |
| Passed | 0 |
| Failed | 1 |

**Failure:** `client launch (xdpos)` — Engine API port `8551` connection refused.

Simulator detail log:

```text
-- client launch (xdpos)
sending engine_forkchoiceUpdatedV3: [...]
client rejected forkchoiceUpdated: Post "http://172.17.0.8:8551": dial tcp 172.17.0.8:8551: connect: connection refused
```

**Reason:** XDPoS is a pre-merge consensus engine. The XDC Core Node does not implement the post-merge Engine API (`engine_*` methods) on port 8551, so `ethereum/rpc-compat` cannot proceed.

### 6.5 XDC-specific smoke simulator

A new simulator was created at `simulators/smoke/xdc/` to test `xdpos` with the real XDC testnet genesis and a suite of XDC-appropriate RPC checks, including a two-node peering test.

Result:

```powershell
.\hive --% -sim smoke/xdc -client xdpos
```

```text
INF simulation smoke/xdc finished suites=1 tests=11 failed=0
```

**Status:** PASS (11/11)

Tests cover:

- `rpc_modules` lists enabled modules.
- Genesis block hash is returned.
- `net_version` returns a valid network version.
- `eth_blockNumber` is `0x0` at startup.
- `eth_getBlockByNumber` and `eth_getBlockByHash` return the genesis block.
- `eth_accounts` returns successfully (empty list accepted).
- Best-effort XDPoS masternode method check.
- `eth_syncing` is `false` at genesis.
- `net_peerCount` returns valid hex.
- Post-merge Engine API methods are not served.
- Two `xdpos` nodes start, connect via `admin_addPeer`, and report one peer each.

### 6.6 XDC RPC-compatibility simulator (`xdc/rpc-compat`)

A new simulator was created at `simulators/xdc/rpc-compat/` to run curated JSON-RPC tests against `xdpos` using the real XDC Apothem testnet genesis.

Result:

```powershell
.\hive --% -sim xdc/rpc-compat -client xdpos
```

```text
INF simulation xdc/rpc-compat finished suites=1 tests=160 failed=0
```

**Status:** PASS (160/160)

Tests now cover 80+ JSON-RPC methods across admin, debug, eth, miner, net, txpool, web3, and XDPoS namespaces. Fixtures exercise genesis state (all genesis balances, validator storage slots), block queries by number/hash/tag/uncle, raw transaction queries, filter lifecycle, XDPoS snapshot/signers, negative cases (unknown accounts, invalid parameters, missing transactions, unsupported methods), and XDPoS-specific behavior.
- `eth_getBlockByNumber`, `eth_getBlockByHash` (genesis, earliest, full transactions, not found)
- `eth_getUncleByBlockNumberAndIndex`, `eth_getUncleCountByBlockNumber`
- `eth_getBalance` for every genesis account
- `eth_getStorageAt` for validator contract slots 0x7-0xf
- `eth_getRawTransactionByHash`, `eth_getRawTransactionByBlockHashAndIndex`, `eth_getRawTransactionByBlockNumberAndIndex`
- `eth_getBlockSignersByHash`, `eth_getBlockSignersByNumber`, `eth_getBlockFinalityByHash`, `eth_getBlockFinalityByNumber`, `eth_getRewardByHash`, `eth_getCandidateStatus`
- `eth_newBlockFilter`, `eth_newFilter`, `eth_newPendingTransactionFilter`, `eth_getFilterChanges`, `eth_getFilterLogs`, `eth_uninstallFilter`
- `eth_pendingTransactions`
- `admin_addPeer`, `admin_removePeer`, `admin_exportChain`, `admin_importChain`
- `debug_getModifiedAccountsByHash`, `debug_getModifiedAccountsByNumber`, `debug_storageRangeAt`, `debug_traceBlock`, `debug_traceBlockByHash`, `debug_traceBlockByNumber`, `debug_traceBlockFromFile`
- `eth_signTransaction`
- `XDPoS_getSnapshot`, `XDPoS_getSnapshotAtHash`, `XDPoS_getSigners`, `XDPoS_getSignersAtHash`
- `eth_getBlockReceipts`, `eth_getProof`, `eth_getCompensation`, `txpool_contentFrom`, `personal_listAccounts`, `eth_chainId`, `eth_feeHistory`, `eth_blobBaseFee`, `eth_maxPriorityFeePerGas` (unsupported on this XDC build)
- `eth_getTransactionByHash`, `eth_getTransactionReceipt` (not found)
- `eth_getBlockTransactionCountByNumber`, `eth_getUncleCountByBlockNumber`
- `eth_syncing`
- `net_listening`, `net_peerCount`, `net_version`
- `rpc_modules`
- `txpool_content`, `txpool_inspect`, `txpool_status`
- `web3_clientVersion`, `web3_sha3`
- `admin_peers`, `admin_datadir`, `admin_nodeInfo`, `admin_exportChain`, `admin_importChain`
- `admin_addTrustedPeer`, `admin_removeTrustedPeer`, `admin_peerEvents` (unsupported on this XDC build)
- `debug_dumpBlock`, `debug_getBadBlocks`, `debug_getBlockRlp`, `debug_printBlock`, `debug_chaindbProperty`, `debug_chaindbCompact`, `debug_preimage`
- `debug_accountRange` (unsupported on this XDC build)
- `miner_setEtherbase`, `miner_setExtra`, `miner_setGasPrice`, `miner_start` (no signer error), `miner_stop`

This suite uses the real XDC Apothem testnet genesis and validates stable RPC responses.

---

## 7. Combined Run Results

### 7.1 smoke/genesis

```powershell
.\hive --% -sim smoke/genesis -client go-ethereum,xdpos,xdc-gateway
```

```text
INF simulation smoke/genesis finished suites=1 tests=18 failed=6
```

| Client | Tests | Passed | Failed |
|--------|-------|--------|--------|
| go-ethereum | 6 | 6 | 0 |
| xdc-gateway | 6 | 6 | 0 |
| xdpos | 6 | 0 | 6 |

### 7.2 ethereum/rpc-compat

```powershell
.\hive --% -sim ethereum/rpc-compat -client go-ethereum,xdc-gateway
```

```text
INF simulation ethereum/rpc-compat finished suites=1 tests=454 failed=2
```

| Client | Tests | Passed | Failed |
|--------|-------|--------|--------|
| go-ethereum | 227 | 226 | 1 (`eth_config/get-config`) |
| xdc-gateway | 227 | 226 | 1 (`eth_config/get-config`) |
| xdpos | — | — | Not applicable (pre-merge, no Engine API) |

### 7.3 smoke/xdc

```powershell
.\hive --% -sim smoke/xdc -client xdpos
```

```text
INF simulation smoke/xdc finished suites=1 tests=11 failed=0
```

| Client | Tests | Passed | Failed |
|--------|-------|--------|--------|
| xdpos | 11 | 11 | 0 |

### 7.4 xdc/rpc-compat

```powershell
.\hive --% -sim xdc/rpc-compat -client xdpos
```

```text
INF simulation xdc/rpc-compat finished suites=1 tests=160 failed=0
```

```powershell
.\hive --% -sim xdc/rpc-compat -client xdc-geth-audit
```

```text
INF simulation xdc/rpc-compat finished suites=1 tests=160 failed=26
```

| Client | Tests | Passed | Failed |
|--------|-------|--------|--------|
| `xdpos` | 160 | 160 | 0 |
| `xdc-geth-audit` | 160 | 134 | 26 |

---

## 8. How to Reproduce

1. Ensure Docker Desktop is running.
2. Open PowerShell in `C:\BlocksScan\hive`.
3. Build Hive:
   ```powershell
   & "C:\Program Files\Go\bin\go.exe" build -o hive.exe .
   ```
4. Run smoke/genesis:
   ```powershell
   .\hive --% -sim smoke/genesis -client go-ethereum,xdpos,xdc-gateway
   ```
5. Run RPC compatibility (go-ethereum + xdc-gateway):
   ```powershell
   .\hive --% -sim ethereum/rpc-compat -client go-ethereum,xdc-gateway
   ```
6. Run XDC-specific smoke test:
   ```powershell
   .\hive --% -sim smoke/xdc -client xdpos
   ```
7. Run XDC RPC-compat against xdc-geth-audit:
   ```powershell
   .\hive --% -sim xdc/rpc-compat -client xdc-geth-audit
   ```

---

## 9. Result Files

After a run, inspect:

```text
workspace/logs/
├── <timestamp>-<hash>.json              # structured suite results
├── <timestamp>-simulator-<id>.log      # simulator output
├── details/<timestamp>-<id>.log        # per-test detail logs
├── hive.json                            # Hive binary version
├── go-ethereum/
│   └── client-<id>.log                 # Geth node logs
├── xdpos/
│   └── client-<id>.log                 # XDC node logs
└── xdc-gateway/
    └── client-<id>.log                 # Gateway + upstream Geth logs
```

---

## 10. Key Log Messages

### Passing test (go-ethereum / xdc-gateway smoke)

```text
INF API: client go-ethereum started suite=0 test=1 container=...
INF API: test ended suite=0 test=1 pass=true
```

### go-ethereum rpc-compat single failure

```text
-- eth_config/get-config (go-ethereum)
<<  {"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"the method eth_config does not exist/is not available"}}
```

### xdc-gateway rpc-compat success

```text
INF API: client xdc-gateway started suite=0 test=228 container=...
-- eth_config/get-config (xdc-gateway)
<<  {"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"the method eth_config does not exist/is not available"}}
```

The only failure matches upstream go-ethereum.

### xdpos smoke/xdc success

```text
INF API: client xdpos started suite=0 test=1 container=...
INF API: test ended suite=0 test=1 pass=true
```

### xdpos smoke/genesis failure (expected)

```text
-- empty genesis (xdpos)
genesis hash 0xbdea51...
wrong genesis hash, want 0x433d0b85...
```

XDC client log shows it starts regardless:

```text
INFO HTTP endpoint opened url=http://0.0.0.0:8545
INFO WebSocket endpoint opened url=ws://[::]:8546
```

### xdpos rpc-compat incompatibility (expected)

```text
-- client launch (xdpos)
client rejected forkchoiceUpdated: Post "http://172.17.0.8:8551": dial tcp 172.17.0.8:8551: connect: connection refused
```

This is expected because XDPoS is pre-merge and does not expose the Engine API.

---

## 11. Summary

| Item | Status |
|------|--------|
| Hive built from source | Done |
| `go-ethereum` smoke tests passing | **6/6** |
| `go-ethereum` rpc-compat passing | **226/227** |
| `xdc-gateway` client added, smoke passing | **6/6** |
| `xdc-gateway` rpc-compat passing | **226/227** |
| `xdpos` client added, XDC-specific smoke passing | **11/11** |
| `xdpos` client added, XDC RPC-compatibility passing | **160/160** |
| `xdpos` standard Ethereum simulators | Incompatible by design (XDPoS genesis / no Engine API) |
| `xdc-geth-audit` client added, XDC RPC-compatibility run | **134/160** (26 failures, see §13) |

---

## 13. Client: xdc-geth-audit (New)

### 13.1 What it is

A new Hive client wrapping the XDC go-ethereum audit fork (`C:\BlocksScan\xdc-geth-audit`). Built from local source at `clients/xdc-geth-audit/`. Uses XDPoS V2 consensus engine, modern go-ethereum flag style (`--http`, `--ws`), and Apothem testnet genesis.

### 13.2 Files added

```text
clients/xdc-geth-audit/
├── Dockerfile           # 2-stage build: golang:1.26-alpine -> alpine
├── hive.yaml            # declares eth1 role
├── geth.sh              # entry point (--apothem + --http/--ws)
├── enode.sh             # enode retriever
├── xdc-geth-audit/      # local source copy of C:\BlocksScan\xdc-geth-audit
└── .dockerignore        # excludes .git, *.md, docs/
```

### 13.3 XDC RPC-compatibility

```powershell
.\hive --% -sim xdc/rpc-compat -client xdc-geth-audit
```

| Metric | Value |
|--------|-------|
| Total  | 160 |
| Passed | 134 |
| Failed | 26 |

**Failures by category:**

| Category | Count | Tests |
|----------|-------|-------|
| Block format (`miner` prefix, missing fields) | 6 | `eth_getBlockByHash` (2), `eth_getBlockByNumber` (4) |
| Block finality/signers | 4 | `eth_getBlockFinalityByHash`, `eth_getBlockFinalityByNumber`, `eth_getBlockSignersByHash`, `eth_getBlockSignersByNumber` |
| XDPoS methods | 5 | `getBlockInfoByEpochNum`, `getRewardByAccount`, `getSigners`, `getSignersAtHash`, `getV2BlockByHeader` |
| Gas/fee | 2 | `eth_gasPrice`, `eth_feeHistory` |
| Error code mismatches | 5 | `eth_getWork`, `eth_sendTransaction`, `miner_start`, `XDPoS_getBlockInfoByEpochNum`, `debug_accountRange` |
| Response shape mismatches | 3 | `eth_syncing`, `eth_getCandidateStatus`, `admin_nodeInfo` |
| Version/net mismatch | 2 | `web3_clientVersion`, `net_version` |

**Root causes:**

1. **Block `miner` address uses `0x` prefix** instead of `xdc` prefix. The audit fork inherits Geth's standard address formatting and does not apply the XDC address remapping. XDC-specific fields (`penalties`, `validator`, `validators`, `totalDifficulty`) are absent from the block response.

2. **Network ID returns `"89"** instead of `"51"`. The audit fork defaults to network 89 (UAE testnet) rather than Apothem (51). The `--apothem` flag is being used but the embedded genesis may differ.

3. **Gas price returns 1 Gwei** (`0xf4240`) instead of the XDC mainnet 18 Gwei (`0x430e23400`). The audit fork uses the default Geth minimum gas price.

4. **`eth_syncing` returns a sync progress object** instead of `false`. The audit fork's Geth version exposes snapshot sync state even when not syncing.

5. **Error codes differ** for several "method not found" / "unsupported" cases. The xdpos baseline returns `-32000` in some cases where the audit fork returns `-32601` (or vice versa).

6. **`eth_feeHistory` fails** at genesis because no block has been produced yet (head 0, requested block 1).

7. **`net_version` returns `"89"`** (UAE testnet ID) — the genesis supplied by the simulator uses chainId 51 but the client's net_version returns a different default.

### 13.4 Detailed Failure Reference

Per-test diffs are recorded in `error_ledger/xdc-geth-audit-20260729.json`. Key diffs:

| Test | Client | Expected | Root Cause |
|------|--------|----------|------------|
| `eth_getBlockByHash/get-genesis` | `miner: "0x0000..."` | `miner: "xdc0000..."` | Address format; missing XDC fields |
| `eth_gasPrice` | `"0xf4240"` | `"0x430e23400"` | Default Geth min gas price |
| `net_version` | `"89"` | `"51"` | Wrong network ID default |
| `eth_syncing` | Sync object | `false` | Geth snapshot sync state |
| `eth_sendTransaction` | `-32000` | `-32600` | Error code mismatch |
| `eth_getWork` | `-32601` | `-32000` | Method not found vs internal error |
| `miner_start` | `-32601` | `-32000` | Method removed vs no signer |
| `eth_getBlockFinalityByHash` | `33` | `0` | Finality value changed |
| `eth_getBlockSignersByHash` | `["0x0000..."]` | `[]` | Empty vs zero-address |
| `web3_clientVersion` | `"XDC/v1.17.4..."` | `"XDC/v1.0.0..."` | Version string diff |

---

## 12. Next Steps

1. **~~Fix `xdc-gateway` Engine API JWT forwarding~~** — Done.
2. **~~Complete `simulators/smoke/xdc`~~** — Done.
3. **~~Run `ethereum/rpc-compat` again~~** — Done; both go-ethereum and xdc-gateway score 226/227.
4. **Update `HIVE_LOG_GUIDE.md`** with the resolved Engine API pattern and the new `smoke/xdc` simulator.
5. **Consider an XDPoS-specific sync/consensus simulator** if deeper XDC coverage is required.
6. **Investigate xdc-geth-audit regressions** — 26 failures against xdpos baseline. Key categories:
   - Block format: `miner` address uses `0x` prefix (not `xdc`), missing XDC fields (`penalties`, `validator`, `validators`, `totalDifficulty`).
   - Network ID: returns `"89"` (UAE testnet) instead of `"51"` (Apothem).
   - Gas price: returns 1 Gwei default instead of 18 Gwei.
   - Syncing: returns snapshot sync progress object instead of `false`.
   - Error code mismatches across multiple methods.
7. **Update fixtures** to match audit fork response format, or backport the missing XDC fields into the audit source.
