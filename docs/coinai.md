# Coin AI CLI (Go)

`cmd/coinai` is a research CLI that:
- fetches Binance spot klines for crypto (`market=coin`),
- or loads private stock candles from local CSV (`market=stock`),
- trains a baseline linear model for next-candle return,
- runs a simple long/short backtest,
- outputs a `BUY` / `SELL` / `HOLD` signal.

## Quick Start

```bash
go run cmd/coinai/main.go
```

Default settings:
- `market=coin`
- `symbol=BTCUSDT`
- `interval=1h`
- `limit=500`
- `train-ratio=0.7`

## Custom Run (Coin / Binance)

```bash
go run cmd/coinai/main.go \
  -market coin \
  -symbol ETHUSDT \
  -interval 15m \
  -limit 800 \
  -epochs 1200 \
  -lr 0.02 \
  -long-threshold 0.0012 \
  -short-threshold -0.0012 \
  -fee-bps 4
```

## Custom Run (Stock / Private CSV)

```bash
go run cmd/coinai/main.go \
  -market stock \
  -symbol AAPL \
  -interval 1d \
  -stock-csv ./data/aapl_1d.csv \
  -limit 1200 \
  -json
```

### CSV Format

Required columns:
- time column: `time` or `timestamp` or `date`
- OHLC columns: `open`, `high`, `low`, `close`

Optional columns:
- `volume`
- `close_time`

Example:

```csv
date,open,high,low,close,volume
2026-01-01,210.5,214.0,209.7,213.2,1023400
2026-01-02,213.2,216.4,212.8,215.9,1156700
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
