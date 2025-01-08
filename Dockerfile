FROM --platform=$BUILDPLATFORM golang:1.23.4@sha256:7ea4c9dcb2b97ff8ee80a67db3d44f98c8ffa0d191399197007d8459c1453041 AS builder

ARG TARGETOS
ARG TARGETARCH
ARG VERSION
ARG CREATED
ARG COMMIT

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./internal ./internal

COPY ./cmd ./cmd

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
  go build -trimpath -ldflags="-s -w \
  -X 'main.version=$VERSION' \
  -X 'main.commit=$COMMIT' \
  -X 'main.date=$CREATED' \
  " -o ./cmd/netbox-ssot/main ./cmd/netbox-ssot/main.go

FROM alpine:3.21.1@sha256:b97e2a89d0b9e4011bb88c02ddf01c544b8c781acf1f4d559e7c8f12f1047ac3

ARG VERSION
ARG CREATED
ARG COMMIT

LABEL \
  org.opencontainers.image.authors="bl4ko" \
  org.opencontainers.image.created=$CREATED \
  org.opencontainers.image.version=$VERSION \
  org.opencontainers.image.revision=$COMMIT \
  org.opencontainers.image.url="https://github.com/bl4ko/netbox-ssot" \
  org.opencontainers.image.documentation="https://github.com/bl4ko/netbox-ssot/blob/main/README.md" \
  org.opencontainers.image.source="https://github.com/bl4ko/netbox-ssot" \
  org.opencontainers.image.title="Netbox-ssot" \
  org.opencontainers.image.description="Microservice for syncing Netbox with multiple external sources."

# Install openssh required for netconf
RUN apk add --no-cache openssh

# Create a netbox user and group
RUN addgroup -S -g 10001 netbox && \
  adduser -S -u 10001 -G netbox netbox && \
  mkdir -p /app && \
  chown -R netbox:netbox /app
USER netbox:netbox

WORKDIR /app

COPY --from=builder --chown=netbox:netbox /app/cmd/netbox-ssot/main ./main

CMD ["./main"]
