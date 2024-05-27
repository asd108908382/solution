FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o solution cmd/main.go
FROM alpine:3.15
COPY --from=builder /app/solution /usr/local/bin/solution

EXPOSE 8080

CMD ["solution"]