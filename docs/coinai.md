# Coin AI CLI (Go)

`cmd/coinai` is a research CLI that:
- fetches Binance spot klines,
- trains a baseline linear model for next-candle return,
- runs a simple long/short backtest,
- outputs a `BUY` / `SELL` / `HOLD` signal.

## Quick Start

```bash
go run cmd/coinai/main.go
```

Default settings:
- `symbol=BTCUSDT`
- `interval=1h`
- `limit=500`
- `train-ratio=0.7`

## Custom Run

```bash
go run cmd/coinai/main.go \
  -symbol ETHUSDT \
  -interval 15m \
  -limit 800 \
  -epochs 1200 \
  -lr 0.02 \
  -long-threshold 0.0012 \
  -short-threshold -0.0012 \
  -fee-bps 4
```

## JSON Report

```bash
go run cmd/coinai/main.go -json
```

## Save Trained Model

```bash
go run cmd/coinai/main.go -model-out tmp/eth_model.json
```

## Notes

- This is a baseline for research, not a production trading system.
- Add walk-forward validation, risk controls, and paper-trading before live capital.
