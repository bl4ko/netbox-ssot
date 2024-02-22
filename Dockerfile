FROM --platform=$BUILDPLATFORM golang:1.22.0@sha256:7b297d9abee021bab9046e492506b3c2da8a3722cbf301653186545ecc1e00bb as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./internal ./internal

COPY ./cmd ./cmd

RUN CGO_ENABLED=0 GOOS=${TARGET_OS} GOARCH=${TARGETARCH} go build  -o ./cmd/netbox-ssot/main ./cmd/netbox-ssot/main.go

FROM alpine:3.19.1

WORKDIR /app

COPY --from=builder /app/cmd/netbox-ssot/main ./main

CMD ["./main"]
