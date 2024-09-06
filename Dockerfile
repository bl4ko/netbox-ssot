FROM --platform=$BUILDPLATFORM golang:1.23.1@sha256:4a3c2bcd243d3dbb7b15237eecb0792db3614900037998c2cd6a579c46888c1e as builder

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
