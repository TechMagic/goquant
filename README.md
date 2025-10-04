
# GoQuant — A Simple & Fast Backtesting Framework in Go

GoQuant is a lightweight, modular backtesting engine for quantitative trading strategies written in **Go**. It supports CSV-based historical data, built-in performance metrics (Sharpe, Max Drawdown, etc.), and is designed for easy extension with new strategies.

Perfect for learning quant development, testing ideas, or building a foundation for a live trading system.

---

## 🚀 Quick Start

### Prerequisites

- [Go](https://golang.org/dl/) 1.16+
- Historical market data in CSV format (see [Data Format](#-data-format) below)

### 1. Clone the repository

```
git clone https://github.com/techmagic/goquant.git
cd goquant
```

### 2. Prepare your data

Place your OHLCV CSV files in the data/ directory:

``` 
mkdir -p data

# Example: download BTC-USD.csv from Yahoo Finance and put it here
cp /path/to/BTC-USD.csv data/
```

### 3. Run a backtest

```
go run cmd/goquantlab/main.go --symbol BTC-USD --ma-fast 10 --ma-slow 30
```

You’ll see a summary like:

```
=== Summary ===
Symbol: BTC-USD
Bars: 4034
Initial capital: 10000.00
Final equity: 2544418.79
Total return: 25344.19%
Max Drawdown: 63.25%
Annualized Return: 52.93%
Annualized Volatility: 39.84%
Sharpe: 1.329
Trades:
2014-10-18 BUY qty=25.5466 price=391.44 cash=0.00
2014-10-28 SELL qty=25.5466 price=357.62 cash=9135.91
...
```

## 📊 Data Format

Your CSV files must follow this format:

- File location: data/<symbol>.csv (e.g., data/BTC-USD.csv)
- Header row: Required (will be skipped)
- Columns (in any order):
  - Date (YYYY-MM-DD)
  - Open
  - High
  - Low
  - Close
  - Volume

✅ Example (data/BTC-USD.csv):

```
Date,Open,High,Low,Close,Volume
2014-09-17,465.86,468.17,452.42,457.33,21056800
2014-09-18,456.86,456.86,413.10,424.44,34483200
...
```

🔍 Note: The parser reads columns by name, not position — so column order doesn’t matter! 
💡 Where to get data?
Download free historical data from Yahoo Finance : 

1. Search for a symbol (e.g., BTC-USD)
2. Go to Historical Data → Download

## ⚙️ Command-Line Options

| Flag      | Default              | Description                                 |
| --------- | -------------------- | ------------------------------------------- |
| --config  | configs/example.yaml | Path to config file (currently minimal use) |
| --symbol  | BTC-USD              | Symbol name (CSV filename without .csv)     |
| --start   | (all data)           | Start date (YYYY-MM-DD)                     |
| --end     | (all data)           | End date (YYYY-MM-DD)                       |
| --initial | 10000                | Initial capital                             |
| --ma-fast | 10                   | Fast moving average period                  |
| --ma-slow | 30                   | Slow moving average period                  |

## 🧩 Architecture Overview

```tree
goquant/
├── cmd/goquantlab/      # CLI entry point
├── internal/
│   ├── data/            # CSV data loader
│   ├── strategy/        # Trading strategies (e.g., MA Crossover)
│   ├── backtest/        # Backtesting engine
│   ├── metrics/         # Performance metrics (Sharpe, Drawdown, etc.)
│   └── utils/           # Config loader
├── data/                # Place your CSV files here
├── configs/             # Example config
└── go.mod               # Go module
```

- Strategy Interface: Easy to add new strategies (implement OnBar)
- Modular Design: Swap data sources, metrics, or position sizing without touching core logic

## 🛠️ Extending GoQuant

**Add a new strategy**

- Create internal/strategy/mystategy.go
- Implement the `Strategy` interface:

``` type MyStrategy struct { ... }
type MyStrategy struct { ... }
func (s *MyStrategy) OnBar(bar data.Bar) int { ... }
```

- Register it in main.go

**Add transaction costs**
Modify `internal/backtest/backtest.go` to deduct fees on trades.

## 📜 License

MIT License — feel free to use, modify, and distribute.

## 🙌 Acknowledgements

Inspired by classic quant backtesting frameworks
Uses standard Go tooling — no external dependencies

Happy backtesting! 📈
If you find this useful, consider giving it a ⭐ on GitHub! 
