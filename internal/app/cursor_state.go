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
	"github.com/licoy/quosparr/internal/cursorusage"
	"github.com/licoy/quosparr/internal/platform"
	"github.com/licoy/quosparr/internal/runway"
	"github.com/licoy/quosparr/internal/store"
)

// cursorHold is the Cursor product snapshot. The panel reads it; the watcher
// and the import button are the only writers. It never contains tokens.
type cursorHold struct {
	mu      sync.Mutex
	path    string
	snap    cursorusage.Snapshot
	missing bool
	err     error
	notice  error
}

type cursorScreen struct {
	Path    string
	Snap    cursorusage.Snapshot
	Missing bool
	Err     error
	Notice  error
}

func (s cursorScreen) ready() bool {
	return s.Err == nil && s.Snap.SchemaVersion == cursorusage.SchemaVersion
}

func (a *App) cursorScreen() cursorScreen {
	a.cursor.mu.Lock()
	s := cursorScreen{Path: a.cursor.path, Snap: a.cursor.snap, Missing: a.cursor.missing, Err: a.cursor.err, Notice: a.cursor.notice}
	a.cursor.mu.Unlock()
	if s.Path == "" {
		if home, err := platform.Home(); err == nil {
			s.Path = cursorusage.UsagePath(home)
		}
	}
	if s.Path == "" {
		s.Path = cursorusage.FileName
	}
	if !s.ready() && s.Err == nil {
		s.Missing = true
	}
	return s
}

func (a *App) cursorPath() string {
	a.cursor.mu.Lock()
	path := a.cursor.path
	a.cursor.mu.Unlock()
	if path != "" {
		return path
	}
	home, err := platform.Home()
	if err != nil {
		return ""
	}
	return cursorusage.UsagePath(home)
}

// watchCursor reloads the snapshot when its modification time or size
// changes. Manual refresh does not wait for the next tick.
func (a *App) watchCursor(ctx context.Context, home string) {
	path := cursorusage.UsagePath(home)
	a.cursor.mu.Lock()
	a.cursor.path = path
	a.cursor.mu.Unlock()
	a.refreshCursor()
	prev := cursorSignature(path)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sig := cursorSignature(path)
			if sig == prev {
				continue
			}
			prev = sig
			a.refreshCursor()
		}
	}
}

func (a *App) refreshCursor() {
	a.loadCursor()
	runOnMain(func() { a.paintCursor() })
}

func (a *App) loadCursor() {
	path := a.cursorPath()
	file, err := cursorusage.Load(path)
	a.cursor.mu.Lock()
	defer a.cursor.mu.Unlock()
	if path != "" {
		a.cursor.path = path
	}
	a.cursor.notice = nil
	if err != nil {
		a.cursor.snap = cursorusage.Snapshot{}
		if errors.Is(err, cursorusage.ErrMissing) {
			a.cursor.missing = true
			a.cursor.err = nil
			return
		}
		a.cursor.missing = false
		a.cursor.err = err
		log.Println("cursor usage snapshot: unreadable")
		return
	}
	a.cursor.snap = file.Snapshot
	a.cursor.missing = false
	a.cursor.err = nil
}

func (a *App) paintCursor() {
	if a.panel != nil {
		a.panel.Invalidate()
	}
	if view, ok := a.cursorTray(); ok && a.tray != nil {
		ApplyTray(a.tray, TrayFromStatus(view))
	}
}

func (a *App) importCursor() {
	title := tx(a.L(), xCursorImport)
	dst := a.cursorPath()
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
				runOnMain(func() { a.finishCursorImport(cursorusage.File{}, cursorusage.ErrRead) })
			}
			return
		}
		file, err := cursorusage.Import(dst, paths[0])
		runOnMain(func() { a.finishCursorImport(file, err) })
	}()
}

func (a *App) finishCursorImport(file cursorusage.File, err error) {
	a.cursor.mu.Lock()
	if err != nil {
		a.cursor.notice = err
	} else {
		a.cursor.path = file.Path
		a.cursor.snap = file.Snapshot
		a.cursor.missing = false
		a.cursor.err = nil
		a.cursor.notice = nil
	}
	a.cursor.mu.Unlock()
	a.paintCursor()
}

func (a *App) cursorMeters(now time.Time) []runway.Meter {
	s := a.cursorScreen()
	if !s.ready() {
		return nil
	}
	used := int(math.Round(s.Snap.Included.UsedPercent))
	if used < 0 {
		used = 0
	}
	if used > 100 {
		used = 100
	}
	minutes := runway.WeeklyMinutes
	if s.Snap.Included.PeriodLabel == "this_month" {
		minutes = runway.MonthlyMinutes
	}
	window := runway.MeterWindow{UsedPercent: used, WindowMinutes: minutes}
	if s.Snap.Included.HasReset {
		window.ResetsAt = s.Snap.Included.ResetAt
	}
	return []runway.Meter{runway.NewMeter("", window, now, nil, runway.SourceStandard)}
}

func (a *App) cursorTray() (runway.StatusBarView, bool) {
	if a.settings == nil {
		return runway.StatusBarView{}, false
	}
	prefs := a.settings.Prefs()
	if prefs.SelectedProvider != store.ProviderCursor || prefs.StatusBarGroups() != nil {
		return runway.StatusBarView{}, false
	}
	now := time.Now()
	in := runway.StatusInputs{
		Prefs: prefs, Provider: store.ProviderCursor, Now: now,
		CursorMeters: a.cursorMeters(now),
	}
	return runway.BuildStatusBar(a.L(), in), true
}

func (a *App) overlayCursorStatus(snap runway.Snapshot) runway.Snapshot {
	if view, ok := a.cursorTray(); ok {
		snap.StatusBar = view
	}
	return snap
}

type cursorSig struct {
	exists bool
	mod    int64
	size   int64
}

func cursorSignature(path string) cursorSig {
	info, err := os.Stat(path)
	if err != nil {
		return cursorSig{}
	}
	return cursorSig{exists: true, mod: info.ModTime().UnixNano(), size: info.Size()}
}
