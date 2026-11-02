package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/licoy/quosparr/internal/cursorusage"
	"github.com/licoy/quosparr/internal/l10n"
	"github.com/licoy/quosparr/internal/store"
)

func TestCursorPanelIsSeparateFromGrokBotAndGrokCLI(t *testing.T) {
	l := l10n.New(l10n.SimplifiedChinese)
	a, tt := newPanel(t, func(p *store.Preferences) {
		p.Language = store.LanguageSimplifiedChinese
	}, panelFixture(), PanelActions{})
	a.SetPanelActions(PanelActions{SelectProvider: func(p store.Provider) { a.settings.SetSelectedProvider(p) }})
	tt.Frame()
	if !hasExact(tt, tx(l, xCursorTab)) || !hasExact(tt, tx(l, xGrokBotTab)) || !hasExact(tt, l.Text(l10n.KeyProviderGrok)) {
		t.Fatalf("tabs: %q", tt.Texts())
	}
	if err := tt.Click(tx(l, xCursorTab)); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if a.settings.Prefs().SelectedProvider != store.ProviderCursor {
		t.Fatalf("selected %s", a.settings.Prefs().SelectedProvider)
	}
	if tt.HasText(l.Text(l10n.KeyQuota)) || hasExact(tt, l.Text(l10n.KeyGrokIncludedQuota)) || tt.HasText(tx(l, xGrokBotEmptyTitle)) || tt.HasText("Ada Lovelace") {
		t.Fatalf("other panel leaked: %q", tt.Texts())
	}
	if !tt.HasText(tx(l, xCursorEmptyTitle)) || !tt.HasText(tx(l, xCursorImport)) || !tt.HasText("cursor-usage.json") || !tt.HasText("grokbot-usage.json") {
		t.Fatalf("empty state: %q", tt.Texts())
	}

	data, err := os.ReadFile(filepath.Join("..", "..", "examples", "cursor-usage.json"))
	if err != nil {
		t.Fatal(err)
	}
	snap, err := cursorusage.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	a.cursor.mu.Lock()
	a.cursor.snap = snap
	a.cursor.missing = false
	a.cursor.err = nil
	a.cursor.path = filepath.Join(t.TempDir(), cursorusage.FileName)
	a.cursor.mu.Unlock()
	tt.Frame()
	for _, want := range []string{"Cursor Pro", "42.5", "Cursor Models", "Other Models", "Includes Cursor Grok and Composer", "Cloud Agents", "60", tx(l, xGrokBotOnDemandUnlimited), tx(l, xCursorTitle), tx(l, xGrokBotBurnTitle), tx(l, xCursorPaceTitle), tx(l, xCursorIncludedInPro), tx(l, xCursorWeeklyRef)} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
	if tt.HasText(tx(l, xGrokBotUsed, "42.5")) || tt.HasText(tx(l, xCursorModelsLabel)) {
		t.Fatalf("pools fell back to the overall bar or a default label: %q", tt.Texts())
	}
	if hasExact(tt, tx(l, xGrokBotTitle)) || hasExact(tt, l.Text(l10n.KeyGrokIncludedQuota)) || tt.HasText(tx(l, xGrokBotRepairTitle)) {
		t.Fatal("cursor reading used a Grok Bot or Grok CLI section")
	}
}

