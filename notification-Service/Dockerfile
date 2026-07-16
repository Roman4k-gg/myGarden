FROM golang:alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/notification-Service
RUN go build -o /app/bin/server ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/server .
CMD ["./server"]
