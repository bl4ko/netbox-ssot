FROM --platform=$BUILDPLATFORM golang:1.22.4@sha256:969349b8121a56d51c74f4c273ab974c15b3a8ae246a5cffc1df7d28b66cf978 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./internal ./internal

COPY ./cmd ./cmd

RUN CGO_ENABLED=0 GOOS=${TARGET_OS} GOARCH=${TARGETARCH} go build  -o ./cmd/netbox-ssot/main ./cmd/netbox-ssot/main.go

FROM alpine:3.20.0@sha256:77726ef6b57ddf65bb551896826ec38bc3e53f75cdde31354fbffb4f25238ebd

# Install openssh required for netconf
RUN apk add openssh

WORKDIR /app

COPY --from=builder /app/cmd/netbox-ssot/main ./main

CMD ["./main"]
