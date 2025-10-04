package backtest

import (
	"errors"
	"goquant/internal/data"
	"goquant/internal/strategy"
	"time"
)

// Trade records a single trade
type Trade struct {
	Time      time.Time // wrapper below: but to avoid cycles, we'll use string in printing
	Side      string    // "BUY" / "SELL"
	Price     float64
	Quantity  float64
	CashAfter float64
	Position  float64
}

// BacktestReports returns the backtest results
type BacktestReports struct {
	Trades      []Trade
	EquityCurve []float64
	FinalEquity float64
}

// Engine is a simple daily backtesting engine (for demonstration only, without fees/slippage)
type Engine struct {
	Initial float64
}

func NewEngine(initial float64) *Engine { return &Engine{Initial: initial} }

func (e *Engine) Run(bars []data.Bar, strat strategy.Strategy) (*BacktestReports, error) {
	if len(bars) == 0 {
		return nil, errors.New("no bars")
	}

	cash := e.Initial
	position := 0.0 // Position size (units: contract quantity or asset shares)
	var trades []Trade
	equity := make([]float64, 0, len(bars))

	for _, b := range bars {
		sig := strat.OnBar(b)
		price := b.Close
		// Simple logic: two states - holding a position or not holding a position
		if sig == 1 && position == 0.0 {
			// buy with all cash
			qty := cash / price
			position = qty
			cash = 0.0
			trades = append(trades, makeTrade(b, "BUY", price, qty, cash, position))
		} else if sig == -1 && position > 0.0 {
			// sell all
			cash = position * price
			qty := position
			position = 0.0
			trades = append(trades, makeTrade(b, "SELL", price, qty, cash, position))
		}
		// equity = cash + position * price
		eq := cash + position*price
		equity = append(equity, eq)
	}

	final := equity[len(equity)-1]
	return &BacktestReports{Trades: trades, EquityCurve: equity, FinalEquity: final}, nil
}

func makeTrade(b data.Bar, side string, price, qty, cash, pos float64) Trade {
	return Trade{
		Time:      b.Date,
		Side:      side,
		Price:     price,
		Quantity:  qty,
		CashAfter: cash,
		Position:  pos,
	}
}
