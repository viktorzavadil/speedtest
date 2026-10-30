FROM --platform=$BUILDPLATFORM golang:1.27.3-alpine3.24 AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build
COPY client/ ./client/
WORKDIR /build/client

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o /build/dist/client ./cmd/main

FROM alpine:3.24

RUN apk upgrade --no-cache && apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /build/dist/client ./

ENTRYPOINT ["./client"]