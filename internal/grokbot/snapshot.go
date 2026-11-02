package grokbot

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/licoy/quosparr/internal/platform"
)

const (
	// SchemaVersion is the only snapshot version this package reads.
	SchemaVersion = 1
	// FileName is the snapshot inside the Quosparr data directory.
	FileName = "grokbot-usage.json"
	// maxSnapshotBytes rejects a file that is larger than a usage snapshot.
	maxSnapshotBytes int64 = 256 * 1024
	maxBots                = 64
	maxActivities          = 64
	maxHistory             = 256
	maxCostLines           = 32
)

// UsagePath is ~/.quosparr/grokbot-usage.json under home
// (%USERPROFILE%\.quosparr\grokbot-usage.json on Windows).
func UsagePath(home string) string {
	return filepath.Join(platform.RunwayDir(home), FileName)
}

// Sentinel errors are safe to show: they never include file bytes, tokens
// or paths.
var (
	ErrMissing   = errors.New("grokbot: snapshot missing")
	ErrTooLarge  = errors.New("grokbot: snapshot is too large")
	ErrRead      = errors.New("grokbot: cannot read snapshot")
	ErrJSON      = errors.New("grokbot: snapshot is not valid JSON")
	ErrSchema    = errors.New("grokbot: unsupported schemaVersion")
	ErrPercent   = errors.New("grokbot: usedPercent out of range")
	ErrTime      = errors.New("grokbot: timestamp is not RFC3339")
	ErrBreakdown = errors.New("grokbot: breakdown is invalid")
	ErrField     = errors.New("grokbot: field is invalid")
)

// Snapshot is schema version 1.
type Snapshot struct {
	SchemaVersion int
	FetchedAt     time.Time
	Plan          string
	UsedPercent   float64
	// ResetAt is set when the snapshot included resetAt.
	ResetAt     time.Time
	HasReset    bool
	OnDemand    string
	PeriodLabel string
	// PeriodStart is the optional start of the quota window.
	PeriodStart    time.Time
	HasPeriodStart bool
	Breakdown      []BotShare
	Notes          string
	History        []HistoryPoint
	Cache          CacheStats
	Usage          UsageStats
	CostLines      []CostLine
}

// BotShare is one bot's share of the weekly usage.
type BotShare struct {
	Bot          string
	SharePercent float64
	Activities   []ActivityShare
}

// ActivityShare is one activity inside a bot.
type ActivityShare struct {
	Name         string
	SharePercent float64
}

// File is a snapshot read from disk.
type File struct {
	Snapshot Snapshot
	Path     string
	ModTime  time.Time
	Size     int64
}

// Load reads and validates path. A missing file is ErrMissing. The returned
// error is one of the sentinel errors.
func Load(path string) (File, error) {
	data, info, err := readLimited(path)
	if err != nil {
		return File{}, err
	}
	snap, err := Parse(data)
	if err != nil {
		return File{}, err
	}
	return File{Snapshot: snap, Path: path, ModTime: info.ModTime(), Size: info.Size()}, nil
}

// Import validates src and atomically replaces dst with the canonical
// snapshot. Unknown JSON fields are dropped. dst is left unchanged when
// validation fails.
func Import(dst, src string) (File, error) {
	data, _, err := readLimited(src)
	if err != nil {
		if errors.Is(err, ErrMissing) {
			return File{}, ErrRead
		}
		return File{}, err
	}
	snap, err := Parse(data)
	if err != nil {
		return File{}, err
	}
	body, err := Marshal(snap)
	if err != nil {
		return File{}, ErrJSON
	}
	if err := platform.WriteFileAtomic(dst, body, 0o600); err != nil {
		return File{}, ErrRead
	}
	return Load(dst)
}

