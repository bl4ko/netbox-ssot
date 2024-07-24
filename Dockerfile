FROM --platform=$BUILDPLATFORM golang:1.22.5@sha256:1b70fa8200ed48367e7d0e0e82c0c1d14cf712d422116f9d69e7c964243653a0 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./internal ./internal

COPY ./cmd ./cmd

RUN CGO_ENABLED=0 GOOS=${TARGET_OS} GOARCH=${TARGETARCH} go build  -o ./cmd/netbox-ssot/main ./cmd/netbox-ssot/main.go

FROM alpine:3.20.2@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5

# Install openssh required for netconf
RUN apk add openssh

WORKDIR /app

COPY --from=builder /app/cmd/netbox-ssot/main ./main

CMD ["./main"]
