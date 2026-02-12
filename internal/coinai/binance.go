package coinai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const defaultBinanceBaseURL = "https://api.binance.com"

type BinanceClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewBinanceClient(baseURL string, timeout time.Duration) *BinanceClient {
	if baseURL == "" {
		baseURL = defaultBinanceBaseURL
	}
	if timeout <= 0 {
		timeout = 15 * time.Second
	}

	return &BinanceClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *BinanceClient) FetchKlines(ctx context.Context, symbol, interval string, limit int) ([]Candle, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}
	if interval == "" {
		return nil, fmt.Errorf("interval is required")
	}
	if limit <= 0 || limit > 1000 {
		return nil, fmt.Errorf("limit must be in range 1..1000")
	}

	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = "/api/v3/klines"
	q := u.Query()
	q.Set("symbol", symbol)
	q.Set("interval", interval)
	q.Set("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request klines: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("binance status %d: %s", resp.StatusCode, string(body))
	}

	var raw [][]any
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode klines: %w", err)
	}

	candles := make([]Candle, 0, len(raw))
	for idx, row := range raw {
		if len(row) < 7 {
			return nil, fmt.Errorf("invalid kline row at index %d", idx)
		}

		openTimeMs, err := parseInt64(row[0])
		if err != nil {
			return nil, fmt.Errorf("parse open time index %d: %w", idx, err)
		}
		open, err := parseFloat(row[1])
		if err != nil {
			return nil, fmt.Errorf("parse open index %d: %w", idx, err)
		}
		high, err := parseFloat(row[2])
		if err != nil {
			return nil, fmt.Errorf("parse high index %d: %w", idx, err)
		}
		low, err := parseFloat(row[3])
		if err != nil {
			return nil, fmt.Errorf("parse low index %d: %w", idx, err)
		}
		closePrice, err := parseFloat(row[4])
		if err != nil {
			return nil, fmt.Errorf("parse close index %d: %w", idx, err)
		}
		volume, err := parseFloat(row[5])
		if err != nil {
			return nil, fmt.Errorf("parse volume index %d: %w", idx, err)
		}
		closeTimeMs, err := parseInt64(row[6])
		if err != nil {
			return nil, fmt.Errorf("parse close time index %d: %w", idx, err)
		}

		candles = append(candles, Candle{
			OpenTime:  time.UnixMilli(openTimeMs).UTC(),
			CloseTime: time.UnixMilli(closeTimeMs).UTC(),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closePrice,
			Volume:    volume,
		})
	}

	return candles, nil
}

func parseFloat(v any) (float64, error) {
	switch t := v.(type) {
	case float64:
		return t, nil
	case string:
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return 0, fmt.Errorf("parse float from string %q: %w", t, err)
		}
		return f, nil
	case json.Number:
		f, err := t.Float64()
		if err != nil {
			return 0, fmt.Errorf("parse float from number %q: %w", t, err)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("unsupported float type %T", v)
	}
}

func parseInt64(v any) (int64, error) {
	switch t := v.(type) {
	case int64:
		return t, nil
	case float64:
		return int64(t), nil
	case string:
		n, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parse int64 from string %q: %w", t, err)
		}
		return n, nil
	case json.Number:
		n, err := t.Int64()
		if err != nil {
			return 0, fmt.Errorf("parse int64 from number %q: %w", t, err)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("unsupported int64 type %T", v)
	}
}
