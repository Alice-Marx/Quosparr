package app

import (
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/licoy/quosparr/internal/cursorusage"
	"github.com/licoy/quosparr/internal/l10n"
)

// cursorSections mirrors the assistant burn, reset, and pace blocks. Cursor
// usage has no cache-hit metric, token price table, or local session index.
func cursorSections(c *ui.Context, l l10n.Localizer, s cursorScreen, rs resetSwitch, first *bool) {
	snap := s.Snap
	now := time.Now()
	cursorBurnSection(c, l, snap, first)
	cursorResetSection(c, l, snap, now, rs, first)
	cursorPaceSection(c, l, snap, first)
}

func cursorBurnSection(c *ui.Context, l l10n.Localizer, snap cursorusage.Snapshot, first *bool) {
	sectionBlock(c, first, func() {
		drawSectionHead(c, l, nil, sectionSpec{title: tx(l, xGrokBotBurnTitle)})
		burn := cursorusage.BurnOf(snap)
		switch burn.Mode {
		case cursorusage.BurnNone:
			dimText(c, tx(l, xGrokBotBurnEmpty))
		default:
			rate := cursorusage.FormatPercent(burn.PercentPerDay)
			ui.Text(c, tx(l, xGrokBotBurnRate, rate)).FontSize(fsCallout).FontWeight(600).FontFeatures("tnum")
			if burn.Mode == cursorusage.BurnCoarse {
				dimText(c, tx(l, xGrokBotBurnCoarse))
			} else {
				dimText(c, tx(l, xGrokBotBurnHistory))
			}
			switch {
			case burn.Full:
				dimText(c, tx(l, xGrokBotBurnFull))
			case burn.HasETA:
				dimText(c, tx(l, xGrokBotBurnETA, burn.ETA.Local().Format("2006-01-02 15:04")))
			case burn.PercentPerDay <= 0:
				dimText(c, tx(l, xGrokBotBurnFlat))
			}
		}
	})
}

func cursorResetSection(c *ui.Context, l l10n.Localizer, snap cursorusage.Snapshot, now time.Time, rs resetSwitch, first *bool) {
	sectionBlock(c, first, func() {
		drawSectionHead(c, l, nil, sectionSpec{title: tx(l, xGrokBotResetTitle), caption: cursorPeriod(l, snap.Included.PeriodLabel)})
		if !snap.Included.HasReset {
			dimText(c, tx(l, xGrokBotResetMissing))
		} else {
			ui.Text(c, snap.Included.ResetAt.Local().Format("2006-01-02 15:04")).FontSize(fsCallout).FontWeight(600).FontFeatures("tnum")
			ui.Row(c).Children(func() {
				ui.Spacer(c).Grow(1)
				resetLine(c, l, "reset-cursor-section", snap.Included.ResetAt, now, rs)
			})
		}
		if snap.HasPeriodStart {
			dimText(c, tx(l, xGrokBotPeriodStart, snap.PeriodStart.Local().Format("2006-01-02 15:04")))
		}
		cursorOnDemandLines(c, l, snap.OnDemand)
	})
}

func cursorPaceSection(c *ui.Context, l l10n.Localizer, snap cursorusage.Snapshot, first *bool) {
	sectionBlock(c, first, func() {
		drawSectionHead(c, l, nil, sectionSpec{title: tx(l, xCursorPaceTitle)})
		pace := cursorusage.PaceOf(snap)
		if pace.Kind == cursorusage.PaceMissing {
			dimText(c, tx(l, xGrokBotPaceMissing, strings.Join(pace.Missing, ", ")))
			return
		}
		ui.Text(c, tx(l, xGrokBotPaceProjected, cursorusage.FormatPercent(pace.Projected))).FontSize(fsCallout).FontWeight(600).FontFeatures("tnum")
		dimText(c, tx(l, xGrokBotPaceElapsed, cursorusage.FormatPercent(pace.Elapsed*100)))
		switch pace.Kind {
		case cursorusage.PaceOver:
			dimText(c, tx(l, xGrokBotPaceOver))
		case cursorusage.PaceUnder:
			dimText(c, tx(l, xGrokBotPaceUnder))
		default:
			dimText(c, tx(l, xGrokBotPaceEven))
		}
		if pace.AssumedWeek {
			dimText(c, tx(l, xGrokBotPaceAssumed))
		}
	})
}
