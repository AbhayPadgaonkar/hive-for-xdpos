#!/bin/bash

# XDC Gateway Hive client entry point.
# Starts a local go-ethereum node as upstream, then exposes it through
# a minimal JSON-RPC gateway proxy on port 8545.

set -e

# ---- Initialize Geth with the test genesis ----
if [ "$HIVE_LOGLEVEL" != "" ]; then
    VERBOSITY="--verbosity=$HIVE_LOGLEVEL"
else
    VERBOSITY="--verbosity=3"
fi

NETWORK_ID="${HIVE_NETWORK_ID:-1337}"

# Copy genesis to working location and map it to Geth-compatible format.
cp /genesis.json /genesis-input.json
jq -f /mapper.jq /genesis-input.json > /genesis-mapped.json

DATADIR=/gateway-datadir
rm -rf "$DATADIR"
mkdir -p "$DATADIR"

echo "Initializing upstream Geth datadir..."
geth --datadir "$DATADIR" init /genesis-mapped.json

# Optional chain import.
if [ -f /chain.rlp ]; then
    echo "Importing chain.rlp..."
    geth --datadir "$DATADIR" --networkid "$NETWORK_ID" import /chain.rlp || true
fi

if [ -d /blocks ]; then
    echo "Importing individual blocks..."
    (cd /blocks && geth --datadir "$DATADIR" --networkid "$NETWORK_ID" import $(ls | sort -n)) || true
fi

# Create the Hive static JWT secret so Geth's engine API accepts simulator requests.
# Hive uses ENGINEAPI_JWT_SECRET = [32]byte{0x73,0x65,0x63,0x72,0x65,0x74,...}
# i.e. the UTF-8 string "secretsecretsecretsecretsecretsecr" (32 bytes).
# Geth accepts either 32 raw bytes or a 64-character hex string.
JWT_SECRET_FILE=/gateway-jwt-secret
rm -f "$JWT_SECRET_FILE"
# Hex encoding of the 32-byte secret.
printf '%s' '7365637265747365637265747365637265747365637265747365637265747365' > "$JWT_SECRET_FILE"
ls -l "$JWT_SECRET_FILE"

# ---- Start upstream Geth on port 8546 ----
GETH_FLAGS="--datadir $DATADIR --networkid $NETWORK_ID"
GETH_FLAGS="$GETH_FLAGS --syncmode full"
GETH_FLAGS="$GETH_FLAGS --port 30303"
GETH_FLAGS="$GETH_FLAGS --http --http.addr 127.0.0.1 --http.port 8546"
GETH_FLAGS="$GETH_FLAGS --http.api admin,debug,eth,miner,net,txpool,web3,testing"
GETH_FLAGS="$GETH_FLAGS --ws --ws.addr 127.0.0.1 --ws.port 8547 --ws.origins '*'"
GETH_FLAGS="$GETH_FLAGS --ws.api admin,debug,eth,miner,net,txpool,web3,testing"
GETH_FLAGS="$GETH_FLAGS --http.vhosts '*' --http.corsdomain '*'"
GETH_FLAGS="$GETH_FLAGS --authrpc.addr 127.0.0.1 --authrpc.port 8651 --authrpc.vhosts '*'"
GETH_FLAGS="$GETH_FLAGS --authrpc.jwtsecret $JWT_SECRET_FILE"
GETH_FLAGS="$GETH_FLAGS $VERBOSITY"

# Determine external IP for NAT / enode.
ip=$(ip addr show eth0 2>/dev/null | grep 'inet ' | awk '{print $2}' | cut -d/ -f1 | head -1)
if [ -n "$ip" ]; then
    GETH_FLAGS="$GETH_FLAGS --nat extip:$ip"
fi

if [ "$HIVE_MINER" != "" ]; then
    GETH_FLAGS="$GETH_FLAGS --mine --miner.etherbase $HIVE_MINER"
fi

echo "Starting upstream Geth..."
geth $GETH_FLAGS &
GETH_PID=$!

# Wait for upstream RPC to be ready.
echo "Waiting for upstream Geth RPC on 8546..."
for i in $(seq 1 60); do
    if curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        http://127.0.0.1:8546 > /dev/null 2>&1; then
        echo "Upstream Geth ready."
        break
    fi
    if [ $i -eq 60 ]; then
        echo "Upstream Geth failed to start"
        kill $GETH_PID || true
        wait $GETH_PID || true
        exit 1
    fi
    sleep 1
done

# ---- Start gateway proxy on port 8545 ----
echo "Starting XDC Gateway proxy on 0.0.0.0:8545..."
export UPSTREAM_URL=http://127.0.0.1:8546
export ENGINE_URL=http://127.0.0.1:8651
export PORT=8545
cd /app && exec node proxy.js
