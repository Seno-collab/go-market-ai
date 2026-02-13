package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go-ai/internal/coinai"
	"log"
	"os"
	"strings"
	"time"
)

const (
	marketCoin  = "coin"
	marketStock = "stock"
)

type config struct {
	Market         string
	Symbol         string
	Interval       string
	StockCSV       string
	Limit          int
	TrainRatio     float64
	Epochs         int
	LearningRate   float64
	L2             float64
	LongThreshold  float64
	ShortThreshold float64
	FeeBPS         float64
	Timeout        time.Duration
	JSONOutput     bool
	ModelOut       string
}

type trainReport struct {
	Market              string                `json:"market"`
	DataSource          string                `json:"data_source"`
	Symbol              string                `json:"symbol"`
	Interval            string                `json:"interval"`
	Candles             int                   `json:"candles"`
	TrainSamples        int                   `json:"train_samples"`
	TestSamples         int                   `json:"test_samples"`
	FeatureNames        []string              `json:"feature_names"`
	TrainLoss           float64               `json:"train_loss"`
	TestMSE             float64               `json:"test_mse"`
	TestDirectionalAcc  float64               `json:"test_directional_acc"`
	Backtest            coinai.BacktestResult `json:"backtest"`
	NextPredictedReturn float64               `json:"next_predicted_return"`
	Signal              coinai.Signal         `json:"signal"`
	GeneratedAt         time.Time             `json:"generated_at"`
}

type savedModel struct {
	Market       string                `json:"market"`
	DataSource   string                `json:"data_source"`
	Symbol       string                `json:"symbol"`
	Interval     string                `json:"interval"`
	FeatureNames []string              `json:"feature_names"`
	Scaler       coinai.StandardScaler `json:"scaler"`
	Model        coinai.LinearModel    `json:"model"`
	TrainedAt    time.Time             `json:"trained_at"`
}

func main() {
	cfg := parseFlags()
	if err := validateConfig(cfg); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	candles, dataSource, err := loadCandles(ctx, cfg)
	if err != nil {
		log.Fatalf("fetch candles: %v", err)
	}

	samples, err := coinai.BuildDataset(candles)
	if err != nil {
		log.Fatalf("build dataset: %v", err)
	}

	trainSamples, testSamples, err := coinai.SplitSequential(samples, cfg.TrainRatio)
	if err != nil {
		log.Fatalf("split dataset: %v", err)
	}

	trainX, trainY := samplesToXY(trainSamples)
	testX, testY := samplesToXY(testSamples)

	scaler := coinai.NewStandardScaler(len(trainX[0]))
	if err := scaler.Fit(trainX); err != nil {
		log.Fatalf("fit scaler: %v", err)
	}
	trainXNorm, err := scaler.TransformBatch(trainX)
	if err != nil {
		log.Fatalf("normalize train data: %v", err)
	}
	testXNorm, err := scaler.TransformBatch(testX)
	if err != nil {
		log.Fatalf("normalize test data: %v", err)
	}

	model := coinai.NewLinearModel(len(trainXNorm[0]))
	stats, err := model.Train(trainXNorm, trainY, coinai.TrainConfig{
		Epochs:       cfg.Epochs,
		LearningRate: cfg.LearningRate,
		L2:           cfg.L2,
	})
	if err != nil {
		log.Fatalf("train model: %v", err)
	}

	preds := model.PredictBatch(testXNorm)
	testMSE := coinai.MeanSquaredError(preds, testY)
	acc := coinai.DirectionalAccuracy(preds, testY)

	backtest, err := coinai.Backtest(preds, testY, coinai.BacktestConfig{
		LongThreshold:  cfg.LongThreshold,
		ShortThreshold: cfg.ShortThreshold,
		FeeRate:        cfg.FeeBPS / 10000,
	})
	if err != nil {
		log.Fatalf("backtest: %v", err)
	}

	latestFeatures, err := coinai.BuildLatestFeatures(candles)
	if err != nil {
		log.Fatalf("latest features: %v", err)
	}
	latestNorm, err := scaler.Transform(latestFeatures)
	if err != nil {
		log.Fatalf("normalize latest feature: %v", err)
	}
	nextPred := model.Predict(latestNorm)
	signal := coinai.SignalFromPrediction(nextPred, cfg.LongThreshold, cfg.ShortThreshold)

	report := trainReport{
		Market:              normalizeMarket(cfg.Market),
		DataSource:          dataSource,
		Symbol:              cfg.Symbol,
		Interval:            cfg.Interval,
		Candles:             len(candles),
		TrainSamples:        len(trainSamples),
		TestSamples:         len(testSamples),
		FeatureNames:        coinai.FeatureNames(),
		TrainLoss:           stats.FinalLoss,
		TestMSE:             testMSE,
		TestDirectionalAcc:  acc,
		Backtest:            backtest,
		NextPredictedReturn: nextPred,
		Signal:              signal,
		GeneratedAt:         time.Now().UTC(),
	}

	if cfg.ModelOut != "" {
		if err := saveModel(cfg, dataSource, scaler, model); err != nil {
			log.Fatalf("save model: %v", err)
		}
	}

	if cfg.JSONOutput {
		output, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			log.Fatalf("marshal report: %v", err)
		}
		fmt.Println(string(output))
		return
	}

	printReport(report, cfg.ModelOut)
}

