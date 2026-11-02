package app

import (
	"cmp"
	"errors"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"
	"github.com/licoy/quosparr/internal/grokbot"
	"github.com/licoy/quosparr/internal/l10n"
	"github.com/licoy/quosparr/internal/runway"
	"github.com/licoy/quosparr/internal/store"
)

// panelGrokBot is the assistant quota page. It is not the Grok CLI panel.
func panelGrokBot(c *ui.Context, a *App, l l10n.Localizer, prefs store.Preferences, st *panelState) {
	s := a.grokBotScreen()
	first := true
	sectionBlock(c, &first, func() {
		spec := sectionSpec{
			title:   tx(l, xGrokBotTitle),
			caption: grokBotPeriod(l, s.Snap.PeriodLabel),
			refresh: func() { a.refreshGrokBot() },
			trailing: func() {
				if SmallPillButton(c, nil, tx(l, xGrokBotImport), false).Clicked() {
					a.importGrokBot()
				}
			},
		}
		drawSectionHead(c, l, st, spec)
		if s.Notice != nil {
			ui.Text(c, tx(l, xGrokBotImportFailed)).FontSize(fsCallout).FontWeight(600).TextColor(c.Theme().Danger)
			dimText(c, grokBotErrText(l, s.Notice))
		}
		switch {
		case s.Err != nil:
			grokBotMessage(c, l, tx(l, xGrokBotInvalid), grokBotErrText(l, s.Err), s.Path)
		case !s.ready():
			grokBotMessage(c, l, tx(l, xGrokBotEmptyTitle), tx(l, xGrokBotEmptyDetail, s.Path), "")
		default:
			grokBotReading(c, l, s, st.resetSwitch(prefs.QuotaResetTimeFormat))
		}
	})
	if s.ready() {
		grokBotSections(c, l, s, st.resetSwitch(prefs.QuotaResetTimeFormat), &first)
	}
}

func grokBotIdentityRow(c *ui.Context, l l10n.Localizer, a *App) {
	s := a.grokBotScreen()
	name := s.Snap.Plan
	if name == "" {
		name = tx(l, xGrokBotTab)
	}
	t := c.Theme()
	ui.Row(c).Padding(3, 6).Gap(6).Margin(8, -6, 0, -6).AlignItems(ui.Center).Children(func() {
		ui.Icon(c, icChat).Size(14, 14).TextColor(inkSecondary(t))
		ui.Text(c, name).FontSize(fsCallout).FontWeight(600).SingleLine().Ellipsis("…").MaxWidth(240)
	})
}

func grokBotMessage(c *ui.Context, l l10n.Localizer, title, detail, path string) {
	t := c.Theme()
	cardColumn(c, surfaceSunkenKind, radiusCard).Gap(7).Padding(12).Children(func() {
		ui.Row(c).Gap(6).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icChat).Size(13, 13).TextColor(inkSecondary(t))
			ui.Text(c, title).FontSize(fsCallout).FontWeight(600)
		})
		dimText(c, detail)
		if path != "" {
			dimText(c, path)
		}
		dimText(c, tx(l, xGrokBotSourceNote))
	})
}

func grokBotReading(c *ui.Context, l l10n.Localizer, s grokBotScreen, rs resetSwitch) {
	snap := s.Snap
	now := time.Now()
	used := grokbot.FormatPercent(snap.UsedPercent)
	left := grokbot.FormatPercent(100 - snap.UsedPercent)
	t := c.Theme()
	ui.Column(c).Gap(10).Children(func() {
		dimText(c, tx(l, xGrokBotLocal)+" · "+s.Path)
		ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
			ui.Text(c, tx(l, xGrokBotUsed, used)).FontSize(fsCallout).FontWeight(600).Grow(1).Basis(0)
			ui.Text(c, tx(l, xGrokBotLeft, left)).FontSize(fsCallout).FontWeight(600).FontFeatures("tnum").TextColor(healthInk(t, grokBotHealth(snap.UsedPercent))).SingleLine()
		})
		usedFractionBar(c, snap.UsedPercent/100, grokBotHealth(snap.UsedPercent))
		if snap.HasReset {
			ui.Row(c).Children(func() {
				ui.Spacer(c).Grow(1)
				resetLine(c, l, "reset-grokbot", snap.ResetAt, now, rs)
			})
		}
		if line := grokBotOnDemand(l, snap.OnDemand); line != "" {
			dimText(c, line)
		}
		dimText(c, tx(l, xGrokBotFetched, snap.FetchedAt.Local().Format("2006-01-02 15:04")))
		if len(snap.Breakdown) > 0 {
			grokBotBreakdown(c, l, snap.Breakdown)
		}
		if snap.Notes != "" {
			dimText(c, snap.Notes)
		}
		dimText(c, tx(l, xGrokBotSourceNote))
	})
}

