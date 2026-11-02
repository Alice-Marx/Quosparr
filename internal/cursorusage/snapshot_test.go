package cursorusage

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const sample = `{
  "schemaVersion": 1,
  "fetchedAt": "2026-11-02T13:00:00+08:00",
  "plan": "Cursor Pro",
  "periodStart": "2026-10-30T17:57:00+08:00",
  "included": {
    "usedPercent": 42.5,
    "resetAt": "2026-11-06T17:57:00+08:00",
    "periodLabel": "this_week"
  },
  "onDemand": { "status": "unlimited", "usedUsd": null, "limitUsd": null },
  "breakdown": [
    { "name": "Cloud Agents", "sharePercent": 60 },
    { "name": "Tab / Inline", "sharePercent": 25 },
    { "name": "Chat", "sharePercent": 15 }
  ],
  "history": [
    { "at": "2026-10-31T13:00:00+08:00", "usedPercent": 20 },
    { "at": "2026-11-02T13:00:00+08:00", "usedPercent": 42.5 }
  ],
  "notes": "Manual or routine snapshot; not an official public API. Grok Bot assistant usage is a separate file.",
  "token": "DROP-ME-EXTRA",
  "cookie": "DROP-ME-EXTRA"
}
`

func TestParseSample(t *testing.T) {
	snap, err := Parse([]byte(sample))
	if err != nil {
		t.Fatal(err)
	}
	if snap.SchemaVersion != 1 || snap.Plan != "Cursor Pro" || snap.Included.UsedPercent != 42.5 ||
		snap.Included.PeriodLabel != "this_week" || snap.OnDemand.Status != "unlimited" || snap.OnDemand.HasUsed || snap.OnDemand.HasLimit {
		t.Fatalf("snapshot = %+v", snap)
	}
	if !snap.Included.HasReset || snap.Included.ResetAt.Format(time.RFC3339) != "2026-11-06T17:57:00+08:00" {
		t.Fatalf("reset = %s has %v", snap.Included.ResetAt, snap.Included.HasReset)
	}
	if !snap.HasPeriodStart || len(snap.Breakdown) != 3 || snap.Breakdown[0].Name != "Cloud Agents" || snap.Breakdown[0].SharePercent != 60 {
		t.Fatalf("breakdown = %+v start %v", snap.Breakdown, snap.HasPeriodStart)
	}
	if len(snap.History) != 2 || snap.History[0].UsedPercent != 20 {
		t.Fatalf("history = %+v", snap.History)
	}
	if snap.HasModelPools() || snap.GrokBotWeeklyRef {
		t.Fatalf("legacy sample invented pools: %+v", snap.IncludedPools)
	}
	body, err := Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	if strings.Contains(text, "DROP-ME-EXTRA") || strings.Contains(text, "token") || strings.Contains(text, "cookie") || strings.Contains(text, "usedUsd") {
		t.Fatalf("canonical snapshot kept an unknown or null field:\n%s", text)
	}
	if burn := BurnOf(snap); burn.Mode != BurnHistory || burn.PercentPerDay <= 0 {
		t.Fatalf("burn = %+v", burn)
	}
	if pace := PaceOf(snap); pace.Kind == PaceMissing || pace.Projected <= 0 {
		t.Fatalf("pace = %+v", pace)
	}
}

func TestExampleFilesMatchSchema(t *testing.T) {
	for _, name := range []string{"cursor-usage.json", "cursor-usage.min.json"} {
		data, err := os.ReadFile(filepath.Join("..", "..", "examples", name))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := Parse(data); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
	}
}

