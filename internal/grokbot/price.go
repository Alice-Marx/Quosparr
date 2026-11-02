package grokbot

import (
	"math"

	"github.com/licoy/quosparr/internal/claude"
	"github.com/licoy/quosparr/internal/grok"
	"github.com/licoy/quosparr/internal/usage"
	"github.com/licoy/quosparr/internal/usage/decimal"
)

// CostRow is one priced line for the assistant panel. Amount is empty when
// the model has no price and the snapshot did not give costUsd.
type CostRow struct {
	Model   string
	Label   string
	Input   int64
	Output  int64
	Cache   int64
	Amount  string
	Missing bool
}

// CostRows prices usage and costLines. An explicit costUsd wins over the
// built-in tables. Input includes cache reads.
func CostRows(s Snapshot) (rows []CostRow, ok bool) {
	if s.Usage.HasAny {
		ok = true
		rows = append(rows, priceTokens(s.Usage.Model, "", s.Usage.InputTokens, s.Usage.OutputTokens, s.Usage.CacheReadTokens, s.Usage.HasCost, s.Usage.CostUSD))
	}
	for _, line := range s.CostLines {
		ok = true
		rows = append(rows, priceTokens(line.Model, line.Label, line.InputTokens, line.OutputTokens, line.CacheReadTokens, line.HasCost, line.CostUSD))
	}
	return rows, ok
}

func priceTokens(model, label string, input, output, cache int64, hasCost bool, cost float64) CostRow {
	row := CostRow{Model: model, Label: label, Input: input, Output: output, Cache: cache}
	if hasCost {
		row.Amount = formatUSD(cost)
		return row
	}
	if amount, ok := lookupPrice(model, input, output, cache); ok {
		row.Amount = amount
		return row
	}
	row.Missing = true
	return row
}

func lookupPrice(model string, input, output, cache int64) (string, bool) {
	if model == "" {
		return "", false
	}
	cached := cache
	if cached > input {
		cached = input
	}
	if amount, ok := usage.CostForUsage(model, usage.TokenUsage{InputTokens: input, CachedInputTokens: cached, OutputTokens: output}); ok {
		return formatDecimal(amount), true
	}
	if input <= math.MaxInt && cached <= math.MaxInt && output <= math.MaxInt {
		if amount, ok := grok.PricingCost(model, int(input), int(cached), int(output)); ok {
			return formatDecimal(amount), true
		}
	}
	uncached := input - cached
	if uncached < 0 {
		uncached = 0
	}
	if amount, ok := claude.PricingCost(model, claude.TokenParts{Input: uncached, CacheRead: cached, Output: output}, ""); ok {
		return formatDecimal(amount), true
	}
	return "", false
}

func formatUSD(v float64) string {
	return "$" + trimFloat(v)
}

func formatDecimal(d decimal.Decimal) string {
	return "$" + d.String()
}