// Parse validates a schema version 1 document. Unknown fields are ignored
// and are not part of the result.
func Parse(data []byte) (Snapshot, error) {
	if int64(len(data)) > maxSnapshotBytes {
		return Snapshot{}, ErrTooLarge
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var raw rawSnapshot
	if err := dec.Decode(&raw); err != nil {
		return Snapshot{}, ErrJSON
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return Snapshot{}, ErrJSON
	}
	return raw.snapshot()
}

// Marshal writes the canonical snapshot. Extra fields are not included.
func Marshal(s Snapshot) ([]byte, error) {
	body, err := json.MarshalIndent(s.encoded(), "", "  ")
	if err != nil {
		return nil, err
	}
	return append(body, '\n'), nil
}

// FormatPercent renders a 0–100 value with at most one decimal place.
func FormatPercent(v float64) string {
	rounded := math.Round(v*10) / 10
	if rounded == math.Trunc(rounded) {
		return strconv.FormatFloat(rounded, 'f', 0, 64)
	}
	return strconv.FormatFloat(rounded, 'f', 1, 64)
}

type rawSnapshot struct {
	SchemaVersion *int          `json:"schemaVersion"`
	FetchedAt     string        `json:"fetchedAt"`
	Plan          string        `json:"plan"`
	UsedPercent   *json.Number  `json:"usedPercent"`
	ResetAt       string        `json:"resetAt"`
	OnDemand      *rawOnDemand  `json:"onDemand"`
	PeriodLabel   string        `json:"periodLabel"`
	PeriodStart   string        `json:"periodStart"`
	Breakdown     []rawBot      `json:"breakdown"`
	Notes         string        `json:"notes"`
	History       []rawHistory  `json:"history"`
	Cache         *rawCache     `json:"cache"`
	Usage         *rawUsage     `json:"usage"`
	CostLines     []rawCostLine `json:"costLines"`
}

type rawOnDemand struct {
	Status string `json:"status"`
}

type rawBot struct {
	Bot          string        `json:"bot"`
	SharePercent *json.Number  `json:"sharePercent"`
	Activities   []rawActivity `json:"activities"`
}

type rawActivity struct {
	Name         string       `json:"name"`
	SharePercent *json.Number `json:"sharePercent"`
}

func (r rawSnapshot) snapshot() (Snapshot, error) {
	if r.SchemaVersion == nil || *r.SchemaVersion != SchemaVersion {
		return Snapshot{}, ErrSchema
	}
	fetched, err := parseTime(r.FetchedAt, true)
	if err != nil {
		return Snapshot{}, err
	}
	used, err := parsePercent(r.UsedPercent, true)
	if err != nil {
		return Snapshot{}, ErrPercent
	}
	reset, hasReset, err := parseOptionalTime(r.ResetAt)
	if err != nil {
		return Snapshot{}, err
	}
	plan, err := optionalLabel(r.Plan, 80)
	if err != nil {
		return Snapshot{}, err
	}
	notes, err := optionalLabel(r.Notes, 400)
	if err != nil {
		return Snapshot{}, err
	}
	period, err := optionalPeriod(r.PeriodLabel)
	if err != nil {
		return Snapshot{}, err
	}
	demand, err := optionalDemand(r.OnDemand)
	if err != nil {
		return Snapshot{}, err
	}
	breakdown, err := parseBreakdown(r.Breakdown)
	if err != nil {
		return Snapshot{}, err
	}
	extra, err := parseExtras(r)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{
		SchemaVersion:  SchemaVersion,
		FetchedAt:      fetched,
		Plan:           plan,
		UsedPercent:    used,
		ResetAt:        reset,
		HasReset:       hasReset,
		OnDemand:       demand,
		PeriodLabel:    period,
		PeriodStart:    extra.PeriodStart,
		HasPeriodStart: extra.HasPeriodStart,
		Breakdown:      breakdown,
		Notes:          notes,
		History:        extra.History,
		Cache:          extra.Cache,
		Usage:          extra.Usage,
		CostLines:      extra.CostLines,
	}, nil
}

func (s Snapshot) encoded() encodedSnapshot {
	out := encodedSnapshot{
		SchemaVersion: SchemaVersion,
		FetchedAt:     s.FetchedAt.Format(time.RFC3339),
		Plan:          s.Plan,
		UsedPercent:   json.Number(FormatPercent(s.UsedPercent)),
		PeriodLabel:   s.PeriodLabel,
		Notes:         s.Notes,
	}
	if s.HasReset {
		out.ResetAt = s.ResetAt.Format(time.RFC3339)
	}
	if s.OnDemand != "" {
		out.OnDemand = &encodedOnDemand{Status: s.OnDemand}
	}
	if s.HasPeriodStart {
		out.PeriodStart = s.PeriodStart.Format(time.RFC3339)
	}
	encodeExtras(&out, s)
	if len(s.Breakdown) > 0 {
		out.Breakdown = make([]encodedBot, len(s.Breakdown))
		for i, bot := range s.Breakdown {
			item := encodedBot{Bot: bot.Bot, SharePercent: json.Number(FormatPercent(bot.SharePercent))}
			if len(bot.Activities) > 0 {
				item.Activities = make([]encodedActivity, len(bot.Activities))
				for j, act := range bot.Activities {
					item.Activities[j] = encodedActivity{Name: act.Name, SharePercent: json.Number(FormatPercent(act.SharePercent))}
				}
			}
			out.Breakdown[i] = item
		}
	}
	return out
}

type encodedSnapshot struct {
	SchemaVersion int               `json:"schemaVersion"`
	FetchedAt     string            `json:"fetchedAt"`
	Plan          string            `json:"plan,omitempty"`
	UsedPercent   json.Number       `json:"usedPercent"`
	ResetAt       string            `json:"resetAt,omitempty"`
	OnDemand      *encodedOnDemand  `json:"onDemand,omitempty"`
	PeriodLabel   string            `json:"periodLabel,omitempty"`
	PeriodStart   string            `json:"periodStart,omitempty"`
	Breakdown     []encodedBot      `json:"breakdown,omitempty"`
	Notes         string            `json:"notes,omitempty"`
	History       []encodedHistory  `json:"history,omitempty"`
	Cache         *encodedCache     `json:"cache,omitempty"`
	Usage         *encodedUsage     `json:"usage,omitempty"`
	CostLines     []encodedCostLine `json:"costLines,omitempty"`
}

type encodedOnDemand struct {
	Status string `json:"status"`
}

type encodedBot struct {
	Bot          string            `json:"bot"`
	SharePercent json.Number       `json:"sharePercent"`
	Activities   []encodedActivity `json:"activities,omitempty"`
}

type encodedActivity struct {
	Name         string      `json:"name"`
	SharePercent json.Number `json:"sharePercent"`
}

func readLimited(path string) ([]byte, os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil, ErrMissing
		}
		return nil, nil, ErrRead
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || info.IsDir() {
		return nil, nil, ErrRead
	}
	if info.Size() > maxSnapshotBytes {
		return nil, nil, ErrTooLarge
	}
	data, err := io.ReadAll(io.LimitReader(f, maxSnapshotBytes+1))
	if err != nil {
		return nil, nil, ErrRead
	}
	if int64(len(data)) > maxSnapshotBytes {
		return nil, nil, ErrTooLarge
	}
	return data, info, nil
}

