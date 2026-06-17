#!/bin/bash
# Retrieve the enode URL from the upstream Geth node.

curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"admin_nodeInfo","params":[],"id":1}' \
    http://127.0.0.1:8546 | jq -r '.result.enode'
