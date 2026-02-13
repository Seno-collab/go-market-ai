package coinai

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type csvColumns struct {
	openTime  int
	closeTime int
	open      int
	high      int
	low       int
	close     int
	volume    int
}

var csvTimeLayouts = []string{
	time.RFC3339,
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02",
}

func LoadCandlesFromCSV(path string, limit int) ([]Candle, error) {
	if path == "" {
		return nil, fmt.Errorf("csv path is required")
	}
	if limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open csv: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read csv header: %w", err)
	}
	columns, err := resolveCSVColumns(header)
	if err != nil {
		return nil, err
	}

	candles := make([]Candle, 0, 256)
	for rowNum := 2; ; rowNum++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read csv row %d: %w", rowNum, err)
		}
		if isEmptyCSVRecord(record) {
			continue
		}

		candle, err := parseCSVCandle(record, columns)
		if err != nil {
			return nil, fmt.Errorf("parse csv row %d: %w", rowNum, err)
		}
		candles = append(candles, candle)
	}

	if len(candles) == 0 {
		return nil, fmt.Errorf("csv has no candle rows")
	}

	sort.SliceStable(candles, func(i, j int) bool {
		return candles[i].OpenTime.Before(candles[j].OpenTime)
	})

	if limit > 0 && len(candles) > limit {
		candles = append([]Candle(nil), candles[len(candles)-limit:]...)
	}

	return candles, nil
}

func parseCSVCandle(record []string, columns csvColumns) (Candle, error) {
	openTime, err := parseCSVTime(csvValue(record, columns.openTime))
	if err != nil {
		return Candle{}, fmt.Errorf("parse time: %w", err)
	}
	closeTime := openTime
	if columns.closeTime >= 0 {
		closeTime, err = parseCSVTime(csvValue(record, columns.closeTime))
		if err != nil {
			return Candle{}, fmt.Errorf("parse close time: %w", err)
		}
	}

	open, err := parseCSVFloat(csvValue(record, columns.open))
	if err != nil {
		return Candle{}, fmt.Errorf("parse open: %w", err)
	}
	high, err := parseCSVFloat(csvValue(record, columns.high))
	if err != nil {
		return Candle{}, fmt.Errorf("parse high: %w", err)
	}
	low, err := parseCSVFloat(csvValue(record, columns.low))
	if err != nil {
		return Candle{}, fmt.Errorf("parse low: %w", err)
	}
	closePrice, err := parseCSVFloat(csvValue(record, columns.close))
	if err != nil {
		return Candle{}, fmt.Errorf("parse close: %w", err)
	}

	volume := 0.0
	if columns.volume >= 0 {
		volumeValue := strings.TrimSpace(csvValue(record, columns.volume))
		if volumeValue != "" {
			volume, err = parseCSVFloat(volumeValue)
			if err != nil {
				return Candle{}, fmt.Errorf("parse volume: %w", err)
			}
		}
	}

	return Candle{
		OpenTime:  openTime.UTC(),
		CloseTime: closeTime.UTC(),
		Open:      open,
		High:      high,
		Low:       low,
		Close:     closePrice,
		Volume:    volume,
	}, nil
}

func resolveCSVColumns(header []string) (csvColumns, error) {
	columns := csvColumns{
		openTime:  -1,
		closeTime: -1,
		open:      -1,
		high:      -1,
		low:       -1,
		close:     -1,
		volume:    -1,
	}

	for idx, raw := range header {
		name := normalizeCSVHeader(raw)
		switch name {
		case "time", "timestamp", "datetime", "date", "opentime", "openat":
			if columns.openTime == -1 {
				columns.openTime = idx
			}
		case "closetime", "closetimestamp":
			if columns.closeTime == -1 {
				columns.closeTime = idx
			}
		case "open", "o":
			if columns.open == -1 {
				columns.open = idx
			}
		case "high", "h":
			if columns.high == -1 {
				columns.high = idx
			}
		case "low", "l":
			if columns.low == -1 {
				columns.low = idx
			}
		case "close", "c":
			if columns.close == -1 {
				columns.close = idx
			}
		case "volume", "vol", "v":
			if columns.volume == -1 {
				columns.volume = idx
			}
		}
	}

	if columns.openTime == -1 {
		return csvColumns{}, fmt.Errorf("missing time column (time/timestamp/date)")
	}
	if columns.open == -1 || columns.high == -1 || columns.low == -1 || columns.close == -1 {
		return csvColumns{}, fmt.Errorf("missing required OHLC columns")
	}
	return columns, nil
}

func normalizeCSVHeader(name string) string {
	replacer := strings.NewReplacer(" ", "", "_", "", "-", "")
	return strings.ToLower(replacer.Replace(strings.TrimSpace(name)))
}

func parseCSVTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, fmt.Errorf("empty time value")
	}

	if n, err := strconv.ParseInt(value, 10, 64); err == nil {
		switch {
		case n > 1_000_000_000_000:
			return time.UnixMilli(n), nil
		case n > 1_000_000_000:
			return time.Unix(n, 0), nil
		}
	}

	for _, layout := range csvTimeLayouts {
		ts, err := time.Parse(layout, value)
		if err == nil {
			return ts, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported time format %q", value)
}

func parseCSVFloat(value string) (float64, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, fmt.Errorf("empty numeric value")
	}
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("parse float %q: %w", value, err)
	}
	return num, nil
}

func csvValue(record []string, idx int) string {
	if idx < 0 || idx >= len(record) {
		return ""
	}
	return record[idx]
}

func isEmptyCSVRecord(record []string) bool {
	for _, field := range record {
		if strings.TrimSpace(field) != "" {
			return false
		}
	}
	return true
}
