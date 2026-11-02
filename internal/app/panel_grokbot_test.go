package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/licoy/quosparr/internal/grokbot"
	"github.com/licoy/quosparr/internal/l10n"
	"github.com/licoy/quosparr/internal/store"
)

func TestGrokBotPanelIsSeparateFromGrokCLI(t *testing.T) {
	l := l10n.New(l10n.SimplifiedChinese)
	a, tt := newPanel(t, func(p *store.Preferences) {
		p.Language = store.LanguageSimplifiedChinese
	}, panelFixture(), PanelActions{})
	a.SetPanelActions(PanelActions{SelectProvider: func(p store.Provider) { a.settings.SetSelectedProvider(p) }})
	tt.Frame()
	if !hasExact(tt, tx(l, xGrokBotTab)) || !hasExact(tt, l.Text(l10n.KeyProviderGrok)) {
		t.Fatalf("tabs: %q", tt.Texts())
	}
	if err := tt.Click(tx(l, xGrokBotTab)); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if a.settings.Prefs().SelectedProvider != store.ProviderGrokBot {
		t.Fatalf("selected %s", a.settings.Prefs().SelectedProvider)
	}
	if tt.HasText(l.Text(l10n.KeyQuota)) || hasExact(tt, l.Text(l10n.KeyGrokIncludedQuota)) || tt.HasText("Ada Lovelace") {
		t.Fatalf("cli panel leaked: %q", tt.Texts())
	}
	if !tt.HasText(tx(l, xGrokBotEmptyTitle)) || !tt.HasText(tx(l, xGrokBotImport)) || !tt.HasText("grokbot-usage.json") {
		t.Fatalf("empty state: %q", tt.Texts())
	}

	data, err := os.ReadFile(filepath.Join("..", "..", "examples", "grokbot-usage.json"))
	if err != nil {
		t.Fatal(err)
	}
	snap, err := grokbot.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	a.grokBot.mu.Lock()
	a.grokBot.snap = snap
	a.grokBot.missing = false
	a.grokBot.err = nil
	a.grokBot.path = filepath.Join(t.TempDir(), grokbot.FileName)
	a.grokBot.mu.Unlock()
	tt.Frame()
	for _, want := range []string{"Cursor Pro", "6.7", "89", "chats", tx(l, xGrokBotOnDemandUnlimited), tx(l, xGrokBotTitle)} {
		if !tt.HasText(want) {
			t.Fatalf("missing %q in %q", want, tt.Texts())
		}
	}
	if hasExact(tt, l.Text(l10n.KeyGrokIncludedQuota)) {
		t.Fatal("assistant reading used the Grok CLI quota title")
	}
}

func TestGrokBotRefreshReadsTheSnapshotFile(t *testing.T) {
	t.Cleanup(setRunOnMain(func(fn func()) { fn() }))
	dir := t.TempDir()
	dst := filepath.Join(dir, grokbot.FileName)
	src := filepath.Join(dir, "in.json")
	body := []byte(`{"schemaVersion":1,"fetchedAt":"2026-11-02T10:00:00+08:00","usedPercent":6.7,"secret":"DROP-ME-EXTRA"}`)
	if err := os.WriteFile(src, body, 0o600); err != nil {
		t.Fatal(err)
	}
	a, _ := newPanel(t, func(p *store.Preferences) {
		p.Language = store.LanguageSimplifiedChinese
		p.SelectedProvider = store.ProviderGrokBot
	}, panelFixture(), PanelActions{})
	a.grokBot.mu.Lock()
	a.grokBot.path = dst
	a.grokBot.mu.Unlock()
	if _, err := grokbot.Import(dst, src); err != nil {
		t.Fatal(err)
	}
	a.refreshGrokBot()
	s := a.grokBotScreen()
	if !s.ready() || s.Snap.UsedPercent != 6.7 {
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
