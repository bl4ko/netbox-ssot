FROM golang:1.21.4 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./internal ./internal 

COPY ./cmd ./cmd

RUN go build -o ./cmd/netbox-ssot/main ./cmd/netbox-ssot/main.go 

FROM alpine:latest 

WORKDIR /app

COPY --from=builder /app/cmd/netbox-ssot/main ./main

CMD ["./main"]