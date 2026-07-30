package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/hive/hivesim"
)

var knownMethods = []string{
	"eth_accounts", "eth_blockNumber", "eth_call", "eth_chainId", "eth_coinbase",
	"eth_createAccessList", "eth_estimateGas", "eth_feeHistory", "eth_gasPrice",
	"eth_getBalance", "eth_getBlockByHash", "eth_getBlockByNumber",
	"eth_getBlockReceipts", "eth_getBlockTransactionCountByHash",
	"eth_getBlockTransactionCountByNumber", "eth_getCode", "eth_getFilterChanges",
	"eth_getFilterLogs", "eth_getLogs", "eth_getProof",
	"eth_getRawTransactionByBlockHashAndIndex", "eth_getRawTransactionByBlockNumberAndIndex",
	"eth_getRawTransactionByHash", "eth_getStorageAt",
	"eth_getTransactionByBlockHashAndIndex", "eth_getTransactionByBlockNumberAndIndex",
	"eth_getTransactionByHash", "eth_getTransactionCount", "eth_getTransactionReceipt",
	"eth_getUncleByBlockHashAndIndex", "eth_getUncleByBlockNumberAndIndex",
	"eth_getUncleCountByBlockHash", "eth_getUncleCountByBlockNumber", "eth_getWork",
	"eth_hashrate", "eth_maxPriorityFeePerGas", "eth_mining", "eth_newBlockFilter",
	"eth_newFilter", "eth_newPendingTransactionFilter", "eth_pendingTransactions",
	"eth_protocolVersion", "eth_sendRawTransaction", "eth_sendTransaction",
	"eth_sign", "eth_signTransaction", "eth_submitHashrate", "eth_submitWork",
	"eth_syncing", "eth_uninstallFilter",
	"eth_blobBaseFee",
	"eth_getBlockFinality", "eth_getBlockFinalityByHash", "eth_getBlockFinalityByNumber",
	"eth_getBlockSignersByHash", "eth_getBlockSignersByNumber",
	"eth_getRewardByHash", "eth_getRewards", "eth_getCompensation",
	"eth_getCandidateStatus", "eth_getMasternodeInfo", "eth_getMasternodes",
	"eth_getVoters",
	"net_listening", "net_peerCount", "net_version",
	"web3_clientVersion", "web3_sha3",
	"admin_addPeer", "admin_addTrustedPeer", "admin_datadir", "admin_exportChain",
	"admin_importChain", "admin_nodeInfo", "admin_peerEvents", "admin_peers",
	"admin_removePeer", "admin_removeTrustedPeer",
	"debug_accountRange", "debug_chaindbCompact", "debug_chaindbProperty",
	"debug_dumpBlock", "debug_getBadBlocks", "debug_getBlockRlp",
	"debug_getModifiedAccountsByHash", "debug_getModifiedAccountsByNumber",
	"debug_preimage", "debug_printBlock", "debug_storageRangeAt", "debug_traceBlock",
	"debug_traceBlockByHash", "debug_traceBlockByNumber", "debug_traceBlockFromFile",
	"debug_traceTransaction",
	"miner_setEtherbase", "miner_setExtra", "miner_setGasPrice", "miner_start", "miner_stop",
	"personal_listAccounts", "personal_importRawKey", "personal_unlockAccount",
	"txpool_content", "txpool_contentFrom", "txpool_inspect", "txpool_status",
	"XDPoS_calculateBlockInfoByV1EpochNum", "XDPoS_config",
	"XDPoS_getBlockInfoByEpochNum", "XDPoS_getBlockInfoByV2EpochNum",
	"XDPoS_getEpochNumbersBetween", "XDPoS_getLatestPoolStatus",
	"XDPoS_getMasternodesByNumber", "XDPoS_getMissedRoundsInEpochByBlockNum",
	"XDPoS_getRewardByAccount", "XDPoS_getSigners", "XDPoS_getSignersAtHash",
	"XDPoS_getSnapshot", "XDPoS_getSnapshotAtHash", "XDPoS_getV2BlockByHash",
	"XDPoS_getV2BlockByHeader", "XDPoS_getV2BlockByNumber", "XDPoS_networkInformation",
	"engine_forkchoiceUpdatedV1", "engine_forkchoiceUpdatedV2", "engine_forkchoiceUpdatedV3",
	"engine_newPayloadV1", "engine_newPayloadV2", "engine_newPayloadV3",
	"engine_getPayloadV1", "engine_getPayloadV2", "engine_getPayloadV3",
	"engine_exchangeTransitionConfigurationV1",
	"rpc_modules",
}

