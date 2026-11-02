package app

import (
	"cmp"
	"errors"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/licoy/quosparr/internal/cursorusage"
	"github.com/licoy/quosparr/internal/l10n"
	"github.com/licoy/quosparr/internal/runway"
	"github.com/licoy/quosparr/internal/store"
)

// panelCursor is the Cursor product usage page. It is not the Grok Bot
// assistant page and not the Grok CLI page.
func panelCursor(c *ui.Context, a *App, l l10n.Localizer, prefs store.Preferences, st *panelState) {
	s := a.cursorScreen()
	first := true
	sectionBlock(c, &first, func() {
		spec := sectionSpec{
			title:   tx(l, xCursorTitle),
			caption: cursorPeriod(l, s.Snap.Included.PeriodLabel),
			refresh: func() { a.refreshCursor() },
			trailing: func() {
				if SmallPillButton(c, nil, tx(l, xCursorImport), false).Clicked() {
					a.importCursor()
				}
			},
		}
		drawSectionHead(c, l, st, spec)
		if s.Notice != nil {
			ui.Text(c, tx(l, xGrokBotImportFailed)).FontSize(fsCallout).FontWeight(600).TextColor(c.Theme().Danger)
			dimText(c, cursorErrText(l, s.Notice))
		}
		switch {
		case s.Err != nil:
			cursorMessage(c, l, tx(l, xGrokBotInvalid), cursorErrText(l, s.Err), s.Path)
		case !s.ready():
			cursorMessage(c, l, tx(l, xCursorEmptyTitle), tx(l, xCursorEmptyDetail, s.Path), "")
		default:
			cursorReading(c, l, s, st.resetSwitch(prefs.QuotaResetTimeFormat))
		}
	})
	if s.ready() {
		cursorSections(c, l, s, st.resetSwitch(prefs.QuotaResetTimeFormat), &first)
	}
}

func cursorIdentityRow(c *ui.Context, l l10n.Localizer, a *App) {
	s := a.cursorScreen()
	name := s.Snap.Plan
	if name == "" {
		name = tx(l, xCursorTab)
	}
	t := c.Theme()
	ui.Row(c).Padding(3, 6).Gap(6).Margin(8, -6, 0, -6).AlignItems(ui.Center).Children(func() {
		ui.Icon(c, icGauge).Size(14, 14).TextColor(inkSecondary(t))
		ui.Text(c, name).FontSize(fsCallout).FontWeight(600).SingleLine().Ellipsis("…").MaxWidth(240)
	})
}

func cursorMessage(c *ui.Context, l l10n.Localizer, title, detail, path string) {
	t := c.Theme()
	cardColumn(c, surfaceSunkenKind, radiusCard).Gap(7).Padding(12).Children(func() {
		ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icGauge).Size(13, 13).TextColor(inkSecondary(t))
			ui.Text(c, title).FontSize(fsCallout).FontWeight(600)
		})
		dimText(c, detail)
		if path != "" {
			dimText(c, path)
		}
		dimText(c, tx(l, xCursorSourceNote))
	})
}

func cursorReading(c *ui.Context, l l10n.Localizer, s cursorScreen, rs resetSwitch) {
	snap := s.Snap
	now := time.Now()
	used := cursorusage.FormatPercent(snap.Included.UsedPercent)
	left := cursorusage.FormatPercent(100 - snap.Included.UsedPercent)
	t := c.Theme()
	ui.Column(c).Gap(10).Children(func() {
		dimText(c, tx(l, xCursorLocal)+" · "+s.Path)
		if snap.HasModelPools() {
			cursorModelPools(c, l, snap.IncludedPools)
		} else {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, tx(l, xGrokBotUsed, used)).FontSize(fsCallout).FontWeight(600).Grow(1).Basis(0)
				ui.Text(c, tx(l, xGrokBotLeft, left)).FontSize(fsCallout).FontWeight(600).FontFeatures("tnum").TextColor(healthInk(t, grokBotHealth(snap.Included.UsedPercent))).SingleLine()
			})
			usedFractionBar(c, snap.Included.UsedPercent/100, grokBotHealth(snap.Included.UsedPercent))
		}
		if snap.Included.HasReset {
			ui.Row(c).Children(func() {
				ui.Spacer(c).Grow(1)
				resetLine(c, l, "reset-cursor", snap.Included.ResetAt, now, rs)
			})
		}
		cursorOnDemandLines(c, l, snap.OnDemand)
		if snap.GrokBotWeeklyRef {
			dimText(c, tx(l, xCursorWeeklyRef))
		}
		dimText(c, tx(l, xGrokBotFetched, snap.FetchedAt.Local().Format("2006-01-02 15:04")))
		if len(snap.Breakdown) > 0 {
			cursorBreakdown(c, l, snap.Breakdown)
		}
		if snap.Notes != "" {
			dimText(c, snap.Notes)
		}
		dimText(c, tx(l, xCursorSourceNote))
	})
}