func TestCursorRefreshReadsTheSnapshotFile(t *testing.T) {
	t.Cleanup(setRunOnMain(func(fn func()) { fn() }))
	dir := t.TempDir()
	dst := filepath.Join(dir, cursorusage.FileName)
	src := filepath.Join(dir, "in.json")
	body := []byte(`{"schemaVersion":1,"fetchedAt":"2026-11-02T13:00:00+08:00","included":{"usedPercent":42.5},"secret":"DROP-ME-EXTRA"}`)
	if err := os.WriteFile(src, body, 0o600); err != nil {
		t.Fatal(err)
	}
	a, _ := newPanel(t, func(p *store.Preferences) {
		p.Language = store.LanguageSimplifiedChinese
		p.SelectedProvider = store.ProviderCursor
	}, panelFixture(), PanelActions{})
	a.cursor.mu.Lock()
	a.cursor.path = dst
	a.cursor.mu.Unlock()
	if _, err := cursorusage.Import(dst, src); err != nil {
		t.Fatal(err)
	}
	a.refreshCursor()
	s := a.cursorScreen()
	if !s.ready() || s.Snap.Included.UsedPercent != 42.5 {
		t.Fatalf("screen = %+v", s)
	}
	written, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(written), "DROP-ME-EXTRA") || strings.Contains(string(written), "secret") {
		t.Fatal("import kept an unknown field")
	}
}

func TestCursorPoolRowsUseDefaultsAndLegacyOverallStays(t *testing.T) {
	l := l10n.New(l10n.SimplifiedChinese)
	a, tt := newPanel(t, func(p *store.Preferences) {
		p.Language = store.LanguageSimplifiedChinese
		p.SelectedProvider = store.ProviderCursor
	}, panelFixture(), PanelActions{})
	a.cursor.mu.Lock()
	a.cursor.path = filepath.Join(t.TempDir(), cursorusage.FileName)
	a.cursor.missing = false
	a.cursor.err = nil
	a.cursor.snap = cursorusage.Snapshot{
		SchemaVersion: 1,
		Plan:          "Pro",
		FetchedAt:     time.Date(2026, 11, 2, 14, 0, 0, 0, time.FixedZone("CST", 8*3600)),
		Included: cursorusage.Included{
			UsedPercent: 14,
			PeriodLabel: "Included in Pro",
		},
		IncludedPools: cursorusage.IncludedPools{
			CursorModels: cursorusage.Pool{Present: true, UsedPercent: 14},
			OtherModels:  cursorusage.Pool{Present: true, UsedPercent: 0},
		},
		GrokBotWeeklyRef: true,
	}
	a.cursor.mu.Unlock()
	tt.Frame()
	for _, want := range []string{tx(l, xCursorModelsLabel), tx(l, xCursorOtherModelsLabel), tx(l, xCursorModelsNote), tx(l, xCursorIncludedInPro), tx(l, xCursorWeeklyRef), "14%", "0%"} {
		if !hasExact(tt, want) && !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
	if tt.HasText("Cursor Models") || tt.HasText("Other Models") || tt.HasText(tx(l, xGrokBotUsed, "14")) {
		t.Fatalf("defaults were replaced by English labels or the overall bar: %q", tt.Texts())
	}

	a.cursor.mu.Lock()
	a.cursor.snap.IncludedPools.CursorModels.Present = false
	a.cursor.mu.Unlock()
	tt.Frame()
	if !hasExact(tt, tx(l, xCursorOtherModelsLabel)) || tt.HasText(tx(l, xCursorModelsLabel)) || tt.HasText(tx(l, xCursorModelsNote)) || tt.HasText(tx(l, xGrokBotUsed, "14")) {
		t.Fatalf("one pool did not stay a dedicated row: %q", tt.Texts())
	}

	a.cursor.mu.Lock()
	a.cursor.snap.IncludedPools = cursorusage.IncludedPools{}
	a.cursor.snap.GrokBotWeeklyRef = false
	a.cursor.snap.Included.PeriodLabel = "this_week"
	a.cursor.mu.Unlock()
	tt.Frame()
	if !tt.HasText(tx(l, xGrokBotUsed, "14")) || !tt.HasText(tx(l, xGrokBotPeriodWeek)) {
		t.Fatalf("legacy reading missing: %q", tt.Texts())
	}
	if tt.HasText(tx(l, xCursorModelsLabel)) || tt.HasText(tx(l, xCursorOtherModelsLabel)) || tt.HasText(tx(l, xCursorWeeklyRef)) {
		t.Fatalf("legacy reading invented pools: %q", tt.Texts())
	}
}