type methodResult struct {
	Method      string `json:"method"`
	Supported   bool   `json:"supported"`
	Error       string `json:"error,omitempty"`
	SampleValue string `json:"sample_value,omitempty"`
}

func main() {
	sim := hivesim.New()
	suite := hivesim.Suite{
		Name:        "rpc-methods",
		Description: "Probes all known RPC methods on a client and reports which are supported.",
	}
	suite.Add(hivesim.ClientTestSpec{
		Role:       "eth1",
		Name:       "rpc method probe",
		Run:        runMethodProbe,
		AlwaysRun:  true,
	})
	hivesim.MustRunSuite(sim, suite)
}

func runMethodProbe(t *hivesim.T, c *hivesim.Client) {
	client, err := rpc.DialHTTP(fmt.Sprintf("http://%s:8545", c.IP))
	if err != nil {
		t.Fatalf("failed to dial RPC: %v", err)
	}
	defer client.Close()

	var version string
	if err := client.Call(&version, "web3_clientVersion"); err == nil {
		t.Logf("Client version: %s", version)
	}

	var modules map[string]string
	if err := client.Call(&modules, "rpc_modules"); err == nil {
		modNames := make([]string, 0, len(modules))
		for k, v := range modules {
			modNames = append(modNames, fmt.Sprintf("%s:%s", k, v))
		}
		t.Logf("RPC modules: %s", strings.Join(modNames, ", "))
	}

	results := probeMethods(t, client)

	supported := 0
	unsupported := 0
	for _, r := range results {
		if r.Supported {
			supported++
		} else {
			unsupported++
		}
	}
	t.Logf("Probe complete: %d supported, %d unsupported out of %d methods",
		supported, unsupported, len(results))

	var unsupportedList []string
	for _, r := range results {
		if !r.Supported {
			unsupportedList = append(unsupportedList, r.Method)
		}
	}
	if len(unsupportedList) > 0 {
		t.Logf("Unsupported methods:\n  %s", strings.Join(unsupportedList, "\n  "))
	}

	out := map[string]interface{}{
		"client":      c.Type,
		"version":     version,
		"modules":     modules,
		"total":       len(results),
		"supported":   supported,
		"unsupported": unsupported,
		"methods":     results,
	}
	data, _ := json.MarshalIndent(out, "", "  ")
	t.Logf("PROBE_RESULT %s", string(data))
}

func probeMethods(t *hivesim.T, client *rpc.Client) []methodResult {
	results := make([]methodResult, 0, len(knownMethods))
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	for _, method := range knownMethods {
		var result interface{}
		err := client.CallContext(ctx, &result, method)
		mr := methodResult{Method: method}
		if err == nil {
			mr.Supported = true
			if result != nil {
				switch v := result.(type) {
				case string:
					if len(v) < 100 {
						mr.SampleValue = v
					}
				case map[string]interface{}, []interface{}:
					if b, e := json.Marshal(v); e == nil && len(b) < 200 {
						mr.SampleValue = string(b)
					}
				default:
					mr.SampleValue = fmt.Sprintf("%v", v)
				}
			}
		} else {
			mr.Supported = false
			errStr := err.Error()
			if len(errStr) > 120 {
				errStr = errStr[:120]
			}
			mr.Error = errStr
		}
		results = append(results, mr)
	}
	return results
}
