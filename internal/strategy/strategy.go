package strategy

import (
	"goquant/internal/data"
	"time"
)

// Strategy 接口：OnBar 返回信号：-1 sell, 0 hold, +1 buy
type Strategy interface {
	OnBar(bar data.Bar) int
	Name() string
}

// TradeSignal 可选封装
type TradeSignal struct {
	Time time.Time
	Side int // -1 sell, +1 buy
}