func cursorOnDemandLines(c *ui.Context, l l10n.Localizer, demand cursorusage.OnDemand) {
	if line := grokBotOnDemand(l, demand.Status); line != "" {
		dimText(c, line)
	}
	switch {
	case demand.HasUsed && demand.HasLimit:
		dimText(c, tx(l, xCursorOnDemandRange, cursorusage.FormatUSD(demand.UsedUSD), cursorusage.FormatUSD(demand.LimitUSD)))
	case demand.HasUsed:
		dimText(c, tx(l, xCursorOnDemandUsed, cursorusage.FormatUSD(demand.UsedUSD)))
	case demand.HasLimit:
		dimText(c, tx(l, xCursorOnDemandCap, cursorusage.FormatUSD(demand.LimitUSD)))
	}
}

// cursorModelPools is the Included-in-Pro strip: one progress row per pool,
// separate from the optional product breakdown below.
func cursorModelPools(c *ui.Context, l l10n.Localizer, pools cursorusage.IncludedPools) {
	ui.Column(c).Gap(10).Children(func() {
		if pools.CursorModels.Present {
			cursorPoolRow(c, cursorPoolLabel(l, pools.CursorModels, xCursorModelsLabel), cursorModelsNote(l, pools.CursorModels), pools.CursorModels.UsedPercent)
		}
		if pools.OtherModels.Present {
			cursorPoolRow(c, cursorPoolLabel(l, pools.OtherModels, xCursorOtherModelsLabel), pools.OtherModels.Note, pools.OtherModels.UsedPercent)
		}
	})
}

func cursorPoolRow(c *ui.Context, label, note string, used float64) {
	t := c.Theme()
	ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
		ui.Text(c, label).FontSize(fsCallout).FontWeight(600).Grow(1).Basis(0)
		ui.Text(c, cursorusage.FormatPercent(used)+"%").FontSize(fsCallout).FontWeight(600).FontFeatures("tnum").TextColor(healthInk(t, grokBotHealth(used))).SingleLine()
	})
	usedFractionBar(c, used/100, grokBotHealth(used))
	if note != "" {
		dimText(c, note)
	}
}

func cursorPoolLabel(l l10n.Localizer, pool cursorusage.Pool, fallback textKey) string {
	if pool.Label != "" {
		return pool.Label
	}
	return tx(l, fallback)
}

func cursorModelsNote(l l10n.Localizer, pool cursorusage.Pool) string {
	if pool.Note != "" {
		return pool.Note
	}
	return tx(l, xCursorModelsNote)
}

func cursorBreakdown(c *ui.Context, l l10n.Localizer, items []cursorusage.Share) {
	ui.Text(c, tx(l, xGrokBotBreakdown)).FontSize(fsCallout).FontWeight(600)
	ui.Column(c).Gap(8).Children(func() {
		for _, item := range orderedCursorShares(items) {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, item.Name).FontSize(fsCaption).FontWeight(600).Grow(1).Basis(0)
				ui.Text(c, cursorusage.FormatPercent(item.SharePercent)+"%").FontSize(fsCaption).FontWeight(600).FontFeatures("tnum").SingleLine()
			})
			usedFractionBar(c, item.SharePercent/100, runway.HealthGreen)
		}
	})
}

func orderedCursorShares(items []cursorusage.Share) []cursorusage.Share {
	out := slices.Clone(items)
	slices.SortStableFunc(out, func(a, b cursorusage.Share) int {
		return cmp.Compare(b.SharePercent, a.SharePercent)
	})
	if len(out) > 8 {
		out = out[:8]
	}
	return out
}

func cursorPeriod(l l10n.Localizer, label string) string {
	switch label {
	case "", "this_week":
		return tx(l, xGrokBotPeriodWeek)
	case "this_month":
		return tx(l, xCursorPeriodMonth)
	case "Included in Pro":
		return tx(l, xCursorIncludedInPro)
	default:
		return label
	}
}

func cursorErrText(l l10n.Localizer, err error) string {
	switch {
	case errors.Is(err, cursorusage.ErrJSON):
		return tx(l, xGrokBotErrJSON)
	case errors.Is(err, cursorusage.ErrSchema):
		return tx(l, xGrokBotErrSchema)
	case errors.Is(err, cursorusage.ErrPercent):
		return tx(l, xGrokBotErrPercent)
	case errors.Is(err, cursorusage.ErrTime):
		return tx(l, xGrokBotErrTime)
	case errors.Is(err, cursorusage.ErrBreakdown):
		return tx(l, xGrokBotErrBreakdown)
	case errors.Is(err, cursorusage.ErrField):
		return tx(l, xGrokBotErrField)
	case errors.Is(err, cursorusage.ErrTooLarge), errors.Is(err, cursorusage.ErrRead), errors.Is(err, cursorusage.ErrMissing):
		return tx(l, xGrokBotErrRead)
	default:
		return tx(l, xGrokBotErrRead)
	}
}
