package grokbot

import (
	"encoding/json"
	"math"
	"regexp"
	"time"
)

// HistoryPoint is one optional used-percent sample.
type HistoryPoint struct {
	At          time.Time
	UsedPercent float64
}

// CacheStats is an optional cache-hit block. Present is false for a missing
// or empty object; the assistant quota itself has no cache metric.
type CacheStats struct {
	Present    bool
	HasHit     bool
	HitPercent float64
	HasHits    bool
	Hits       int64
	HasMisses  bool
	Misses     int64
	Notes      string
}

// UsageStats is an optional token and cost block. InputTokens includes
// cache reads, matching the other Quosparr price tables.
type UsageStats struct {
	HasAny          bool
	HasInput        bool
	InputTokens     int64
	HasOutput       bool
	OutputTokens    int64
	HasCache        bool
	CacheReadTokens int64
	Model           string
	HasCost         bool
	CostUSD         float64
}

// CostLine is one optional model row.
type CostLine struct {
	Model           string
	Label           string
	InputTokens     int64
	OutputTokens    int64
	CacheReadTokens int64
	HasCost         bool
	CostUSD         float64
}

type snapshotExtras struct {
	PeriodStart    time.Time
	HasPeriodStart bool
	History        []HistoryPoint
	Cache          CacheStats
	Usage          UsageStats
	CostLines      []CostLine
}

type rawHistory struct {
	At          string       `json:"at"`
	UsedPercent *json.Number `json:"usedPercent"`
}

type rawCache struct {
	HitPercent *json.Number `json:"hitPercent"`
	Hits       *json.Number `json:"hits"`
	Misses     *json.Number `json:"misses"`
	Notes      string       `json:"notes"`
}

type rawUsage struct {
	InputTokens     *json.Number `json:"inputTokens"`
	OutputTokens    *json.Number `json:"outputTokens"`
	CacheReadTokens *json.Number `json:"cacheReadTokens"`
	Model           string       `json:"model"`
	CostUSD         *json.Number `json:"costUsd"`
}

type rawCostLine struct {
	Model           string       `json:"model"`
	Label           string       `json:"label"`
	InputTokens     *json.Number `json:"inputTokens"`
	OutputTokens    *json.Number `json:"outputTokens"`
	CacheReadTokens *json.Number `json:"cacheReadTokens"`
	CostUSD         *json.Number `json:"costUsd"`
}

type encodedHistory struct {
	At          string      `json:"at"`
	UsedPercent json.Number `json:"usedPercent"`
}

type encodedCache struct {
	HitPercent *json.Number `json:"hitPercent,omitempty"`
	Hits       *json.Number `json:"hits,omitempty"`
	Misses     *json.Number `json:"misses,omitempty"`
	Notes      string       `json:"notes,omitempty"`
}

type encodedUsage struct {
	InputTokens     *json.Number `json:"inputTokens,omitempty"`
	OutputTokens    *json.Number `json:"outputTokens,omitempty"`
	CacheReadTokens *json.Number `json:"cacheReadTokens,omitempty"`
	Model           string       `json:"model,omitempty"`
	CostUSD         *json.Number `json:"costUsd,omitempty"`
}

type encodedCostLine struct {
	Model           string       `json:"model,omitempty"`
	Label           string       `json:"label,omitempty"`
	InputTokens     json.Number  `json:"inputTokens"`
	OutputTokens    json.Number  `json:"outputTokens"`
	CacheReadTokens json.Number  `json:"cacheReadTokens"`
	CostUSD         *json.Number `json:"costUsd,omitempty"`
}

var modelPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/-]{0,79}$`)

func parseExtras(r rawSnapshot) (snapshotExtras, error) {
	var out snapshotExtras
	start, hasStart, err := parseOptionalTime(r.PeriodStart)
	if err != nil {
		return out, err
	}
	out.PeriodStart, out.HasPeriodStart = start, hasStart
	history, err := parseHistory(r.History)
	if err != nil {
		return out, err
	}
	cache, err := parseCache(r.Cache)
	if err != nil {
		return out, err
	}
	usage, err := parseUsage(r.Usage)
	if err != nil {
		return out, err
	}
	lines, err := parseCostLines(r.CostLines)
	if err != nil {
		return out, err
	}
	out.History, out.Cache, out.Usage, out.CostLines = history, cache, usage, lines
	return out, nil
}

func parseHistory(raw []rawHistory) ([]HistoryPoint, error) {
	if len(raw) > maxHistory {
		return nil, ErrField
	}
	out := make([]HistoryPoint, 0, len(raw))
	for _, item := range raw {
		at, err := parseTime(item.At, true)
		if err != nil {
			return nil, err
		}
		used, err := parsePercent(item.UsedPercent, true)
		if err != nil {
			return nil, ErrPercent
		}
		out = append(out, HistoryPoint{At: at, UsedPercent: used})
	}
	return out, nil
}

func parseCache(raw *rawCache) (CacheStats, error) {
	if raw == nil {
		return CacheStats{}, nil
	}
	var out CacheStats
	if raw.HitPercent != nil {
		v, err := parsePercent(raw.HitPercent, true)
		if err != nil {
			return CacheStats{}, ErrPercent
		}
		out.HasHit, out.HitPercent, out.Present = true, v, true
	}
	if raw.Hits != nil {
		n, err := parseCount(raw.Hits)
		if err != nil {
			return CacheStats{}, err
		}
		out.HasHits, out.Hits, out.Present = true, n, true
	}
	if raw.Misses != nil {
		n, err := parseCount(raw.Misses)
		if err != nil {
			return CacheStats{}, err
		}
		out.HasMisses, out.Misses, out.Present = true, n, true
	}
	notes, err := optionalLabel(raw.Notes, 400)
	if err != nil {
		return CacheStats{}, err
	}
	if notes != "" {
		out.Notes, out.Present = notes, true
	}
	return out, nil
}

func parseUsage(raw *rawUsage) (UsageStats, error) {
	if raw == nil {
		return UsageStats{}, nil
	}
	var out UsageStats
	if raw.InputTokens != nil {
		n, err := parseCount(raw.InputTokens)
		if err != nil {
			return UsageStats{}, err
		}
		out.HasInput, out.InputTokens, out.HasAny = true, n, true
	}
	if raw.OutputTokens != nil {
		n, err := parseCount(raw.OutputTokens)
		if err != nil {
			return UsageStats{}, err
		}
		out.HasOutput, out.OutputTokens, out.HasAny = true, n, true
	}
	if raw.CacheReadTokens != nil {
		n, err := parseCount(raw.CacheReadTokens)
		if err != nil {
			return UsageStats{}, err
		}
		out.HasCache, out.CacheReadTokens, out.HasAny = true, n, true
	}
	model, err := optionalModel(raw.Model)
	if err != nil {
		return UsageStats{}, err
	}
	if model != "" {
		out.Model, out.HasAny = model, true
	}
	if raw.CostUSD != nil {
		v, err := parseMoney(raw.CostUSD)
		if err != nil {
			return UsageStats{}, err
		}
		out.HasCost, out.CostUSD, out.HasAny = true, v, true
	}
	return out, nil
}

func parseCostLines(raw []rawCostLine) ([]CostLine, error) {
	if len(raw) > maxCostLines {
		return nil, ErrField
	}
	out := make([]CostLine, 0, len(raw))
	for _, item := range raw {
		model, err := optionalModel(item.Model)
		if err != nil {
			return nil, err
		}
		label, err := optionalLabel(item.Label, 80)
		if err != nil {
			return nil, err
		}
		in, err := countOrZero(item.InputTokens)
		if err != nil {
			return nil, err
		}
		outN, err := countOrZero(item.OutputTokens)
		if err != nil {
			return nil, err
		}
		cache, err := countOrZero(item.CacheReadTokens)
		if err != nil {
			return nil, err
		}
		line := CostLine{Model: model, Label: label, InputTokens: in, OutputTokens: outN, CacheReadTokens: cache}
		if item.CostUSD != nil {
			v, err := parseMoney(item.CostUSD)
			if err != nil {
				return nil, err
			}
			line.HasCost, line.CostUSD = true, v
		}
		out = append(out, line)
	}
	return out, nil
}

func encodeExtras(out *encodedSnapshot, s Snapshot) {
	if len(s.History) > 0 {
		out.History = make([]encodedHistory, len(s.History))
		for i, p := range s.History {
			out.History[i] = encodedHistory{At: p.At.Format(time.RFC3339), UsedPercent: json.Number(FormatPercent(p.UsedPercent))}
		}
	}
	if s.Cache.Present {
		cache := &encodedCache{Notes: s.Cache.Notes}
		if s.Cache.HasHit {
			n := json.Number(FormatPercent(s.Cache.HitPercent))
			cache.HitPercent = &n
		}
		if s.Cache.HasHits {
			n := json.Number(formatCount(s.Cache.Hits))
			cache.Hits = &n
		}
		if s.Cache.HasMisses {
			n := json.Number(formatCount(s.Cache.Misses))
			cache.Misses = &n
		}
		out.Cache = cache
	}
	if s.Usage.HasAny {
		usage := &encodedUsage{Model: s.Usage.Model}
		if s.Usage.HasInput {
			n := json.Number(formatCount(s.Usage.InputTokens))
			usage.InputTokens = &n
		}
		if s.Usage.HasOutput {
			n := json.Number(formatCount(s.Usage.OutputTokens))
			usage.OutputTokens = &n
		}
		if s.Usage.HasCache {
			n := json.Number(formatCount(s.Usage.CacheReadTokens))
			usage.CacheReadTokens = &n
		}
		if s.Usage.HasCost {
			n := json.Number(formatMoney(s.Usage.CostUSD))
			usage.CostUSD = &n
		}
		out.Usage = usage
	}
	if len(s.CostLines) > 0 {
		out.CostLines = make([]encodedCostLine, len(s.CostLines))
		for i, line := range s.CostLines {
			item := encodedCostLine{
				Model:           line.Model,
				Label:           line.Label,
				InputTokens:     json.Number(formatCount(line.InputTokens)),
				OutputTokens:    json.Number(formatCount(line.OutputTokens)),
				CacheReadTokens: json.Number(formatCount(line.CacheReadTokens)),
			}
			if line.HasCost {
				n := json.Number(formatMoney(line.CostUSD))
				item.CostUSD = &n
			}
			out.CostLines[i] = item
		}
	}
}

func parseCount(n *json.Number) (int64, error) {
	if n == nil {
		return 0, ErrField
	}
	if v, err := n.Int64(); err == nil {
		if v < 0 {
			return 0, ErrField
		}
		return v, nil
	}
	f, err := n.Float64()
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) || f < 0 || f > math.MaxInt64 || f != math.Trunc(f) {
		return 0, ErrField
	}
	return int64(f), nil
}

func countOrZero(n *json.Number) (int64, error) {
	if n == nil {
		return 0, nil
	}
	return parseCount(n)
}

func parseMoney(n *json.Number) (float64, error) {
	v, err := n.Float64()
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
		return 0, ErrField
	}
	return v, nil
}

func optionalModel(raw string) (string, error) {
	raw, err := optionalLabel(raw, 80)
	if err != nil || raw == "" {
		return raw, err
	}
	if !modelPattern.MatchString(raw) {
		return "", ErrField
	}
	return raw, nil
}

func formatCount(n int64) string {
	return json.Number(itoa(n)).String()
}

func itoa(n int64) string {
	return jsonNumber(n)
}

func jsonNumber(n int64) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func formatMoney(v float64) string {
	rounded := math.Round(v*1e6) / 1e6
	s := trimFloat(rounded)
	return s
}

func trimFloat(v float64) string {
	s := json.Number(strconvFormat(v)).String()
	return s
}

func strconvFormat(v float64) string {
	// Up to 6 decimal places, trailing zeros removed by the caller via FormatPercent-like logic.
	whole := math.Trunc(v)
	frac := math.Round((v - whole) * 1e6)
	if frac >= 1e6 {
		whole++
		frac = 0
	}
	out := jsonNumber(int64(whole))
	if frac == 0 {
		return out
	}
	digits := jsonNumber(int64(frac))
	for len(digits) < 6 {
		digits = "0" + digits
	}
	for len(digits) > 0 && digits[len(digits)-1] == '0' {
		digits = digits[:len(digits)-1]
	}
	return out + "." + digits
}
