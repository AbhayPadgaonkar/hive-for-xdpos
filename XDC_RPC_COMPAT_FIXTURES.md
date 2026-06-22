# XDC RPC-compatibility Fixture Tracker

This document tracks which RPC fixtures are implemented for the `xdpos` client in `simulators/xdc/rpc-compat/tests/`.

Last updated: 2026-06-22
Current passing count: **160/160**

## Legend

- `[x]` — fixture implemented and passing
- `[ ]` — fixture not yet implemented
- `[-]` — method unsupported / not applicable in current `xdpos` image

---

## admin

| Fixture | Status | Notes |
|---------|--------|-------|
| `admin_addPeer` | [x] | valid enode returns true |
| `admin_addTrustedPeer` | [x] | unsupported negative fixture |
| `admin_datadir` | [x] | |
| `admin_exportChain` | [x] | |
| `admin_importChain` | [x] | |
| `admin_nodeInfo` | [x] | |
| `admin_peerEvents` | [x] | unsupported negative fixture |
| `admin_peers` | [x] | |
| `admin_removePeer` | [x] | valid enode returns true |
| `admin_removeTrustedPeer` | [x] | unsupported negative fixture |
| `admin_startHTTP` | [-] | unsupported on this XDC build |
| `admin_startRPC` | [-] | unsupported on this XDC build |
| `admin_startWS` | [-] | unsupported on this XDC build |
| `admin_stopHTTP` | [-] | unsupported on this XDC build |
| `admin_stopRPC` | [-] | unsupported on this XDC build |
| `admin_stopWS` | [-] | unsupported on this XDC build |

## debug

| Fixture | Status | Notes |
|---------|--------|-------|
| `debug_accountRange` | [x] | unsupported negative fixture |
| `debug_chaindbCompact` | [x] | |
| `debug_chaindbProperty` | [x] | asserts only presence of result |
| `debug_dbGet` | [-] | unsupported on this XDC build |
| `debug_dumpBlock` | [x] | |
| `debug_getBadBlocks` | [x] | |
| `debug_getBlockRlp` | [x] | |
| `debug_getModifiedAccountsByHash` | [x] | genesis error negative |
| `debug_getModifiedAccountsByNumber` | [x] | genesis error negative |
| `debug_intermediateRoots` | [-] | unsupported on this XDC build |
| `debug_preimage` | [x] | missing preimage negative |
| `debug_printBlock` | [x] | asserts only presence of result |
| `debug_setHead` | [-] | local-only (PrivateDebugAPI) |
| `debug_storageRangeAt` | [x] | genesis error negative |
| `debug_traceBadBlock` | [-] | unsupported on this XDC build |
| `debug_traceBlock` | [x] | genesis error negative |
| `debug_traceBlockByHash` | [x] | genesis error negative |
| `debug_traceBlockByNumber` | [x] | genesis error negative |
| `debug_traceBlockFromFile` | [x] | missing file negative |
| `debug_traceCall` | [-] | unsupported on this XDC build |
| `debug_traceChain` | [-] | unsupported on this XDC build |
| `debug_traceTransaction` | [x] | not-found negative |
| profiling methods | [-] | exist but stateful; skipped |

## eth