func parseTime(raw string, required bool) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		if required {
			return time.Time{}, ErrTime
		}
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, ErrTime
	}
	return t, nil
}

func parseOptionalTime(raw string) (time.Time, bool, error) {
	if strings.TrimSpace(raw) == "" {
		return time.Time{}, false, nil
	}
	t, err := parseTime(raw, true)
	if err != nil {
		return time.Time{}, false, err
	}
	return t, true, nil
}

func parsePercent(n *json.Number, required bool) (float64, error) {
	if n == nil {
		if required {
			return 0, ErrPercent
		}
		return 0, ErrBreakdown
	}
	v, err := n.Float64()
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v < 0 || v > 100 {
		return 0, ErrPercent
	}
	return v, nil
}

func optionalLabel(raw string, maxRunes int) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}
	if !validLabel(raw, maxRunes) {
		return "", ErrField
	}
	return raw, nil
}

func validLabel(raw string, maxRunes int) bool {
	if raw == "" || utf8.RuneCountInString(raw) > maxRunes {
		return false
	}
	for _, r := range raw {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

var (
	periodPattern = regexp.MustCompile(`^[a-z0-9_]{1,32}$`)
	statusPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,32}$`)
)

func optionalPeriod(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}
	if !periodPattern.MatchString(raw) {
		return "", ErrField
	}
	return raw, nil
}

func optionalDemand(raw *rawOnDemand) (string, error) {
	if raw == nil {
		return "", nil
	}
	status := strings.TrimSpace(raw.Status)
	if !statusPattern.MatchString(status) {
		return "", ErrField
	}
	return status, nil
}

func parseBreakdown(raw []rawBot) ([]BotShare, error) {
	if len(raw) > maxBots {
		return nil, ErrBreakdown
	}
	out := make([]BotShare, 0, len(raw))
	for _, item := range raw {
		name, err := optionalLabel(item.Bot, 80)
		if err != nil || name == "" {
			return nil, ErrBreakdown
		}
		share, err := parsePercent(item.SharePercent, false)
		if err != nil {
			return nil, ErrBreakdown
		}
		if len(item.Activities) > maxActivities {
			return nil, ErrBreakdown
		}
		acts := make([]ActivityShare, 0, len(item.Activities))
		for _, act := range item.Activities {
			actName, err := optionalLabel(act.Name, 80)
			if err != nil || actName == "" {
				return nil, ErrBreakdown
			}
			actShare, err := parsePercent(act.SharePercent, false)
			if err != nil {
				return nil, ErrBreakdown
			}
			acts = append(acts, ActivityShare{Name: actName, SharePercent: actShare})
		}
		out = append(out, BotShare{Bot: name, SharePercent: share, Activities: acts})
	}
	return out, nil
}
