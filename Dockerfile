FROM golang:1.25.1 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o app ./cmd/api
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /app/app .


RUN chmod +x /app/app

EXPOSE 8080

CMD ["./app"]
