FROM --platform=$BUILDPLATFORM golang:1.26.4@sha256:792443b89f65105abba56b9bd5e97f680a80074ac62fc844a584212f8c8102c3 AS builder

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

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

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

# Upgrade base packages to pick up security fixes and install openssh for netconf
RUN apk upgrade --no-cache && apk add --no-cache openssh

# Create a netbox user and group
RUN addgroup -S -g 10001 netbox && \
  adduser -S -u 10001 -G netbox netbox && \
  mkdir -p /app && \
  chown -R netbox:netbox /app
USER netbox:netbox

# Also allow deprecated ssh algorithims for older devices
# See https://github.com/bl4ko/netbox-ssot/issues/498
RUN mkdir -p /home/netbox/.ssh/ && \
cat <<EOF > /home/netbox/.ssh/config
Host *
  HostKeyAlgorithms +ssh-rsa
  PubkeyAcceptedKeyTypes +ssh-rsa
EOF

WORKDIR /app

COPY --from=builder --chown=netbox:netbox /app/cmd/netbox-ssot/main ./main

CMD ["./main"]
