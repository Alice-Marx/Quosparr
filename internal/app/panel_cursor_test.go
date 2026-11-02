package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	for _, want := range []string{"Cursor Pro", "42.5", "Cloud Agents", "60", tx(l, xGrokBotOnDemandUnlimited), tx(l, xCursorTitle), tx(l, xGrokBotBurnTitle), tx(l, xCursorPaceTitle)} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
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
