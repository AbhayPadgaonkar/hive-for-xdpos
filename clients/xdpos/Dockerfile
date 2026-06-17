## XDC Core Node client for Hive
## This wraps the official XDC testnet/mainnet Docker image.

ARG baseimage=xinfinorg/xinfin-testnet-node
ARG tag=apothem_network

FROM $baseimage:$tag

# Install jq via static binary since apt repos in the old base image may be broken.
RUN curl -fsSL -o /usr/local/bin/jq https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64 \
    && chmod +x /usr/local/bin/jq

# Add startup script, genesis mapper, real XDC genesis files, and enode retriever.
ADD xdpos.sh /xdpos.sh
ADD mapper.jq /mapper.jq
ADD genesis.json /genesis.json
ADD genesis-mainnet.json /genesis-mainnet.json
ADD genesis-testnet.json /genesis-testnet.json
ADD enode.sh /hive-bin/enode.sh
RUN chmod +x /xdpos.sh /hive-bin/enode.sh

# Generate version.txt if possible.
RUN (/usr/bin/XDC version 2>/dev/null || echo "unknown") | head -1 > /version.txt

EXPOSE 8545 8546 8547 8551 30303 30303/udp

ENTRYPOINT ["/xdpos.sh"]
