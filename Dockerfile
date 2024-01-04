FROM golang:1.21.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY  ./pkg ./pkg 

COPY ./main.go ./main.go

RUN go build -o main .

CMD ["/app/main"]