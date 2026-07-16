FROM golang:alpine AS builder
WORKDIR /app
COPY . .
WORKDIR /app/user-Service
RUN go build -o /app/bin/server ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/server .
COPY user-Service/migrations ./migrations
CMD ["./server"]
