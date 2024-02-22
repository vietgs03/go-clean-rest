FROM golang:1.17-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]