FROM golang:1.16-alpine

RUN apk add --no-cache git
RUN go get -u -d github.com/golang-migrate/migrate/cmd/migrate

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]