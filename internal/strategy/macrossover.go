package strategy

import (
	"container/list"
	"fmt"

	"goquant/internal/data"
)

// 简单双均线交叉策略
type MACrossover struct {
	fast int
	slow int
	// rolling windows
	fastWindow *list.List
	slowWindow *list.List
	sumFast    float64
	sumSlow    float64
	lastSignal int
}

func NewMACrossover(fast, slow int) *MACrossover {
	if fast >= slow {
		panic("fast must be < slow")
	}
	return &MACrossover{
		fast:       fast,
		slow:       slow,
		fastWindow: list.New(),
		slowWindow: list.New(),
	}
}

func (m *MACrossover) Name() string { return fmt.Sprintf("MA%d/%d", m.fast, m.slow) }

func (m *MACrossover) OnBar(bar data.Bar) int {
	// update windows
	m.fastWindow.PushBack(bar.Close)
	m.sumFast += bar.Close
	if m.fastWindow.Len() > m.fast {
		front := m.fastWindow.Front()
		m.sumFast -= front.Value.(float64)
		m.fastWindow.Remove(front)
	}

	m.slowWindow.PushBack(bar.Close)
	m.sumSlow += bar.Close
	if m.slowWindow.Len() > m.slow {
		front := m.slowWindow.Front()
		m.sumSlow -= front.Value.(float64)
		m.slowWindow.Remove(front)
	}

	if m.fastWindow.Len() < m.fast || m.slowWindow.Len() < m.slow {
		return 0 // not enough data
	}

	maFast := m.sumFast / float64(m.fastWindow.Len())
	maSlow := m.sumSlow / float64(m.slowWindow.Len())

	var sig int
	if maFast > maSlow {
		sig = 1
	} else if maFast < maSlow {
		sig = -1
	} else {
		sig = 0
	}

	// only trigger on cross (change in signal)
	if sig != m.lastSignal {
		m.lastSignal = sig
		return sig
	}
	return 0
}
