#!/bin/bash
set -e

PAYLOAD='{"jsonrpc":"2.0","method":"admin_nodeInfo","params":[],"id":1}'
RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" --data "$PAYLOAD" http://localhost:8545)
echo "$RESPONSE" | jq -r '.result.enode'
