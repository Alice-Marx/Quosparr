package grokbot

import (
	"strings"
	"testing"
	"time"
)

func TestParseKeepsOptionalFieldsAndDropsSecrets(t *testing.T) {
	raw := `{
	  "schemaVersion": 1,
	  "fetchedAt": "2026-11-02T10:00:00+08:00",
	  "usedPercent": 6.7,
	  "resetAt": "2026-11-06T17:57:00+08:00",
	  "periodStart": "2026-10-30T17:57:00+08:00",
	  "periodLabel": "this_week",
	  "apiKey": "DROP-ME-EXTRA",
	  "history": [
	    {"at": "2026-10-31T10:00:00+08:00", "usedPercent": 2.1},
	    {"at": "2026-11-02T10:00:00+08:00", "usedPercent": 6.7}
	  ],
	  "cache": {"hitPercent": 12.5, "hits": 10, "misses": 70, "notes": "optional"},
	  "usage": {"inputTokens": 1200, "outputTokens": 300, "cacheReadTokens": 400, "model": "gpt-4o", "costUsd": 0.0125},
	  "costLines": [{"model": "gpt-4o", "label": "demo", "inputTokens": 1200, "outputTokens": 300, "cacheReadTokens": 400, "costUsd": 0.0125}]
	}`
	snap, err := Parse([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if !snap.HasPeriodStart || len(snap.History) != 2 || !snap.Cache.Present || snap.Cache.HitPercent != 12.5 || snap.Cache.Hits != 10 {
		t.Fatalf("optional fields %+v", snap)
	}
	if !snap.Usage.HasAny || snap.Usage.InputTokens != 1200 || !snap.Usage.HasCost || snap.Usage.CostUSD != 0.0125 || len(snap.CostLines) != 1 {
		t.Fatalf("usage %+v lines %+v", snap.Usage, snap.CostLines)
	}
	body, err := Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "DROP-ME-EXTRA") || strings.Contains(string(body), "apiKey") {
		t.Fatalf("secret written back: %s", body)
	}
	again, err := Parse(body)
	if err != nil || again.Usage.Model != "gpt-4o" || again.History[0].UsedPercent != 2.1 {
		t.Fatalf("round trip %+v %v", again, err)
	}
	emptyCache, err := Parse([]byte(`{"schemaVersion":1,"fetchedAt":"2026-11-02T10:00:00Z","usedPercent":1,"cache":{}}`))
	if err != nil || emptyCache.Cache.Present {
		t.Fatalf("empty cache %+v %v", emptyCache.Cache, err)
	}
	if _, err := Parse([]byte(`{"schemaVersion":1,"fetchedAt":"2026-11-02T10:00:00Z","usedPercent":1,"history":[{"at":"nope","usedPercent":1}]}`)); err != ErrTime {
		t.Fatalf("bad history %v", err)
	}
}

func TestBurnAndPace(t *testing.T) {
	fetched := time.Date(2026, 11, 2, 2, 0, 0, 0, time.UTC)
	reset := time.Date(2026, 11, 6, 9, 57, 0, 0, time.UTC)
	start := time.Date(2026, 10, 30, 9, 57, 0, 0, time.UTC)
	snap := Snapshot{
		FetchedAt: fetched, UsedPercent: 6.7, HasReset: true, ResetAt: reset,
		HasPeriodStart: true, PeriodStart: start, PeriodLabel: "this_week",
		History: []HistoryPoint{
			{At: fetched.Add(-48 * time.Hour), UsedPercent: 2.1},
			{At: fetched, UsedPercent: 6.7},
		},
	}
	burn := BurnOf(snap)
	if burn.Mode != BurnHistory || burn.PercentPerDay <= 0 || !burn.HasETA {
		t.Fatalf("burn %+v", burn)
	}
	pace := PaceOf(snap)
	if pace.Kind != PaceUnder || pace.Projected <= 0 || pace.AssumedWeek {
		t.Fatalf("pace %+v", pace)
	}
	bare := Snapshot{FetchedAt: fetched, UsedPercent: 6.7}
	if BurnOf(bare).Mode != BurnNone {
		t.Fatal("bare snapshot invented a rate")
	}
	if missing := PaceOf(bare).Missing; len(missing) != 1 || missing[0] != "resetAt" {
		t.Fatalf("missing %+v", PaceOf(bare))
	}
	week := Snapshot{FetchedAt: fetched, UsedPercent: 50, HasReset: true, ResetAt: fetched.Add(24 * time.Hour), PeriodLabel: "this_week"}
	got := PaceOf(week)
	if !got.AssumedWeek || got.Kind == PaceMissing {
		t.Fatalf("assumed week %+v", got)
	}
}

func TestCostRowsPreferExplicitAndMarkMissingPrice(t *testing.T) {
	rows, ok := CostRows(Snapshot{Usage: UsageStats{HasAny: true, HasInput: true, InputTokens: 10, Model: "not-a-real-model"}})
	if !ok || len(rows) != 1 || !rows[0].Missing || rows[0].Amount != "" {
		t.Fatalf("missing price %+v ok=%v", rows, ok)
	}
	rows, ok = CostRows(Snapshot{Usage: UsageStats{HasAny: true, HasCost: true, CostUSD: 1.5, Model: "not-a-real-model"}})
	if !ok || rows[0].Missing || rows[0].Amount != "$1.5" {
		t.Fatalf("explicit %+v", rows)
	}
	if _, ok := CostRows(Snapshot{UsedPercent: 6.7}); ok {
		t.Fatal("used percent is not a token bill")
	}
}
