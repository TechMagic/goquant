package metrics

import (
	"math"
)

// Metrics calculates annualized return, volatility, Sharpe ratio, and maximum drawdown
type Metrics struct {
	AnnualReturn float64
	AnnualVol    float64
	Sharpe       float64
	MaxDrawdown  float64
}

// FromEquitySeries takes a daily/periodic equity series and periodsPerYear (e.g., 252)
func FromEquitySeries(equity []float64, periodsPerYear float64) Metrics {
	n := len(equity)
	if n < 2 {
		return Metrics{}
	}

	// returns: daily/simple returns
	rets := make([]float64, 0, n-1)
	for i := 1; i < n; i++ {
		r := equity[i]/equity[i-1] - 1.0
		rets = append(rets, r)
	}

	// mean and std
	var sum float64
	for _, r := range rets {
		sum += r
	}
	mean := sum / float64(len(rets))
	var sdsum float64
	for _, r := range rets {
		sdsum += (r - mean) * (r - mean)
	}
	std := math.Sqrt(sdsum / float64(len(rets)))

	annRet := math.Pow(1+mean, periodsPerYear) - 1
	annVol := std * math.Sqrt(periodsPerYear)
	sharpe := 0.0
	if annVol > 0 {
		sharpe = annRet / annVol
	}

	// max drawdown
	peak := equity[0]
	maxdd := 0.0
	for _, v := range equity {
		if v > peak {
			peak = v
		}
		dd := (peak - v) / peak
		if dd > maxdd {
			maxdd = dd
		}
	}

	return Metrics{AnnualReturn: annRet, AnnualVol: annVol, Sharpe: sharpe, MaxDrawdown: maxdd}
}
