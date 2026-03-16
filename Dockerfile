FROM --platform=$BUILDPLATFORM golang:1.25.8@sha256:bd1e2df4e6259b2bd5b1de0e6b22ca414502cd6e7276a5dd5dd414b65063be58 AS builder

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

FROM alpine:3.21.6@sha256:c3f8e73fdb79deaebaa2037150150191b9dcbfba68b4a46d70103204c53f4709

ARG VERSION
ARG CREATED
ARG COMMIT

LABEL \
  org.opencontainers.image.authors="src-doo" \
  org.opencontainers.image.created=$CREATED \
  org.opencontainers.image.version=$VERSION \
  org.opencontainers.image.revision=$COMMIT \
  org.opencontainers.image.url="https://github.com/src-doo/netbox-ssot" \
  org.opencontainers.image.documentation="https://github.com/src-doo/netbox-ssot/blob/main/README.md" \
  org.opencontainers.image.source="https://github.com/src-doo/netbox-ssot" \
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

# Also allow deprecated ssh algorithims for older devices
# See https://github.com/SRC-doo/netbox-ssot/issues/498
RUN mkdir -p /home/netbox/.ssh/ && \
cat <<EOF > /home/netbox/.ssh/config
Host *
  HostKeyAlgorithms +ssh-rsa
  PubkeyAcceptedKeyTypes +ssh-rsa
EOF

WORKDIR /app

COPY --from=builder --chown=netbox:netbox /app/cmd/netbox-ssot/main ./main

CMD ["./main"]
