## Build xdc-geth-audit from local source

FROM golang:1.26-alpine AS builder

ARG local_path=xdc-geth-audit

RUN apk add --no-cache gcc musl-dev linux-headers git jq curl

COPY $local_path/go.mod $local_path/go.sum /go-ethereum/
WORKDIR /go-ethereum
RUN go mod download

COPY $local_path /go-ethereum
WORKDIR /go-ethereum
RUN go run build/ci.go install -static ./cmd/geth

FROM alpine:latest

RUN apk add --no-cache ca-certificates jq curl bash

COPY --from=builder /go-ethereum/build/bin/geth /usr/local/bin/geth

RUN geth version | head -1 > /version.txt

ADD geth.sh /geth.sh
ADD enode.sh /hive-bin/enode.sh
ADD genesis.json /genesis.json
ADD genesis-testnet.json /genesis-testnet.json
RUN chmod +x /geth.sh /hive-bin/enode.sh

EXPOSE 8545 8546 8547 8551 30303 30303/udp

ENTRYPOINT ["/geth.sh"]
