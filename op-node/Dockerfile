FROM --platform=$BUILDPLATFORM golang:1.21.5-alpine3.19 as builder

ARG VERSION=nubit-v1.9.0

RUN apk add --no-cache make gcc musl-dev linux-headers git jq bash

# build op-node with the shared go.mod & go.sum files
COPY ./op-node /app/op-node
COPY ./op-nubit /app/op-nubit
COPY ./op-chain-ops /app/op-chain-ops
COPY ./op-service /app/op-service
COPY ./op-plasma /app/op-plasma
COPY ./op-conductor /app/op-conductor
COPY ./packages /app/packages
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app/op-node

RUN go mod download

ARG TARGETOS TARGETARCH

RUN make op-node VERSION="$VERSION" GOOS=$TARGETOS GOARCH=$TARGETARCH

FROM alpine:3.19

COPY --from=builder /app/op-node/bin/op-node /usr/local/bin

CMD ["op-node"]
