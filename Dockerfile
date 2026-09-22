FROM --platform=$BUILDPLATFORM golang:1.27.1@sha256:3680233e3204827fbdc66088528ae6d4b3d034f51d03a99d454f6de034888244 AS builder

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

FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6

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
