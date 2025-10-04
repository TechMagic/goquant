package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"goquant/internal/backtest"
	"goquant/internal/data"
	"goquant/internal/metrics"
	"goquant/internal/strategy"
	"goquant/internal/utils"
)

func main() {
	cfgPath := flag.String("config", "configs/example.yaml", "path to config file")
	symbol := flag.String("symbol", "BTC-USD", "symbol to backtest (csv filename under data/ without extension)")
	start := flag.String("start", "", "start date YYYY-MM-DD")
	end := flag.String("end", "", "end date YYYY-MM-DD")
	initial := flag.Float64("initial", 10000.0, "initial capital")
	maFast := flag.Int("ma-fast", 10, "fast MA period")
	maSlow := flag.Int("ma-slow", 30, "slow MA period")
	flag.Parse()

	cfg, err := utils.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	_ = cfg // currently not used extensively but kept for extension

	provider := data.NewCSVDataProvider(cfg.DataDir) // ✅ 使用 cfg.DataDir

	var st strategy.Strategy
	st = strategy.NewMACrossover(*maFast, *maSlow)

	var s time.Time
	var e time.Time
	if *start != "" {
		s, _ = time.Parse("2006-01-02", *start)
	}
	if *end != "" {
		e, _ = time.Parse("2006-01-02", *end)
	}

	bars, err := provider.LoadBars(*symbol, s, e) // ✅ 传入 time.Time

	if err != nil {
		log.Fatalf("load bars: %v", err)
	}

	engine := backtest.NewEngine(*initial)
	reports, err := engine.Run(bars, st)
	if err != nil {
		log.Fatalf("backtest run: %v", err)
	}

	fmt.Println("=== Summary ===")
	fmt.Printf("Symbol: %s\n", *symbol)
	fmt.Printf("Bars: %d\n", len(bars))
	fmt.Printf("Initial capital: %.2f\n", *initial)
	fmt.Printf("Final equity: %.2f\n", reports.FinalEquity)
	fmt.Printf("Total return: %.2f%%\n", (reports.FinalEquity-*initial)/(*initial)*100)

	m := metrics.FromEquitySeries(reports.EquityCurve, 252.0)
	fmt.Printf("Max Drawdown: %.2f%%\n", m.MaxDrawdown*100)
	fmt.Printf("Annualized Return: %.2f%%\n", m.AnnualReturn*100)
	fmt.Printf("Annualized Volatility: %.2f%%\n", m.AnnualVol*100)
	fmt.Printf("Sharpe: %.3f\n", m.Sharpe)

	fmt.Println("Trades:")
	for _, t := range reports.Trades {
		fmt.Printf("%s %s qty=%.4f price=%.2f cash=%.2f\n", t.Time.Format("2006-01-02"), t.Side, t.Quantity, t.Price, t.CashAfter)
	}
}
