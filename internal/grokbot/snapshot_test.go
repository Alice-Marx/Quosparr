package grokbot

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
  "fetchedAt": "2026-11-02T10:00:00+08:00",
  "plan": "Cursor Pro",
  "usedPercent": 6.7,
  "resetAt": "2026-11-06T17:57:00+08:00",
  "onDemand": { "status": "unlimited" },
  "periodLabel": "this_week",
  "breakdown": [
    { "bot": "Grok Bot", "sharePercent": 89, "activities": [ { "name": "chats", "sharePercent": 87 } ] }
  ],
  "notes": "Written by Grok Bot routine or manual import; not an official public API."
}
`

func TestParseSample(t *testing.T) {
	snap, err := Parse([]byte(sample))
	if err != nil {
		t.Fatal(err)
	}
	if snap.SchemaVersion != 1 || snap.Plan != "Cursor Pro" || snap.UsedPercent != 6.7 || snap.OnDemand != "unlimited" || snap.PeriodLabel != "this_week" {
		t.Fatalf("snapshot = %+v", snap)
	}
	if !snap.HasReset || snap.ResetAt.Format(time.RFC3339) != "2026-11-06T17:57:00+08:00" {
		t.Fatalf("reset = %s has %v", snap.ResetAt, snap.HasReset)
	}
	if len(snap.Breakdown) != 1 || snap.Breakdown[0].Bot != "Grok Bot" || snap.Breakdown[0].SharePercent != 89 ||
		len(snap.Breakdown[0].Activities) != 1 || snap.Breakdown[0].Activities[0].Name != "chats" || snap.Breakdown[0].Activities[0].SharePercent != 87 {
		t.Fatalf("breakdown = %+v", snap.Breakdown)
	}
	if FormatPercent(snap.UsedPercent) != "6.7" || FormatPercent(89) != "89" {
		t.Fatalf("format = %s %s", FormatPercent(snap.UsedPercent), FormatPercent(89))
	}
}

func TestExampleFileMatchesSchema(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "examples", "grokbot-usage.json"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Parse(data); err != nil {
		t.Fatal(err)
	}
}

func TestParseRejectsUnsafeOrInvalidDocuments(t *testing.T) {
	cases := []struct {
		name string
		body string
		want error
	}{
		{"schema", `{"schemaVersion":2,"fetchedAt":"2026-11-02T10:00:00Z","usedPercent":1}`, ErrSchema},
		{"percent", `{"schemaVersion":1,"fetchedAt":"2026-11-02T10:00:00Z","usedPercent":101}`, ErrPercent},
		{"time", `{"schemaVersion":1,"fetchedAt":"yesterday","usedPercent":1}`, ErrTime},
		{"json", `{`, ErrJSON},
		{"missing percent", `{"schemaVersion":1,"fetchedAt":"2026-11-02T10:00:00Z","secret":"DROP-ME-EXTRA"}`, ErrPercent},
		{"control in plan", `{"schemaVersion":1,"fetchedAt":"2026-11-02T10:00:00Z","usedPercent":1,"plan":"a\nb"}`, ErrField},
		{"bad breakdown", `{"schemaVersion":1,"fetchedAt":"2026-11-02T10:00:00Z","usedPercent":1,"breakdown":[{"bot":"","sharePercent":1}]}`, ErrBreakdown},
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
	  "fetchedAt": "2026-11-02T10:00:00+08:00",
	  "usedPercent": 6.7,
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
	if file.Snapshot.UsedPercent != 6.7 || file.Path != dst {
		t.Fatalf("file = %+v", file)
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

func TestUsagePathUsesQuosparrDir(t *testing.T) {
	got := UsagePath(filepath.Join(string(os.PathSeparator), "home", "ada"))
	want := filepath.Join(string(os.PathSeparator), "home", "ada", ".quosparr", FileName)
	if got != want {
		t.Fatalf("path = %s, want %s", got, want)
	}
}