func TestParseRejectsUnsafeOrInvalidDocuments(t *testing.T) {
	base := `"schemaVersion":1,"fetchedAt":"2026-11-02T13:00:00Z","included":{"usedPercent":1}`
	cases := []struct {
		name string
		body string
		want error
	}{
		{"schema", `{"schemaVersion":2,"fetchedAt":"2026-11-02T13:00:00Z","included":{"usedPercent":1}}`, ErrSchema},
		{"percent", `{"schemaVersion":1,"fetchedAt":"2026-11-02T13:00:00Z","included":{"usedPercent":101}}`, ErrPercent},
		{"missing included", `{"schemaVersion":1,"fetchedAt":"2026-11-02T13:00:00Z","usedPercent":1}`, ErrField},
		{"time", `{"schemaVersion":1,"fetchedAt":"yesterday","included":{"usedPercent":1}}`, ErrTime},
		{"json", `{`, ErrJSON},
		{"control in plan", `{` + base + `,"plan":"a\nb"}`, ErrField},
		{"bad breakdown", `{` + base + `,"breakdown":[{"name":"","sharePercent":1}]}`, ErrBreakdown},
		{"negative usd", `{` + base + `,"onDemand":{"status":"limited","usedUsd":-1}}`, ErrField},
		{"bad status", `{` + base + `,"onDemand":{"status":"not a status"}}`, ErrField},
		{"pool percent", `{` + base + `,"includedPools":{"cursorModels":{"usedPercent":101}}}`, ErrPercent},
		{"pool missing percent", `{` + base + `,"includedPools":{"otherModels":{"label":"Other Models"}}}`, ErrPercent},
		{"pool label", `{` + base + `,"includedPools":{"cursorModels":{"usedPercent":1,"label":"a\nb"}}}`, ErrField},
		{"period label", `{` + base + `,"included":{"usedPercent":1,"periodLabel":"Included\nin Pro"}}`, ErrField},
		{"weekly ref type", `{` + base + `,"grokBotWeeklyRef":"yes"}`, ErrJSON},
		{"pools type", `{` + base + `,"includedPools":[]}`, ErrJSON},
		{"grok bot shape", `{"schemaVersion":1,"fetchedAt":"2026-11-02T13:00:00Z","usedPercent":6.7,"breakdown":[{"bot":"Grok Bot","sharePercent":89}]}`, ErrField},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse([]byte(tc.body))
			if !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestImportDropsUnknownFieldsAndReplacesTheSnapshot(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, FileName)
	src := filepath.Join(dir, "incoming.json")
	incoming := `{
	  "schemaVersion": 1,
	  "fetchedAt": "2026-11-02T13:00:00+08:00",
	  "included": { "usedPercent": 42.5 },
	  "onDemand": { "status": "limited", "usedUsd": 12.5, "limitUsd": 20 },
	  "apiKey": "DROP-ME-EXTRA",
	  "cookie": "DROP-ME-EXTRA"
	}`
	if err := os.WriteFile(src, []byte(incoming), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, []byte("keep-me"), 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := Import(dst, src)
	if err != nil {
		t.Fatal(err)
	}
	if file.Snapshot.Included.UsedPercent != 42.5 || !file.Snapshot.OnDemand.HasUsed || file.Snapshot.OnDemand.UsedUSD != 12.5 || file.Snapshot.OnDemand.LimitUSD != 20 || file.Path != dst {
		t.Fatalf("file = %+v", file)
	}
	if FormatUSD(file.Snapshot.OnDemand.UsedUSD) != "$12.5" || FormatUSD(20) != "$20" {
		t.Fatalf("usd format = %s %s", FormatUSD(file.Snapshot.OnDemand.UsedUSD), FormatUSD(20))
	}
	written, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(written), "DROP-ME-EXTRA") || strings.Contains(string(written), "apiKey") || strings.Contains(string(written), "cookie") {
		t.Fatal("import stored an unknown field")
	}
	if info, err := os.Stat(dst); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("perm = %v err %v", info.Mode().Perm(), err)
	}

	bad := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(bad, []byte(`{"schemaVersion":9}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Import(dst, bad); !errors.Is(err, ErrSchema) {
		t.Fatalf("bad import err = %v", err)
	}
	again, err := os.ReadFile(dst)
	if err != nil || string(again) != string(written) {
		t.Fatal("failed import replaced the snapshot")
	}
}

func TestLoadMissing(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), FileName))
	if !errors.Is(err, ErrMissing) {
		t.Fatal(err)
	}
}

func TestParseIncludedPoolsStripsUnknownFields(t *testing.T) {
	body := `{
	  "schemaVersion": 1,
	  "fetchedAt": "2026-11-02T14:00:00+08:00",
	  "plan": "Pro",
	  "included": {
	    "usedPercent": 14,
	    "resetAt": "2026-11-30T00:00:00+08:00",
	    "periodLabel": "Included in Pro"
	  },
	  "includedPools": {
	    "cursorModels": {
	      "usedPercent": 14,
	      "label": "Cursor Models",
	      "note": "Includes Cursor Grok and Composer",
	      "token": "DROP-ME-EXTRA"
	    },
	    "otherModels": { "usedPercent": 0, "label": "Other Models" }
	  },
	  "onDemand": { "status": "unlimited", "usedUsd": 0, "limitUsd": null },
	  "grokBotWeeklyRef": true,
	  "weeklyUsedPercent": 14
	}`
	snap, err := Parse([]byte(body))
	if err != nil {
		t.Fatal(err)
	}
	if snap.Included.UsedPercent != 14 || snap.Included.PeriodLabel != "Included in Pro" || !snap.GrokBotWeeklyRef || !snap.HasModelPools() {
		t.Fatalf("snapshot = %+v", snap)
	}
	if snap.IncludedPools.CursorModels.UsedPercent != 14 || snap.IncludedPools.CursorModels.Label != "Cursor Models" ||
		snap.IncludedPools.CursorModels.Note != "Includes Cursor Grok and Composer" ||
		snap.IncludedPools.OtherModels.UsedPercent != 0 || snap.IncludedPools.OtherModels.Note != "" {
		t.Fatalf("pools = %+v", snap.IncludedPools)
	}
	if !snap.OnDemand.HasUsed || snap.OnDemand.UsedUSD != 0 || snap.OnDemand.HasLimit {
		t.Fatalf("on demand = %+v", snap.OnDemand)
	}
	text := string(mustMarshal(t, snap))
	if strings.Contains(text, "DROP-ME-EXTRA") || strings.Contains(text, "token") || strings.Contains(text, "weeklyUsedPercent") || strings.Contains(text, "limitUsd") {
		t.Fatalf("canonical snapshot kept an unknown or null field:\n%s", text)
	}
	if !strings.Contains(text, `"grokBotWeeklyRef": true`) || !strings.Contains(text, `"cursorModels"`) || !strings.Contains(text, `"otherModels"`) {
		t.Fatalf("canonical snapshot dropped a pool field:\n%s", text)
	}
	again, err := Parse([]byte(text))
	if err != nil || again.IncludedPools.CursorModels.UsedPercent != 14 || !again.GrokBotWeeklyRef {
		t.Fatalf("round trip = %+v err %v", again, err)
	}
}

func TestIncludedPoolsMayBeOmitted(t *testing.T) {
	snap, err := Parse([]byte(`{"schemaVersion":1,"fetchedAt":"2026-11-02T13:00:00Z","included":{"usedPercent":14},"includedPools":{},"grokBotWeeklyRef":null}`))
	if err != nil {
		t.Fatal(err)
	}
	if snap.HasModelPools() || snap.GrokBotWeeklyRef || snap.Included.UsedPercent != 14 {
		t.Fatalf("snapshot = %+v", snap)
	}
	if strings.Contains(string(mustMarshal(t, snap)), "includedPools") || strings.Contains(string(mustMarshal(t, snap)), "grokBotWeeklyRef") {
		t.Fatal("empty pools were written back")
	}
}

func TestImportRejectsInvalidPoolsWithoutReplacingTheFile(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, FileName)
	if err := os.WriteFile(dst, []byte("keep-me"), 0o600); err != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(bad, []byte(`{"schemaVersion":1,"fetchedAt":"2026-11-02T13:00:00Z","included":{"usedPercent":1},"includedPools":{"cursorModels":{"usedPercent":140}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Import(dst, bad); !errors.Is(err, ErrPercent) {
		t.Fatalf("bad import err = %v", err)
	}
	again, err := os.ReadFile(dst)
	if err != nil || string(again) != "keep-me" {
		t.Fatal("failed pool import replaced the snapshot")
	}
}

func mustMarshal(t *testing.T, snap Snapshot) []byte {
	t.Helper()
	body, err := Marshal(snap)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func TestUsagePathUsesQuosparrDir(t *testing.T) {
	got := UsagePath(filepath.Join(string(os.PathSeparator), "home", "ada"))
	want := filepath.Join(string(os.PathSeparator), "home", "ada", ".quosparr", FileName)
	if got != want {
		t.Fatalf("path = %s, want %s", got, want)
	}
}

func TestMinSnapshotHasNoPace(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "examples", "cursor-usage.min.json"))
	if err != nil {
		t.Fatal(err)
	}
	snap, err := Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if !snap.HasModelPools() || snap.IncludedPools.CursorModels.Label != "" || snap.IncludedPools.CursorModels.Note != "" ||
		snap.IncludedPools.CursorModels.UsedPercent != 42.5 || snap.IncludedPools.OtherModels.UsedPercent != 0 || snap.GrokBotWeeklyRef {
		t.Fatalf("minimum pools = %+v weekly %v", snap.IncludedPools, snap.GrokBotWeeklyRef)
	}
	if BurnOf(snap).Mode != BurnNone {
		t.Fatal("minimum snapshot invented a burn rate")
	}
	pace := PaceOf(snap)
	if pace.Kind != PaceMissing || len(pace.Missing) != 1 || pace.Missing[0] != "included.resetAt" {
		t.Fatalf("pace = %+v", pace)
	}
}
