package strategy

import (
	"goquant/internal/data"
	"time"
)

// Strategy interface: OnBar returns signal: -1 sell, 0 hold, +1 buy
type Strategy interface {
	OnBar(bar data.Bar) int
	Name() string
}

// TradeSignal optional wrapper
type TradeSignal struct {
	Time time.Time
	Side int // -1 sell, +1 buy
}
