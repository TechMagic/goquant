// internal/data/data.go
package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Bar struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

type CSVDataProvider struct {
	basePath string
}

// 注意：方法名改为 NewCSVDataProvider（更清晰），结构体名也同步
func NewCSVDataProvider(basePath string) *CSVDataProvider {
	return &CSVDataProvider{basePath: basePath}
}

// LoadBars 加载指定 symbol 的数据，并按 start/end 时间范围过滤
// 如果 start 或 end 为零值（time.Time{}），则不过滤对应边界
func (p *CSVDataProvider) LoadBars(symbol string, start, end time.Time) ([]Bar, error) {
	path := fmt.Sprintf("%s/%s.csv", p.basePath, symbol)
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file %s: %w", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("csv has no data rows")
	}

	// === 新增：解析表头 ===
	header := rows[0]
	colIndex := make(map[string]int)
	for i, name := range header {
		colIndex[name] = i
	}

	// 检查必要列是否存在
	requiredCols := []string{"Date", "Open", "High", "Low", "Close", "Volume"}
	for _, col := range requiredCols {
		if _, exists := colIndex[col]; !exists {
			return nil, fmt.Errorf("missing required column: %s", col)
		}
	}

	var allBars []Bar
	for i := 1; i < len(rows); i++ { // 从第2行开始（跳过表头）
		row := rows[i]
		if len(row) == 0 {
			continue
		}

		// === 按列名读取，不再依赖固定位置 ===
		dateStr := row[colIndex["Date"]]
		openStr := row[colIndex["Open"]]
		highStr := row[colIndex["High"]]
		lowStr := row[colIndex["Low"]]
		closeStr := row[colIndex["Close"]]
		volumeStr := row[colIndex["Volume"]]

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("parse date '%s' at row %d: %w", dateStr, i, err)
		}

		open, err := strconv.ParseFloat(openStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parse open '%s' at row %d: %w", openStr, i, err)
		}

		high, err := strconv.ParseFloat(highStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parse high '%s' at row %d: %w", highStr, i, err)
		}

		low, err := strconv.ParseFloat(lowStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parse low '%s' at row %d: %w", lowStr, i, err)
		}

		close, err := strconv.ParseFloat(closeStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parse close '%s' at row %d: %w", closeStr, i, err)
		}

		volume, err := strconv.ParseFloat(volumeStr, 64)
		if err != nil {
			return nil, fmt.Errorf("parse volume '%s' at row %d: %w", volumeStr, i, err)
		}

		allBars = append(allBars, Bar{
			Date:   date,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}

	// === 时间过滤（保持不变）===
	var filtered []Bar
	for _, bar := range allBars {
		if !start.IsZero() && bar.Date.Before(start) {
			continue
		}
		if !end.IsZero() && bar.Date.After(end) {
			continue
		}
		filtered = append(filtered, bar)
	}

	return filtered, nil
}
