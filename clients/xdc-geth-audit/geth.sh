#!/bin/bash
set -e

GETH_BIN="${GETH_BIN:-$(command -v geth || true)}"
if [ -z "$GETH_BIN" ]; then
    for path in /usr/local/bin/geth /usr/bin/geth /bin/geth; do
        if [ -x "$path" ]; then
            GETH_BIN="$path"
            break
        fi
    done
fi
if [ -z "$GETH_BIN" ]; then
    echo "geth binary not found"
    find / -name "geth" -type f 2>/dev/null | head -5
    exit 1
fi

echo "Using geth binary: $GETH_BIN"

FLAGS=""

if [ "$HIVE_LOGLEVEL" != "" ]; then
    FLAGS="$FLAGS --verbosity $HIVE_LOGLEVEL"
fi

if [ "$HIVE_NETWORK_ID" != "" ]; then
    FLAGS="$FLAGS --networkid $HIVE_NETWORK_ID"
fi

case "$HIVE_NODETYPE" in
    "" | full)
        FLAGS="$FLAGS --syncmode full" ;;
    archive)
        FLAGS="$FLAGS --syncmode full --gcmode archive" ;;
    snap)
        FLAGS="$FLAGS --syncmode snap" ;;
esac

if [ -n "$HIVE_BOOTNODE" ]; then
    FLAGS="$FLAGS --bootnodes=$HIVE_BOOTNODE"
fi

# Determine genesis: use --apothem unless /genesis.json has XDPoS config
if [ -f /genesis.json ]; then
    has_xdpos=$(jq '.config.XDPoS != null' /genesis.json 2>/dev/null || echo false)
    if [ "$has_xdpos" != "true" ]; then
        echo "Simulator genesis lacks XDPoS config, using built-in Apothem genesis."
        FLAGS="$FLAGS --apothem"
    fi
else
    FLAGS="$FLAGS --apothem"
fi

DATADIR=/xdcdatadir
rm -rf "$DATADIR"
mkdir -p "$DATADIR"

# Initialize only if using a custom genesis file (not --apothem)
if [ -f /genesis.json ] && ! echo "$FLAGS" | grep -q -- "--apothem"; then
    echo "Initializing datadir with genesis..."
    $GETH_BIN --datadir "$DATADIR" init /genesis.json
fi

if [ -f /chain.rlp ]; then
    echo "Importing chain.rlp..."
    $GETH_BIN --datadir "$DATADIR" $FLAGS import /chain.rlp || true
fi

if [ -d /blocks ]; then
    echo "Importing individual blocks..."
    (cd /blocks && $GETH_BIN --datadir "$DATADIR" $FLAGS import $(ls | sort -n)) || true
fi

ip=$(ip addr show eth0 2>/dev/null | grep 'inet ' | awk '{print $2}' | cut -d/ -f1 | head -1)
if [ -n "$ip" ]; then
    FLAGS="$FLAGS --nat extip:$ip"
fi

FLAGS="$FLAGS --http --http.addr 0.0.0.0 --http.port 8545 --http.api admin,debug,eth,miner,net,txpool,web3,XDPoS"
FLAGS="$FLAGS --http.corsdomain '*' --http.vhosts '*'"
FLAGS="$FLAGS --ws --ws.addr 0.0.0.0 --ws.port 8546 --ws.origins '*' --ws.api admin,debug,eth,miner,net,txpool,web3,XDPoS"

ETHERBASE="${HIVE_MINER:-0x0000000000000000000000000000000000000000}"
if [ "$HIVE_MINER" != "" ]; then
    FLAGS="$FLAGS --mine --miner.etherbase $ETHERBASE"
fi

echo "Starting geth with flags: $FLAGS"
exec $GETH_BIN --datadir "$DATADIR" $FLAGS
