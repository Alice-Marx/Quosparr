package app

import (
	"strconv"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/licoy/quosparr/internal/grokbot"
	"github.com/licoy/quosparr/internal/l10n"
)

func grokBotSections(c *ui.Context, l l10n.Localizer, s grokBotScreen, rs resetSwitch, first *bool) {
	snap := s.Snap
	now := time.Now()
	grokBotBurnSection(c, l, snap, first)
	grokBotResetSection(c, l, snap, now, rs, first)
	grokBotPaceSection(c, l, snap, first)
	grokBotCacheSection(c, l, snap, first)
	grokBotCostSection(c, l, snap, first)
	grokBotRepairSection(c, l, first)
}

func grokBotBurnSection(c *ui.Context, l l10n.Localizer, snap grokbot.Snapshot, first *bool) {
	sectionBlock(c, first, func() {
		drawSectionHead(c, l, nil, sectionSpec{title: tx(l, xGrokBotBurnTitle)})
		burn := grokbot.BurnOf(snap)
		switch burn.Mode {
		case grokbot.BurnNone:
			dimText(c, tx(l, xGrokBotBurnEmpty))
		default:
			rate := grokbot.FormatPercent(burn.PercentPerDay)
			ui.Text(c, tx(l, xGrokBotBurnRate, rate)).FontSize(fsCallout).FontWeight(600).FontFeatures("tnum")
			if burn.Mode == grokbot.BurnCoarse {
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

func grokBotResetSection(c *ui.Context, l l10n.Localizer, snap grokbot.Snapshot, now time.Time, rs resetSwitch, first *bool) {
	sectionBlock(c, first, func() {
		drawSectionHead(c, l, nil, sectionSpec{title: tx(l, xGrokBotResetTitle), caption: grokBotPeriod(l, snap.PeriodLabel)})
		if !snap.HasReset {
			dimText(c, tx(l, xGrokBotResetMissing))
		} else {
			ui.Text(c, snap.ResetAt.Local().Format("2006-01-02 15:04")).FontSize(fsCallout).FontWeight(600).FontFeatures("tnum")
			ui.Row(c).Children(func() {
				ui.Spacer(c).Grow(1)
				resetLine(c, l, "reset-grokbot-section", snap.ResetAt, now, rs)
			})
		}
		if snap.HasPeriodStart {
			dimText(c, tx(l, xGrokBotPeriodStart, snap.PeriodStart.Local().Format("2006-01-02 15:04")))
		}
		if line := grokBotOnDemand(l, snap.OnDemand); line != "" {
			dimText(c, line)
		}
	})
}

func grokBotPaceSection(c *ui.Context, l l10n.Localizer, snap grokbot.Snapshot, first *bool) {
	sectionBlock(c, first, func() {
		drawSectionHead(c, l, nil, sectionSpec{title: tx(l, xGrokBotPaceTitle)})
		pace := grokbot.PaceOf(snap)
		if pace.Kind == grokbot.PaceMissing {
			dimText(c, tx(l, xGrokBotPaceMissing, strings.Join(pace.Missing, ", ")))
			return
		}
		ui.Text(c, tx(l, xGrokBotPaceProjected, grokbot.FormatPercent(pace.Projected))).FontSize(fsCallout).FontWeight(600).FontFeatures("tnum")
		dimText(c, tx(l, xGrokBotPaceElapsed, grokbot.FormatPercent(pace.Elapsed*100)))
		switch pace.Kind {
		case grokbot.PaceOver:
			dimText(c, tx(l, xGrokBotPaceOver))
		case grokbot.PaceUnder:
			dimText(c, tx(l, xGrokBotPaceUnder))
		default:
			dimText(c, tx(l, xGrokBotPaceEven))
		}
		if pace.AssumedWeek {
			dimText(c, tx(l, xGrokBotPaceAssumed))
		}
	})
}

func grokBotCacheSection(c *ui.Context, l l10n.Localizer, snap grokbot.Snapshot, first *bool) {
	sectionBlock(c, first, func() {
		drawSectionHead(c, l, nil, sectionSpec{title: tx(l, xGrokBotCacheTitle)})
		if !snap.Cache.Present {
			dimText(c, tx(l, xGrokBotCacheNA))
			return
		}
		cache := snap.Cache
		if cache.HasHit {
			ui.Text(c, tx(l, xGrokBotCacheHit, grokbot.FormatPercent(cache.HitPercent))).FontSize(fsCallout).FontWeight(600).FontFeatures("tnum")
		}
		if cache.HasHits || cache.HasMisses {
			dimText(c, tx(l, xGrokBotCacheCounts, strconv.FormatInt(cache.Hits, 10), strconv.FormatInt(cache.Misses, 10)))
		}
		if cache.Notes != "" {
			dimText(c, cache.Notes)
		}
		dimText(c, tx(l, xGrokBotCacheOptional))
	})
}

func grokBotCostSection(c *ui.Context, l l10n.Localizer, snap grokbot.Snapshot, first *bool) {
	sectionBlock(c, first, func() {
		drawSectionHead(c, l, nil, sectionSpec{title: tx(l, xGrokBotCostTitle)})
		rows, ok := grokbot.CostRows(snap)
		if !ok {
			dimText(c, tx(l, xGrokBotCostEmpty))
			return
		}
		for _, row := range rows {
			name := row.Model
			if row.Label != "" {
				if name != "" {
					name += " · " + row.Label
				} else {
					name = row.Label
				}
			}
			if name == "" {
				name = tx(l, xGrokBotCostUntitled)
			}
			amount := "—"
			if row.Amount != "" {
				amount = row.Amount
			}
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, name).FontSize(fsCaption).FontWeight(600).Grow(1).Basis(0).SingleLine().Ellipsis("…")
				ui.Text(c, amount).FontSize(fsCaption).FontFeatures("tnum").SingleLine()
			})
			dimText(c, tx(l, xGrokBotCostTokens, strconv.FormatInt(row.Input, 10), strconv.FormatInt(row.Output, 10), strconv.FormatInt(row.Cache, 10)))
			if row.Missing {
				dimText(c, tx(l, xGrokBotCostMissingPrice))
			}
		}
	})
}

func grokBotRepairSection(c *ui.Context, l l10n.Localizer, first *bool) {
	sectionBlock(c, first, func() {
		drawSectionHead(c, l, nil, sectionSpec{title: tx(l, xGrokBotRepairTitle)})
		dimText(c, tx(l, xGrokBotRepairNA))
	})
}