func parseFlags() config {
	cfg := config{}

	flag.StringVar(&cfg.Market, "market", marketCoin, "market type: coin | stock")
	flag.StringVar(&cfg.Symbol, "symbol", "BTCUSDT", "trading pair symbol")
	flag.StringVar(&cfg.Interval, "interval", "1h", "candle interval label (e.g. 15m, 1h, 1d)")
	flag.StringVar(&cfg.StockCSV, "stock-csv", "", "CSV path for stock OHLCV data when market=stock")
	flag.IntVar(&cfg.Limit, "limit", 500, "number of latest candles to use")
	flag.Float64Var(&cfg.TrainRatio, "train-ratio", 0.7, "sequential train split ratio")
	flag.IntVar(&cfg.Epochs, "epochs", 800, "training epochs")
	flag.Float64Var(&cfg.LearningRate, "lr", 0.03, "learning rate")
	flag.Float64Var(&cfg.L2, "l2", 0.001, "L2 regularization")
	flag.Float64Var(&cfg.LongThreshold, "long-threshold", 0.0015, "predicted return threshold for BUY")
	flag.Float64Var(&cfg.ShortThreshold, "short-threshold", -0.0015, "predicted return threshold for SELL")
	flag.Float64Var(&cfg.FeeBPS, "fee-bps", 4, "transaction fee in basis points")
	flag.DurationVar(&cfg.Timeout, "timeout", 20*time.Second, "network timeout")
	flag.BoolVar(&cfg.JSONOutput, "json", false, "print output as JSON")
	flag.StringVar(&cfg.ModelOut, "model-out", "", "optional file path to save trained model JSON")

	flag.Parse()
	return cfg
}

func validateConfig(cfg config) error {
	market := normalizeMarket(cfg.Market)
	switch {
	case market != marketCoin && market != marketStock:
		return fmt.Errorf("market must be coin or stock")
	case cfg.Symbol == "":
		return fmt.Errorf("symbol is required")
	case cfg.Interval == "":
		return fmt.Errorf("interval is required")
	case cfg.Limit <= 0:
		return fmt.Errorf("limit must be greater than 0")
	case market == marketCoin && cfg.Limit > 1000:
		return fmt.Errorf("limit for coin must be in range 1..1000")
	case market == marketStock && strings.TrimSpace(cfg.StockCSV) == "":
		return fmt.Errorf("stock-csv is required when market=stock")
	case cfg.TrainRatio <= 0 || cfg.TrainRatio >= 1:
		return fmt.Errorf("train-ratio must be in (0,1)")
	case cfg.LongThreshold <= cfg.ShortThreshold:
		return fmt.Errorf("long-threshold must be greater than short-threshold")
	case cfg.FeeBPS < 0:
		return fmt.Errorf("fee-bps cannot be negative")
	}
	return nil
}

func loadCandles(ctx context.Context, cfg config) ([]coinai.Candle, string, error) {
	switch normalizeMarket(cfg.Market) {
	case marketCoin:
		baseURL := os.Getenv("BINANCE_BASE_URL")
		client := coinai.NewBinanceClient(baseURL, cfg.Timeout)
		candles, err := client.FetchKlines(ctx, cfg.Symbol, cfg.Interval, cfg.Limit)
		if err != nil {
			return nil, "", err
		}
		return candles, "binance", nil
	case marketStock:
		candles, err := coinai.LoadCandlesFromCSV(strings.TrimSpace(cfg.StockCSV), cfg.Limit)
		if err != nil {
			return nil, "", err
		}
		return candles, "csv", nil
	default:
		return nil, "", fmt.Errorf("unsupported market %q", cfg.Market)
	}
}

func normalizeMarket(market string) string {
	return strings.ToLower(strings.TrimSpace(market))
}

func samplesToXY(samples []coinai.Sample) ([][]float64, []float64) {
	x := make([][]float64, 0, len(samples))
	y := make([]float64, 0, len(samples))
	for _, sample := range samples {
		x = append(x, sample.Features)
		y = append(y, sample.Target)
	}
	return x, y
}

func printReport(report trainReport, modelPath string) {
	fmt.Printf("Coin AI report [%s | %s %s]\n", report.Market, report.Symbol, report.Interval)
	fmt.Printf("Data source: %s\n", report.DataSource)
	fmt.Printf("Candles: %d | train: %d | test: %d\n", report.Candles, report.TrainSamples, report.TestSamples)
	fmt.Printf("Train loss: %.8f\n", report.TrainLoss)
	fmt.Printf("Test MSE: %.8f\n", report.TestMSE)
	fmt.Printf("Directional accuracy: %.2f%%\n", report.TestDirectionalAcc*100)
	fmt.Printf("Backtest total return: %.2f%%\n", report.Backtest.TotalReturn*100)
	fmt.Printf("Backtest win rate: %.2f%%\n", report.Backtest.WinRate*100)
	fmt.Printf("Backtest max drawdown: %.2f%%\n", report.Backtest.MaxDrawdown*100)
	fmt.Printf("Backtest sharpe: %.3f\n", report.Backtest.Sharpe)
	fmt.Printf("Backtest trades: %d\n", report.Backtest.Trades)
	fmt.Printf("Predicted next return: %.4f%%\n", report.NextPredictedReturn*100)
	fmt.Printf("Signal: %s\n", report.Signal)
	if modelPath != "" {
		fmt.Printf("Model saved to: %s\n", modelPath)
	}
}

func saveModel(cfg config, dataSource string, scaler *coinai.StandardScaler, model *coinai.LinearModel) error {
	payload := savedModel{
		Market:       normalizeMarket(cfg.Market),
		DataSource:   dataSource,
		Symbol:       cfg.Symbol,
		Interval:     cfg.Interval,
		FeatureNames: coinai.FeatureNames(),
		Scaler:       *scaler,
		Model:        *model,
		TrainedAt:    time.Now().UTC(),
	}
	bytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal model: %w", err)
	}
	if err := os.WriteFile(cfg.ModelOut, bytes, 0o644); err != nil {
		return fmt.Errorf("write model file: %w", err)
	}
	return nil
}
