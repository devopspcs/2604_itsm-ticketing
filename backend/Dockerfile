FROM golang:1.18-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git

# Copy go.mod dan go.sum
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata postgresql-client
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=5s --retries=3 CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1
CMD ["./server"]
