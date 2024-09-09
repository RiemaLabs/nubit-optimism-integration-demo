FROM --platform=$BUILDPLATFORM golang:1.21.5-alpine3.19 as builder

ARG VERSION=nubit-v1.9.0

RUN apk add --no-cache make gcc musl-dev linux-headers git jq bash

# build op-batcher with the shared go.mod & go.sum files
COPY ./op-batcher /app/op-batcher
COPY ./op-nubit /app/op-nubit
COPY ./op-node /app/op-node
COPY ./op-service /app/op-service
COPY ./op-plasma /app/op-plasma
COPY ./op-conductor /app/op-conductor
COPY ./packages /app/packages
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app/op-batcher

RUN go mod download

ARG TARGETOS TARGETARCH

RUN make op-batcher VERSION="$VERSION" GOOS=$TARGETOS GOARCH=$TARGETARCH

FROM alpine:3.19

COPY --from=builder /app/op-batcher/bin/op-batcher /usr/local/bin

CMD ["op-batcher"]
