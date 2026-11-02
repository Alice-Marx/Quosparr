package app

import (
	"context"
	"errors"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/egoist/mygo"
	"github.com/licoy/quosparr/internal/grokbot"
	"github.com/licoy/quosparr/internal/platform"
	"github.com/licoy/quosparr/internal/runway"
	"github.com/licoy/quosparr/internal/store"
)

// grokBotHold is the assistant snapshot. The panel reads it; the watcher
// and the import button are the only writers. It never contains tokens.
type grokBotHold struct {
	mu      sync.Mutex
	path    string
	snap    grokbot.Snapshot
	missing bool
	err     error
	notice  error
}

type grokBotScreen struct {
	Path    string
	Snap    grokbot.Snapshot
	Missing bool
	Err     error
	Notice  error
}

func (s grokBotScreen) ready() bool {
	return s.Err == nil && s.Snap.SchemaVersion == grokbot.SchemaVersion
}

func (a *App) grokBotScreen() grokBotScreen {
	a.grokBot.mu.Lock()
	s := grokBotScreen{Path: a.grokBot.path, Snap: a.grokBot.snap, Missing: a.grokBot.missing, Err: a.grokBot.err, Notice: a.grokBot.notice}
	a.grokBot.mu.Unlock()
	if s.Path == "" {
		if home, err := platform.Home(); err == nil {
			s.Path = grokbot.UsagePath(home)
		}
	}
	if s.Path == "" {
		s.Path = grokbot.FileName
	}
	if !s.ready() && s.Err == nil {
		s.Missing = true
	}
	return s
}

func (a *App) grokBotPath() string {
	a.grokBot.mu.Lock()
	path := a.grokBot.path
	a.grokBot.mu.Unlock()
	if path != "" {
		return path
	}
	home, err := platform.Home()
	if err != nil {
		return ""
	}
	return grokbot.UsagePath(home)
}

// watchGrokBot reloads the snapshot when its modification time or size
// changes. Manual refresh does not wait for the next tick.
func (a *App) watchGrokBot(ctx context.Context, home string) {
	path := grokbot.UsagePath(home)
	a.grokBot.mu.Lock()
	a.grokBot.path = path
	a.grokBot.mu.Unlock()
	a.refreshGrokBot()
	prev := grokBotSignature(path)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sig := grokBotSignature(path)
			if sig == prev {
				continue
			}
			prev = sig
			a.refreshGrokBot()
		}
	}
}

func (a *App) refreshGrokBot() {
	a.loadGrokBot()
	runOnMain(func() { a.paintGrokBot() })
}

func (a *App) loadGrokBot() {
	path := a.grokBotPath()
	file, err := grokbot.Load(path)
	a.grokBot.mu.Lock()
	defer a.grokBot.mu.Unlock()
	if path != "" {
		a.grokBot.path = path
	}
	a.grokBot.notice = nil
	if err != nil {
		a.grokBot.snap = grokbot.Snapshot{}
		if errors.Is(err, grokbot.ErrMissing) {
			a.grokBot.missing = true
			a.grokBot.err = nil
			return
		}
		a.grokBot.missing = false
		a.grokBot.err = err
		log.Println("grokbot snapshot: unreadable")
		return
	}
	a.grokBot.snap = file.Snapshot
	a.grokBot.missing = false
	a.grokBot.err = nil
}

func (a *App) paintGrokBot() {
	if a.panel != nil {
		a.panel.Invalidate()
	}
	if view, ok := a.grokBotTray(); ok && a.tray != nil {
		ApplyTray(a.tray, TrayFromStatus(view))
	}
}

func (a *App) importGrokBot() {
	title := tx(a.L(), xGrokBotImport)
	dst := a.grokBotPath()
	go func() {
		paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
			Title: title,
			Filters: []mygo.FileFilter{{
				Name:       "JSON",
				Extensions: []string{"json"},
			}},
		})
		if err != nil || len(paths) == 0 || dst == "" {
			if dst == "" && err == nil && len(paths) > 0 {
				runOnMain(func() { a.finishGrokBotImport(grokbot.File{}, grokbot.ErrRead) })
			}
			return
		}
		file, err := grokbot.Import(dst, paths[0])
		runOnMain(func() { a.finishGrokBotImport(file, err) })
	}()
}

func (a *App) finishGrokBotImport(file grokbot.File, err error) {
	a.grokBot.mu.Lock()
	if err != nil {
		a.grokBot.notice = err
	} else {
		a.grokBot.path = file.Path
		a.grokBot.snap = file.Snapshot
		a.grokBot.missing = false
		a.grokBot.err = nil
		a.grokBot.notice = nil
	}
	a.grokBot.mu.Unlock()
	a.paintGrokBot()
}

func (a *App) grokBotMeters(now time.Time) []runway.Meter {
	s := a.grokBotScreen()
	if !s.ready() {
		return nil
	}
	used := int(math.Round(s.Snap.UsedPercent))
	if used < 0 {
		used = 0
	}
	if used > 100 {
		used = 100
	}
	window := runway.MeterWindow{UsedPercent: used, WindowMinutes: runway.WeeklyMinutes}
	if s.Snap.HasReset {
		window.ResetsAt = s.Snap.ResetAt
	}
	return []runway.Meter{runway.NewMeter("", window, now, nil, runway.SourceStandard)}
}

func (a *App) grokBotTray() (runway.StatusBarView, bool) {
	if a.settings == nil {
		return runway.StatusBarView{}, false
	}
	prefs := a.settings.Prefs()
	if prefs.SelectedProvider != store.ProviderGrokBot || prefs.StatusBarGroups() != nil {
		return runway.StatusBarView{}, false
	}
	now := time.Now()
	in := runway.StatusInputs{
		Prefs: prefs, Provider: store.ProviderGrokBot, Now: now,
		GrokBotMeters: a.grokBotMeters(now),
	}
	return runway.BuildStatusBar(a.L(), in), true
}

func (a *App) overlayGrokBotStatus(snap runway.Snapshot) runway.Snapshot {
	if view, ok := a.grokBotTray(); ok {
		snap.StatusBar = view
	}
	return snap
}

type grokBotSig struct {
	exists bool
	mod    int64
	size   int64
}

func grokBotSignature(path string) grokBotSig {
	info, err := os.Stat(path)
	if err != nil {
		return grokBotSig{}
	}
	return grokBotSig{exists: true, mod: info.ModTime().UnixNano(), size: info.Size()}
}
