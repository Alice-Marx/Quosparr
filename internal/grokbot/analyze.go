package grokbot

import (
	"sort"
	"time"
)

// BurnMode says where a percent-per-day figure came from.
type BurnMode int

const (
	// BurnNone means the snapshot cannot support a rate.
	BurnNone BurnMode = iota
	// BurnHistory is a line through two or more samples.
	BurnHistory
	// BurnCoarse is usedPercent divided by the time since periodStart.
	BurnCoarse
)

// Burn is the rate section of the assistant panel.
type Burn struct {
	Mode          BurnMode
	PercentPerDay float64
	ETA           time.Time
	HasETA        bool
	Full          bool
}

// BurnOf estimates how fast usedPercent is moving. History wins when two
// samples are at least a minute apart. Otherwise a known periodStart before
// fetchedAt yields a coarse rate. A negative rate is returned without an ETA.
func BurnOf(s Snapshot) Burn {
	points := make([]HistoryPoint, 0, len(s.History)+1)
	points = append(points, s.History...)
	points = append(points, HistoryPoint{At: s.FetchedAt, UsedPercent: s.UsedPercent})
	sort.SliceStable(points, func(i, j int) bool { return points[i].At.Before(points[j].At) })
	deduped := points[:0]
	for _, p := range points {
		if p.At.IsZero() {
			continue
		}
		if n := len(deduped); n > 0 && deduped[n-1].At.Equal(p.At) {
			deduped[n-1] = p
			continue
		}
		deduped = append(deduped, p)
	}
	if len(deduped) >= 2 {
		first, last := deduped[0], deduped[len(deduped)-1]
		if d := last.At.Sub(first.At); d >= time.Minute {
			days := d.Hours() / 24
			return burnFrom((last.UsedPercent-first.UsedPercent)/days, s.UsedPercent, s.FetchedAt, BurnHistory)
		}
	}
	if s.HasPeriodStart && s.PeriodStart.Before(s.FetchedAt) {
		if d := s.FetchedAt.Sub(s.PeriodStart); d >= time.Minute {
			days := d.Hours() / 24
			if days > 0 {
				return burnFrom(s.UsedPercent/days, s.UsedPercent, s.FetchedAt, BurnCoarse)
			}
		}
	}
	return Burn{Mode: BurnNone}
}

func burnFrom(rate, used float64, at time.Time, mode BurnMode) Burn {
	b := Burn{Mode: mode, PercentPerDay: rate, Full: used >= 100}
	if b.Full || rate <= 0 || used >= 100 {
		return b
	}
	days := (100 - used) / rate
	b.ETA = at.Add(time.Duration(days * float64(24*time.Hour)))
	b.HasETA = true
	return b
}

// PaceKind is whether the current pace fits the remaining window.
type PaceKind int

const (
	// PaceMissing means a required input is absent.
	PaceMissing PaceKind = iota
	// PaceOver means the linear projection passes 100 before reset.
	PaceOver
	// PaceUnder means the projection stays comfortably under 100.
	PaceUnder
	// PaceEven means the projection is close to the window.
	PaceEven
)

// Pace projects usedPercent across the quota window.
type Pace struct {
	Kind        PaceKind
	Missing     []string
	Elapsed     float64
	Projected   float64
	WindowStart time.Time
	WindowEnd   time.Time
	AssumedWeek bool
}

// PaceOf compares usedPercent with the fraction of the window that has
// elapsed at fetchedAt. The window ends at resetAt. Its start is periodStart
// when that is earlier than resetAt; otherwise a this_week (or unlabeled)
// period assumes the seven days ending at resetAt.
func PaceOf(s Snapshot) Pace {
	if !s.HasReset {
		return Pace{Kind: PaceMissing, Missing: []string{"resetAt"}}
	}
	var start time.Time
	assumed := false
	switch {
	case s.HasPeriodStart && s.PeriodStart.Before(s.ResetAt):
		start = s.PeriodStart
	case s.PeriodLabel == "" || s.PeriodLabel == "this_week":
		start = s.ResetAt.Add(-7 * 24 * time.Hour)
		assumed = true
	default:
		return Pace{Kind: PaceMissing, Missing: []string{"periodStart"}}
	}
	window := s.ResetAt.Sub(start)
	if window <= time.Minute {
		return Pace{Kind: PaceMissing, Missing: []string{"resetAt"}}
	}
	elapsed := s.FetchedAt.Sub(start)
	if elapsed < 0 {
		elapsed = 0
	}
	if elapsed > window {
		elapsed = window
	}
	frac := float64(elapsed) / float64(window)
	if frac < 0.001 {
		return Pace{Kind: PaceMissing, Missing: []string{"periodStart"}, AssumedWeek: assumed, WindowStart: start, WindowEnd: s.ResetAt}
	}
	projected := s.UsedPercent / frac
	kind := PaceEven
	switch {
	case projected > 100:
		kind = PaceOver
	case projected < 90:
		kind = PaceUnder
	}
	return Pace{
		Kind: kind, Elapsed: frac, Projected: projected,
		WindowStart: start, WindowEnd: s.ResetAt, AssumedWeek: assumed,
	}
}