| Fixture | Status | Notes |
|---------|--------|-------|
| `eth_accounts` | [x] | |
| `eth_blobBaseFee` | [x] | unsupported negative fixture |
| `eth_blockNumber` | [x] | |
| `eth_call` | [x] | |
| `eth_chainId` | [x] | unsupported negative fixture |
| `eth_coinbase` | [x] | |
| `eth_createAccessList` | [-] | unsupported on this XDC build |
| `eth_estimateGas` | [x] | |
| `eth_feeHistory` | [x] | unsupported negative fixture |
| `eth_gasPrice` | [x] | |
| `eth_getAccountInfo` | [-] | unsupported on this XDC build |
| `eth_getBalance` | [x] | all genesis accounts covered |
| `eth_getBlockByHash` | [x] | |
| `eth_getBlockByNumber` | [x] | |
| `eth_getBlockFinality` | [x] | unsupported negative fixture |
| `eth_getBlockFinalityByHash` | [x] | |
| `eth_getBlockFinalityByNumber` | [x] | |
| `eth_getBlockReceipts` | [x] | unsupported negative fixture |
| `eth_getBlockSignersByHash` | [x] | |
| `eth_getBlockSignersByNumber` | [x] | |
| `eth_getBlockTransactionCountByHash` | [x] | |
| `eth_getBlockTransactionCountByNumber` | [x] | |
| `eth_getCandidateStatus` | [x] | |
| `eth_getCandidates` | [-] | unsupported on this XDC build |
| `eth_getCode` | [x] | |
| `eth_getCompensation` | [x] | unsupported negative fixture |
| `eth_getFilterChanges` | [x] | non-existent filter error |
| `eth_getFilterLogs` | [x] | non-existent filter error |
| `eth_getHeaderByHash` | [-] | unsupported on this XDC build |
| `eth_getHeaderByNumber` | [-] | unsupported on this XDC build |
| `eth_getLogs` | [x] | |
| `eth_getMasternodeInfo` | [x] | unsupported negative fixture |
| `eth_getMasternodes` | [x] | bad-param negative |
| `eth_getMaxPriorityFeePerGas` | [x] | unsupported negative fixture |
| `eth_getOwnerByCoinbase` | [-] | unsupported on this XDC build |
| `eth_getProof` | [x] | unsupported negative fixture |
| `eth_getRawTransactionByBlockHashAndIndex` | [x] | |
| `eth_getRawTransactionByBlockNumberAndIndex` | [x] | |
| `eth_getRawTransactionByHash` | [x] | |
| `eth_getRewardByHash` | [x] | |
| `eth_getRewards` | [x] | unsupported negative fixture |
| `eth_getStorageAt` | [x] | validator slots 0x7-0xf |
| `eth_getStorageValues` | [-] | unsupported on this XDC build |
| `eth_getTransactionAndReceiptProof` | [-] | unsupported on this XDC build |
| `eth_getTransactionByBlockHashAndIndex` | [x] | |
| `eth_getTransactionByBlockNumberAndIndex` | [x] | |
| `eth_getTransactionByHash` | [x] | |
| `eth_getTransactionCount` | [x] | |
| `eth_getTransactionReceipt` | [x] | |
| `eth_getUncleByBlockHashAndIndex` | [x] | |
| `eth_getUncleByBlockNumberAndIndex` | [x] | |
| `eth_getUncleCountByBlockHash` | [x] | |
| `eth_getUncleCountByBlockNumber` | [x] | |
| `eth_getVoters` | [x] | unsupported negative fixture |
| `eth_getWork` | [x] | |
| `eth_hashrate` | [x] | |
| `eth_mining` | [x] | |
| `eth_newBlockFilter` | [x] | asserts only presence of id |
| `eth_newFilter` | [x] | asserts only presence of id |
| `eth_newPendingTransactionFilter` | [x] | asserts only presence of id |
| `eth_pendingTransactions` | [x] | |
| `eth_protocolVersion` | [x] | |
| `eth_sendRawTransaction` | [x] | invalid tx negative |
| `eth_sendTransaction` | [x] | IPC-only negative |
| `eth_sign` | [x] | no account negative |
| `eth_signTransaction` | [x] | no account negative |
| `eth_simulateV1` | [-] | unsupported on this XDC build |
| `eth_submitHashrate` | [x] | |
| `eth_submitWork` | [x] | |
| `eth_syncing` | [x] | |
| `eth_uninstallFilter` | [x] | non-existent filter returns false |

## miner

