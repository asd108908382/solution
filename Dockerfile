FROM ubuntu:latest
FROM golang:1.20 AS builder
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /app/solution cmd/main.go
RUN ls -la /app
FROM scratch AS deploy
COPY --from=builder /app/solution /solution
CMD ["/solution"]