## XDC Gateway client for Hive
## This wraps a lightweight RPC gateway that fronts a local go-ethereum node.

# Use the same Geth binary that already works in the go-ethereum Hive client.
FROM hive/clients/go-ethereum:latest AS geth

FROM node:20-alpine

RUN apk add --update bash curl jq

# Copy Geth binary from the official image.
COPY --from=geth /usr/local/bin/geth /usr/local/bin/geth

# Install proxy dependencies.
WORKDIR /app
ADD proxy.js /app/proxy.js
ADD package.json /app/package.json
RUN npm install

# Generate version file.
RUN echo "xdc-gateway-hive-client" > /version.txt

ADD xdc-gateway.sh /xdc-gateway.sh
ADD mapper.jq /mapper.jq
ADD enode.sh /hive-bin/enode.sh
ADD genesis.json /genesis.json
RUN chmod +x /xdc-gateway.sh /hive-bin/enode.sh

EXPOSE 8545 8546 8551 8651 30303 30303/udp

ENTRYPOINT ["/xdc-gateway.sh"]
