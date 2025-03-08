FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o trading-bot ./cmd/trading-bot/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/trading-bot .
COPY configs/ ./configs/

# Add proper entrypoint
ENTRYPOINT ["./trading-bot"]
CMD ["-log-level=info"]