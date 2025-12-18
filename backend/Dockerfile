FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod first for caching
COPY go.mod ./
RUN go mod download

# Copy source and build
COPY . .
RUN go build -o server .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
