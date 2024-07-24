FROM --platform=$BUILDPLATFORM golang:1.22.5@sha256:86a3c48a61915a8c62c0e1d7594730399caa3feb73655dfe96c7bc17710e96cf as builder

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