| Fixture | Status | Notes |
|---------|--------|-------|
| `miner_getHashrate` | [x] | |
| `miner_getWork` | [-] | unsupported on this XDC build (eth_getWork used) |
| `miner_setEtherbase` | [x] | |
| `miner_setExtra` | [x] | |
| `miner_setGasPrice` | [x] | |
| `miner_start` | [x] | no signer error |
| `miner_stop` | [x] | |
| `miner_submitHashrate` | [-] | unsupported on this XDC build |
| `miner_submitWork` | [-] | unsupported on this XDC build |

## net

| Fixture | Status | Notes |
|---------|--------|-------|
| `net_listening` | [x] | |
| `net_peerCount` | [x] | |
| `net_version` | [x] | |

## personal

| Fixture | Status | Notes |
|---------|--------|-------|
| `personal_listAccounts` | [x] | unsupported negative fixture |

## rpc

| Fixture | Status | Notes |
|---------|--------|-------|
| `rpc_modules` | [x] | |

## txpool

| Fixture | Status | Notes |
|---------|--------|-------|
| `txpool_content` | [x] | |
| `txpool_contentFrom` | [x] | unsupported negative fixture |
| `txpool_inspect` | [x] | |
| `txpool_status` | [x] | |

## web3

| Fixture | Status | Notes |
|---------|--------|-------|
| `web3_clientVersion` | [x] | |
| `web3_sha3` | [x] | |

## XDPoS (requires `--rpcapi XDPoS`)

| Fixture | Status | Notes |
|---------|--------|-------|
| `XDPoS_getSnapshot` | [x] | genesis snapshot |
| `XDPoS_getSnapshotAtHash` | [x] | genesis snapshot by hash |
| `XDPoS_getSigners` | [x] | genesis signers |
| `XDPoS_getSignersAtHash` | [x] | genesis signers by hash |
| `XDPoS_getMasternodesByNumber` | [x] | unsupported negative fixture |
| `XDPoS_getLatestPoolStatus` | [x] | unsupported negative fixture |
| `XDPoS_getV2BlockByHeader` | [x] | unsupported negative fixture |
| `XDPoS_getV2BlockByNumber` | [x] | unsupported negative fixture |
| `XDPoS_getV2BlockByHash` | [x] | unsupported negative fixture |
| `XDPoS_networkInformation` | [x] | unsupported negative fixture |
| `XDPoS_config` | [x] | unsupported negative fixture |
| `XDPoS_getMissedRoundsInEpochByBlockNum` | [x] | unsupported negative fixture |
| `XDPoS_getRewardByAccount` | [x] | unsupported negative fixture |
| `XDPoS_getEpochNumbersBetween` | [x] | unsupported negative fixture |
| `XDPoS_getBlockInfoByV2EpochNum` | [x] | unsupported negative fixture |
| `XDPoS_calculateBlockInfoByV1EpochNum` | [x] | unsupported negative fixture |
| `XDPoS_getBlockInfoByEpochNum` | [x] | unsupported negative fixture |

---

## Running the suite

```powershell
.\hive.exe --sim "xdc/rpc-compat" --client=xdpos --results-root=workspace/logs
```

Expected output:

```text
INF simulation xdc/rpc-compat finished suites=1 tests=160 failed=0
```

## What changed recently

- Enabled `XDPoS` namespace in `clients/xdpos/xdpos.sh`.
- Added `admin_addPeer`, `admin_removePeer`.
- Added `debug_getModifiedAccountsByHash`, `debug_getModifiedAccountsByNumber`, `debug_storageRangeAt`, `debug_traceBlock`, `debug_traceBlockByHash`, `debug_traceBlockByNumber`, `debug_traceBlockFromFile`.
- Added `eth_signTransaction`.
- Added `XDPoS_getSnapshot`, `XDPoS_getSnapshotAtHash`, `XDPoS_getSigners`, `XDPoS_getSignersAtHash` and unsupported XDPoS negative fixtures.
