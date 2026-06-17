# This simulation runs JSON-RPC API tests against XDC-compatible clients.
FROM golang:1-alpine as builder
ARG GOPROXY
ENV GOPROXY=${GOPROXY}

RUN apk add --update gcc musl-dev linux-headers

# Build the simulator executable.
ADD . /source
WORKDIR /source
RUN go build -v -o xdc-rpc-compat .

# Build the simulator run container.
FROM alpine:latest
ADD . /source
WORKDIR /source
COPY --from=builder /source/xdc-rpc-compat .
COPY --from=builder /source/tests ./tests

ENTRYPOINT ["./xdc-rpc-compat"]
