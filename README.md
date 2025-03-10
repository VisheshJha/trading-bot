# Upstox Algorithmic Trading Bot

A robust algorithmic trading system implementing EMA crossover strategy on Upstox platform.

## Features

- Real-time market data streaming
- EMA (9/21) crossover strategy
- Automatic order execution
- Risk management system
- Docker containerization
- Paper trading support

## Prerequisites

- Go 1.23+
- Docker 20.10+
- Upstox API credentials
- Basic understanding of algorithmic trading

## Quick Start

```bash
# Clone repository
git clone https://github.com/yourusername/trading-bot.git
cd trading-bot

# Build and run
docker-compose up --build

# Local run
go build -o trading-bot ./cmd/trading-bot/
./trading-bot -config ./configs/config.yaml

# Docker run
# Build and run
docker-compose up --build

# View logs
docker-compose logs -f

# For logs
docker-compose run trading-bot -log-level=debug


# down and up
docker compose down
docker rmi trading-bot --force
docker compose build
docker compose up 
docker compose logs -f

