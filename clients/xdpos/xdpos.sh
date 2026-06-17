#!/bin/bash

# XDC Core Node startup script for Hive.
# Uses a real XDC genesis (mainnet/testnet) by default.
# If /genesis.json from the simulator looks like a small test genesis, fall back
# to the real XDC testnet genesis.

set -e

XDC_BIN="${XDC_BIN:-$(command -v XDC || true)}"
if [ -z "$XDC_BIN" ]; then
    for path in /usr/bin/XDC /usr/local/bin/XDC /work/xdcchain/XDC /xdc/XDC /opt/xdc/XDC /bin/XDC; do
        if [ -x "$path" ]; then
            XDC_BIN="$path"
            break
        fi
    done
fi
if [ -z "$XDC_BIN" ]; then
    echo "XDC binary not found"
    find / -name "XDC" -type f 2>/dev/null | head -5
    exit 1
fi

echo "Using XDC binary: $XDC_BIN"

FLAGS=""

# Logging verbosity.
if [ "$HIVE_LOGLEVEL" != "" ]; then
    FLAGS="$FLAGS --verbosity $HIVE_LOGLEVEL"
fi

# Network ID / chain ID.
if [ "$HIVE_NETWORK_ID" != "" ]; then
    FLAGS="$FLAGS --networkid $HIVE_NETWORK_ID"
fi

# Node type / sync mode.
case "$HIVE_NODETYPE" in
    "" | full)
        FLAGS="$FLAGS --syncmode full" ;;
    archive)
        FLAGS="$FLAGS --syncmode full --gcmode archive" ;;
    snap)
        FLAGS="$FLAGS --syncmode snap" ;;
esac

# Bootnode.
if [ -n "$HIVE_BOOTNODE" ]; then
    FLAGS="$FLAGS --bootnodes=$HIVE_BOOTNODE"
fi

# Determine which genesis to use.
# If the simulator provided a small test genesis, use the real XDC testnet genesis.
GENESIS=/genesis.json
if [ -f /genesis.json ]; then
    genesis_keys=$(jq 'keys | length' /genesis.json 2>/dev/null || echo 0)
    if [ "$genesis_keys" -lt 5 ]; then
        echo "Simulator genesis is too small, using real XDC testnet genesis."
        GENESIS=/genesis-testnet.json
    fi
else
    GENESIS=/genesis-testnet.json
fi

echo "Using genesis: $GENESIS"

# Initialize datadir.
DATADIR=/xdcdatadir
rm -rf "$DATADIR"
mkdir -p "$DATADIR"

echo "Initializing XDC datadir with genesis..."
$XDC_BIN --datadir "$DATADIR" init "$GENESIS"

# Import chain.rlp if present.
if [ -f /chain.rlp ]; then
    echo "Importing chain.rlp..."
    $XDC_BIN --datadir "$DATADIR" $FLAGS import /chain.rlp || true
fi

# Import individual blocks.
if [ -d /blocks ]; then
    echo "Importing individual blocks..."
    (cd /blocks && $XDC_BIN --datadir "$DATADIR" $FLAGS import $(ls | sort -n)) || true
fi

# JSON-RPC / WS.
FLAGS="$FLAGS --rpc --rpcaddr 0.0.0.0 --rpcport 8545 --rpcapi admin,debug,eth,miner,net,txpool,web3,testing"
FLAGS="$FLAGS --ws --wsaddr 0.0.0.0 --wsport 8546 --wsorigins '*' --wsapi admin,debug,eth,miner,net,txpool,web3,testing"
FLAGS="$FLAGS --rpcvhosts '*'"

# Determine external IP for NAT.
ip=$(ip addr show eth0 2>/dev/null | grep 'inet ' | awk '{print $2}' | cut -d/ -f1 | head -1)
if [ -n "$ip" ]; then
    FLAGS="$FLAGS --nat extip:$ip"
fi

# XDC requires an explicit coinbase/etherbase, especially for XDPoS masternode check.
ETHERBASE="${HIVE_MINER:-xdc0000000000000000000000000000000000000000}"
# XDC expects xdc-prefixed addresses; keep the prefix if the variable already has one.
if ! echo "$ETHERBASE" | grep -qE '^xdc'; then
    ETHERBASE="xdc${ETHERBASE#0x}"
fi
FLAGS="$FLAGS --etherbase $ETHERBASE"

echo "Starting XDC with flags: $FLAGS"
exec $XDC_BIN --datadir "$DATADIR" $FLAGS
