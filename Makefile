.PHONY: build run deploy clean

build:
    @go build -o bin/trading-bot ./cmd/trading-bot/

run: build
    @./bin/trading-bot

deploy:
    @docker compose up --build -d

clean:
    @docker compose down
    @rm -rf bin/