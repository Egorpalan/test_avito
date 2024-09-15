FROM golang:1.22.3-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/app/main.go

COPY .env ./

EXPOSE 8080

CMD ["./main"]