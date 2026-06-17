# XDC Core Node client for Hive

This directory wraps the official XDC Core Node (XDPoS) so it can be tested by Hive.

## Files

- `Dockerfile` — wraps `xinfinorg/xinfin-testnet-node:latest`
- `Dockerfile.git` — builds XDC Core Node from source
- `hive.yaml` — declares the `eth1` role
- `xdpos.sh` — container entry point; initializes genesis and starts RPC
- `mapper.jq` — converts Geth genesis to XDC genesis format
- `enode.sh` — returns the running node's `enode://` URL
- `genesis.json` — default genesis template

## How to test

```bash
# Smoke tests
.\hive --% -sim smoke/genesis -client xdpos
.\hive --% -sim smoke/genesis -client go-ethereum,xdpos

# RPC compatibility
.\hive --% -sim ethereum/rpc-compat -client go-ethereum,xdpos

# Sync test
.\hive --% -sim ethereum/sync -client go-ethereum,xdpos
```

## Known limitations

- XDPoS consensus (masternodes, epochs, penalties) is not exercised by standard
  Ethereum simulators. For XDPoS-specific behavior a custom simulator is needed.
- Some Ethereum consensus tests may fail because XDC's fork schedule and
  consensus rules differ from mainnet Ethereum.