func grokBotBreakdown(c *ui.Context, l l10n.Localizer, items []grokbot.BotShare) {
	t := c.Theme()
	ui.Text(c, tx(l, xGrokBotBreakdown)).FontSize(fsCallout).FontWeight(600)
	ui.Column(c).Gap(8).Children(func() {
		for _, bot := range orderedBots(items) {
			ui.Row(c).Gap(8).AlignItems(ui.Center).Children(func() {
				ui.Text(c, bot.Bot).FontSize(fsCaption).FontWeight(600).Grow(1).Basis(0)
				ui.Text(c, grokbot.FormatPercent(bot.SharePercent)+"%").FontSize(fsCaption).FontWeight(600).FontFeatures("tnum").SingleLine()
			})
			usedFractionBar(c, bot.SharePercent/100, runway.HealthGreen)
			for _, act := range bot.Activities {
				ui.Row(c).Gap(8).Padding(0, 0, 0, 12).AlignItems(ui.Center).Children(func() {
					ui.Text(c, act.Name).FontSize(fsCaption).TextColor(inkSecondary(t)).Grow(1).Basis(0)
					ui.Text(c, grokbot.FormatPercent(act.SharePercent)+"%").FontSize(fsCaption).FontFeatures("tnum").TextColor(inkSecondary(t)).SingleLine()
				})
			}
		}
	})
}

func orderedBots(items []grokbot.BotShare) []grokbot.BotShare {
	out := slices.Clone(items)
	slices.SortStableFunc(out, func(a, b grokbot.BotShare) int {
		return cmp.Compare(b.SharePercent, a.SharePercent)
	})
	if len(out) > 8 {
		out = out[:8]
	}
	for i := range out {
		acts := slices.Clone(out[i].Activities)
		slices.SortStableFunc(acts, func(a, b grokbot.ActivityShare) int {
			return cmp.Compare(b.SharePercent, a.SharePercent)
		})
		if len(acts) > 6 {
			acts = acts[:6]
		}
		out[i].Activities = acts
	}
	return out
}

func usedFractionBar(c *ui.Context, fraction float64, health runway.Health) {
	if fraction < 0 {
		fraction = 0
	}
	if fraction > 1 {
		fraction = 1
	}
	t := c.Theme()
	fill := healthFill(t, health)
	ui.Box(c).Height(quotaBarHeight).Draw(func(p *ui.Painter, r ui.Rect) {
		p.Fill(r, surfaceSunken(t), r.H/2)
		w := max(4, r.W*float32(fraction))
		p.Fill(ui.Rect{X: r.X, Y: r.Y, W: min(w, r.W), H: r.H}, fill, r.H/2)
	})
}

func grokBotHealth(used float64) runway.Health {
	rounded := int(used + 0.5)
	if rounded < 0 {
		rounded = 0
	}
	if rounded > 100 {
		rounded = 100
	}
	return runway.HealthForUsedPercent(rounded)
}

func grokBotPeriod(l l10n.Localizer, label string) string {
	if label == "" || label == "this_week" {
		return tx(l, xGrokBotPeriodWeek)
	}
	return label
}

func grokBotOnDemand(l l10n.Localizer, status string) string {
	var value string
	switch status {
	case "":
		return ""
	case "unlimited":
		value = tx(l, xGrokBotOnDemandUnlimited)
	case "limited":
		value = tx(l, xGrokBotOnDemandLimited)
	case "off", "disabled", "none":
		value = tx(l, xGrokBotOnDemandOff)
	default:
		value = status
	}
	return tx(l, xGrokBotOnDemand, value)
}

func grokBotErrText(l l10n.Localizer, err error) string {
	switch {
	case errors.Is(err, grokbot.ErrJSON):
		return tx(l, xGrokBotErrJSON)
	case errors.Is(err, grokbot.ErrSchema):
		return tx(l, xGrokBotErrSchema)
	case errors.Is(err, grokbot.ErrPercent):
		return tx(l, xGrokBotErrPercent)
	case errors.Is(err, grokbot.ErrTime):
		return tx(l, xGrokBotErrTime)
	case errors.Is(err, grokbot.ErrBreakdown):
		return tx(l, xGrokBotErrBreakdown)
	case errors.Is(err, grokbot.ErrField):
		return tx(l, xGrokBotErrField)
	case errors.Is(err, grokbot.ErrTooLarge), errors.Is(err, grokbot.ErrRead), errors.Is(err, grokbot.ErrMissing):
		return tx(l, xGrokBotErrRead)
	default:
		return tx(l, xGrokBotErrRead)
	}
}
