# XDC RPC-compatibility Fixtures

This document lists all static .io fixtures in simulators/xdc/rpc-compat/tests/ and what each one verifies against the xdpos client.

Total fixtures: **102**

| Fixture | Description |
|---------|-------------|
| `admin_datadir/get-datadir.io` | admin_datadir returns the data directory path. |
| `admin_nodeInfo/get-node-info.io` | admin_nodeInfo returns node identity information. |
| `admin_peers/get-peers.io` | admin_peers returns the list of connected peers. |
| `debug_dumpBlock/get-genesis.io` | debug_dumpBlock returns genesis block state dump. |
| `debug_getBadBlocks/get-empty.io` | debug_getBadBlocks returns an empty list on a fresh node. |
| `debug_getBlockRlp/get-genesis.io` | debug_getBlockRlp returns the RLP of the genesis block. |
| `debug_traceTransaction/trace-notfound.io` | debug_traceTransaction returns error for a missing transaction. |
| `eth_accounts/get-accounts.io` | eth_accounts returns the list of owned accounts. |
| `eth_blockNumber/get-block-number.io` | eth_blockNumber returns the current block number at startup. |
| `eth_call/call-empty.io` | eth_call to a non-existent account returns empty result. |
| `eth_chainId/get-chain-id.io` | eth_chainId is not supported by this XDC build and returns method not found. |
| `eth_coinbase/get-coinbase.io` | eth_coinbase returns the configured coinbase / foundation wallet address. |
| `eth_estimateGas/simple-transfer.io` | eth_estimateGas for a simple transfer to a new account. |
| `eth_feeHistory/fee-history.io` | eth_feeHistory is not supported by this XDC build and returns method not found. |
| `eth_gasPrice/get-gas-price.io` | eth_gasPrice returns the current gas price. |
| `eth_getBalance/get-earliest-foundation.io` | eth_getBalance returns the foundation wallet balance at earliest block. |
| `eth_getBalance/get-earliest-genesis-account.io` | eth_getBalance returns a genesis account balance at earliest block. |
| `eth_getBalance/get-foundation-wallet.io` | eth_getBalance for foundation wallet in genesis. |
| `eth_getBalance/get-genesis-account.io` | eth_getBalance for genesis account with max balance. |
| `eth_getBalance/get-latest-03c0.io` | eth_getBalance returns a genesis account balance at latest. |
| `eth_getBalance/get-latest-38b2.io` | eth_getBalance returns a genesis account balance at latest. |
| `eth_getBalance/get-latest-42e6.io` | eth_getBalance returns a genesis account balance at latest. |
| `eth_getBalance/get-latest-4398.io` | eth_getBalance returns a genesis account balance at latest. |
| `eth_getBalance/get-latest-4e37.io` | eth_getBalance returns a genesis account balance at latest. |
| `eth_getBalance/get-latest-664c.io` | eth_getBalance returns a genesis account balance at latest. |
| `eth_getBalance/get-latest-7aa1.io` | eth_getBalance returns a genesis account balance at latest. |
| `eth_getBalance/get-latest-9a37.io` | eth_getBalance returns a genesis account balance at latest. |
| `eth_getBalance/get-latest-account.io` | eth_getBalance for genesis account at latest (still same since no new blocks). |
| `eth_getBalance/get-latest-validator.io` | eth_getBalance returns the validator contract balance at latest. |
| `eth_getBalance/get-unknown-account.io` | eth_getBalance for an unknown account returns zero. |
| `eth_getBalance/invalid-address.io` | eth_getBalance with invalid address format returns error. |
| `eth_getBlockByHash/get-full-genesis.io` | eth_getBlockByHash returns the genesis block with full transactions. |
| `eth_getBlockByHash/get-genesis.io` | eth_getBlockByHash returns the genesis block by its hash. |
| `eth_getBlockByHash/get-notfound.io` | eth_getBlockByHash returns null for an unknown block hash. |
| `eth_getBlockByNumber/get-earliest.io` | eth_getBlockByNumber returns the genesis block by earliest tag. |
| `eth_getBlockByNumber/get-full-genesis.io` | eth_getBlockByNumber returns the genesis block with full transactions. |
| `eth_getBlockByNumber/get-genesis.io` | eth_getBlockByNumber returns the genesis block. |
| `eth_getBlockByNumber/get-latest.io` | eth_getBlockByNumber latest returns the genesis block because no blocks have been mined. |
| `eth_getBlockByNumber/get-notfound.io` | eth_getBlockByNumber returns null for a block that does not exist. |
| `eth_getBlockFinality/get-genesis.io` | eth_getBlockFinality is not supported by this XDC build and returns method not found. |
| `eth_getBlockReceipts/get-genesis.io` | eth_getBlockReceipts is not supported by this XDC build and returns method not found. |
| `eth_getBlockTransactionCountByHash/get-genesis.io` | eth_getBlockTransactionCountByHash returns zero for genesis block. |
| `eth_getBlockTransactionCountByHash/get-notfound.io` | eth_getBlockTransactionCountByHash returns zero for an unknown block hash. |
| `eth_getBlockTransactionCountByNumber/get-genesis.io` | eth_getBlockTransactionCountByNumber returns zero for genesis. |
| `eth_getCandidateStatus/get-empty.io` | eth_getCandidateStatus returns empty string for non-candidate at genesis. |
| `eth_getCode/get-empty-at-genesis.io` | eth_getCode for a non-contract account at genesis returns empty. |
| `eth_getCode/get-unknown-account.io` | eth_getCode for an unknown account returns empty. |
| `eth_getCompensation/get-empty.io` | eth_getCompensation is not supported by this XDC build. |
| `eth_getLogs/get-genesis-range.io` | eth_getLogs returns empty array for genesis block range. |
| `eth_getLogs/invalid-filter.io` | eth_getLogs with invalid filter returns error. |
| `eth_getMasternodeInfo/get-empty.io` | eth_getMasternodeInfo is not supported by this XDC build and returns method not found. |
| `eth_getMasternodes/get-empty.io` | eth_getMasternodes exists but rejects a string block argument on this XDC build. |
| `eth_getProof/get-genesis-account.io` | eth_getProof is not supported by this XDC build and returns method not found. |
| `eth_getRewards/get-empty.io` | eth_getRewards is not supported by this XDC build and returns method not found. |
| `eth_getStorageAt/get-unknown-account.io` | eth_getStorageAt for an unknown account returns zero. |
| `eth_getStorageAt/get-validator-slot-7.io` | eth_getStorageAt reads validator contract slot 0x7. |
| `eth_getStorageAt/get-validator-slot-8.io` | eth_getStorageAt reads validator contract slot 0x8. |
| `eth_getStorageAt/get-validator-slot-9.io` | eth_getStorageAt reads validator contract slot 0x9. |
| `eth_getStorageAt/get-validator-slot-a.io` | eth_getStorageAt reads validator contract slot 0xa. |
| `eth_getStorageAt/get-validator-slot-b.io` | eth_getStorageAt for validator contract storage slot 0xb. |
| `eth_getStorageAt/get-validator-slot-c.io` | eth_getStorageAt for validator contract storage slot 0xc. |
| `eth_getStorageAt/get-validator-slot-d.io` | eth_getStorageAt reads validator contract slot 0xd. |
| `eth_getStorageAt/get-validator-slot-e.io` | eth_getStorageAt reads validator contract slot 0xe. |
| `eth_getStorageAt/get-validator-slot-f.io` | eth_getStorageAt reads validator contract slot 0xf. |
| `eth_getStorageAt/invalid-key.io` | eth_getStorageAt coerces an invalid storage key to zero on this XDC build. |
| `eth_getTransactionByBlockHashAndIndex/get-genesis.io` | eth_getTransactionByBlockHashAndIndex returns null for genesis block with no transactions. |
| `eth_getTransactionByBlockNumberAndIndex/get-genesis.io` | eth_getTransactionByBlockNumberAndIndex returns null for genesis block with no transactions. |
| `eth_getTransactionByHash/get-notfound.io` | eth_getTransactionByHash returns null for a missing transaction hash. |
| `eth_getTransactionCount/get-genesis-account.io` | eth_getTransactionCount returns zero for a genesis account with no transactions. |
| `eth_getTransactionCount/get-unknown-account.io` | eth_getTransactionCount for an unknown account returns zero. |
| `eth_getTransactionCount/get-validator.io` | eth_getTransactionCount returns zero for the validator contract. |
| `eth_getTransactionReceipt/get-notfound.io` | eth_getTransactionReceipt returns null for a missing transaction hash. |
| `eth_getUncleByBlockHashAndIndex/get-genesis.io` | eth_getUncleByBlockHashAndIndex returns null for genesis block. |
| `eth_getUncleByBlockNumberAndIndex/get-genesis.io` | eth_getUncleByBlockNumberAndIndex returns null for genesis block. |
| `eth_getUncleCountByBlockHash/get-genesis.io` | eth_getUncleCountByBlockHash returns zero for genesis block. |
| `eth_getUncleCountByBlockHash/get-notfound.io` | eth_getUncleCountByBlockHash returns zero for an unknown block hash. |
| `eth_getUncleCountByBlockNumber/get-genesis.io` | eth_getUncleCountByBlockNumber returns zero for the genesis block. |
| `eth_getVoters/get-empty.io` | eth_getVoters is not supported by this XDC build and returns method not found. |
| `eth_getWork/get-work.io` | eth_getWork returns mining work or error depending on node state. |
| `eth_hashrate/get-hashrate.io` | eth_hashrate returns the current hashrate. |
| `eth_mining/check-mining.io` | eth_mining returns whether the node is actively mining. |
| `eth_protocolVersion/get-version.io` | eth_protocolVersion returns the current protocol version. |
| `eth_sendRawTransaction/invalid-tx.io` | eth_sendRawTransaction rejects an empty transaction. |
| `eth_sendTransaction/send-no-account.io` | eth_sendTransaction is only supported over IPC on this XDC build. |
| `eth_sign/sign-no-account.io` | eth_sign requires an unlocked account and fails here. |
| `eth_submitHashrate/submit-hashrate.io` | eth_submitHashrate reports hashrate to the node with a 32-byte id. |
| `eth_submitWork/submit-work.io` | eth_submitWork returns false when no work is pending. |
| `eth_syncing/check-syncing.io` | eth_syncing returns false at genesis. |
| `miner_setEtherbase/set-etherbase.io` | miner_setEtherbase returns true for a valid address. |
| `miner_stop/stop-miner.io` | miner_stop returns true and stops the miner. |
| `net_listening/check-listening.io` | net_listening returns true when the node is accepting connections. |
| `net_peerCount/get-peer-count.io` | net_peerCount returns a valid hex number. |
| `net_version/get-network-id.io` | net_version returns the network ID. |
| `personal_listAccounts/list-accounts.io` | personal_listAccounts is not available over HTTP on this XDC build. |
| `rpc_modules/get-modules.io` | rpc_modules lists enabled RPC modules. |
| `txpool_content/get-content.io` | txpool_content returns empty pending and queued maps. |
| `txpool_contentFrom/get-empty.io` | txpool_contentFrom is not supported by this XDC build and returns method not found. |
| `txpool_inspect/get-inspect.io` | txpool_inspect returns empty pending and queued maps. |
| `txpool_status/get-status.io` | txpool_status returns pool counters. |
| `web3_clientVersion/get-client-version.io` | web3_clientVersion returns the client version string. |
| `web3_sha3/get-sha3.io` | web3_sha3 returns the keccak256 hash of the input. |

Last verified: 2026-06-22 — all 102 fixtures pass against xdpos.
