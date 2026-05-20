FROM --platform=$BUILDPLATFORM golang:1.26.3@sha256:cc9a5d7a008cfe2cbc7ffc752b0d6636ad30fc16e4a648d2e4aac00fd8b25ca3 AS builder

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

FROM alpine:3.23.4@sha256:5b10f432ef3da1b8d4c7eb6c487f2f5a8f096bc91145e68878dd4a5019afde11

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
