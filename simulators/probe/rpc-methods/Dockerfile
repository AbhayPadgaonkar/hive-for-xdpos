FROM golang:1-alpine AS builder
RUN apk --no-cache add gcc musl-dev linux-headers
ADD . /probe
WORKDIR /probe
RUN go build -o rpc-methods .

FROM alpine:latest
COPY --from=builder /probe/rpc-methods /
ENTRYPOINT ["/rpc-methods"]
