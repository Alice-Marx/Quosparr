package cursorusage

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
	FileName = "cursor-usage.json"
	// maxSnapshotBytes rejects a file that is larger than a usage snapshot.
	maxSnapshotBytes int64 = 256 * 1024
	maxShares              = 64
	maxHistory             = 256
	// maxUSD rejects a money field that cannot be a usage amount.
	maxUSD = 1e9
)

// UsagePath is ~/.quosparr/cursor-usage.json under home
// (%USERPROFILE%\.quosparr\cursor-usage.json on Windows).
func UsagePath(home string) string {
	return filepath.Join(platform.RunwayDir(home), FileName)
}

// Sentinel errors are safe to show: they never include file bytes, tokens
// or paths.
var (
	ErrMissing   = errors.New("cursorusage: snapshot missing")
	ErrTooLarge  = errors.New("cursorusage: snapshot is too large")
	ErrRead      = errors.New("cursorusage: cannot read snapshot")
	ErrJSON      = errors.New("cursorusage: snapshot is not valid JSON")
	ErrSchema    = errors.New("cursorusage: unsupported schemaVersion")
	ErrPercent   = errors.New("cursorusage: usedPercent out of range")
	ErrTime      = errors.New("cursorusage: timestamp is not RFC3339")
	ErrBreakdown = errors.New("cursorusage: breakdown is invalid")
	ErrField     = errors.New("cursorusage: field is invalid")
)

// Snapshot is schema version 1 of the Cursor product usage file.
type Snapshot struct {
	SchemaVersion int
	FetchedAt     time.Time
	Plan          string
	Included      Included
	OnDemand      OnDemand
	Breakdown     []Share
	Notes         string
	// PeriodStart is the optional start of the included window.
	PeriodStart    time.Time
	HasPeriodStart bool
	History        []HistoryPoint
}

// Included is the plan allowance shown on the Cursor usage page.
type Included struct {
	UsedPercent float64
	ResetAt     time.Time
	HasReset    bool
	PeriodLabel string
}

// OnDemand is the on-demand allowance. Present is false when the object was
// omitted. UsedUSD and LimitUSD are set only when the file included a number;
// JSON null is the same as omitted.
type OnDemand struct {
	Present  bool
	Status   string
	HasUsed  bool
	UsedUSD  float64
	HasLimit bool
	LimitUSD float64
}

// Share is one product's share of included usage.
type Share struct {
	Name         string
	SharePercent float64
}

// HistoryPoint is one optional used-percent sample.
type HistoryPoint struct {
	At          time.Time
	UsedPercent float64
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

// FormatUSD renders a non-negative dollar amount with at most two decimal
// places and a leading "$".
func FormatUSD(v float64) string {
	return "$" + formatMoney(v)
}

func formatMoney(v float64) string {
	rounded := math.Round(v*100) / 100
	s := strconv.FormatFloat(rounded, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" || s == "-0" {
		return "0"
	}
	return s
}

type rawSnapshot struct {
	SchemaVersion *int         `json:"schemaVersion"`
	FetchedAt     string       `json:"fetchedAt"`
	Plan          string       `json:"plan"`
	PeriodStart   string       `json:"periodStart"`
	Included      *rawIncluded `json:"included"`
	OnDemand      *rawOnDemand `json:"onDemand"`
	Breakdown     []rawShare   `json:"breakdown"`
	History       []rawHistory `json:"history"`
	Notes         string       `json:"notes"`
}

type rawIncluded struct {
	UsedPercent *json.Number `json:"usedPercent"`
	ResetAt     string       `json:"resetAt"`
	PeriodLabel string       `json:"periodLabel"`
}

type rawOnDemand struct {
	Status   string       `json:"status"`
	UsedUSD  *json.Number `json:"usedUsd"`
	LimitUSD *json.Number `json:"limitUsd"`
}

type rawShare struct {
	Name         string       `json:"name"`
	SharePercent *json.Number `json:"sharePercent"`
}

type rawHistory struct {
	At          string       `json:"at"`
	UsedPercent *json.Number `json:"usedPercent"`
}

func (r rawSnapshot) snapshot() (Snapshot, error) {
	if r.SchemaVersion == nil || *r.SchemaVersion != SchemaVersion {
		return Snapshot{}, ErrSchema
	}
	fetched, err := parseTime(r.FetchedAt, true)
	if err != nil {
		return Snapshot{}, err
	}
	if r.Included == nil {
		return Snapshot{}, ErrField
	}
	used, err := parsePercent(r.Included.UsedPercent, true)
	if err != nil {
		return Snapshot{}, ErrPercent
	}
	reset, hasReset, err := parseOptionalTime(r.Included.ResetAt)
	if err != nil {
		return Snapshot{}, err
	}
	period, err := optionalPeriod(r.Included.PeriodLabel)
	if err != nil {
		return Snapshot{}, err
	}
	start, hasStart, err := parseOptionalTime(r.PeriodStart)
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
	demand, err := parseOnDemand(r.OnDemand)
	if err != nil {
		return Snapshot{}, err
	}
	breakdown, err := parseBreakdown(r.Breakdown)
	if err != nil {
		return Snapshot{}, err
	}
	history, err := parseHistory(r.History)
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{
		SchemaVersion: SchemaVersion,
		FetchedAt:     fetched,
		Plan:          plan,
		Included: Included{
			UsedPercent: used,
			ResetAt:     reset,
			HasReset:    hasReset,
			PeriodLabel: period,
		},
		OnDemand:       demand,
		Breakdown:      breakdown,
		Notes:          notes,
		PeriodStart:    start,
		HasPeriodStart: hasStart,
		History:        history,
	}, nil
}

func (s Snapshot) encoded() encodedSnapshot {
	out := encodedSnapshot{
		SchemaVersion: SchemaVersion,
		FetchedAt:     s.FetchedAt.Format(time.RFC3339),
		Plan:          s.Plan,
		Included: encodedIncluded{
			UsedPercent: json.Number(FormatPercent(s.Included.UsedPercent)),
			PeriodLabel: s.Included.PeriodLabel,
		},
		Notes: s.Notes,
	}
	if s.Included.HasReset {
		out.Included.ResetAt = s.Included.ResetAt.Format(time.RFC3339)
	}
	if s.HasPeriodStart {
		out.PeriodStart = s.PeriodStart.Format(time.RFC3339)
	}
	if s.OnDemand.Present {
		item := encodedOnDemand{Status: s.OnDemand.Status}
		if s.OnDemand.HasUsed {
			n := json.Number(formatMoney(s.OnDemand.UsedUSD))
			item.UsedUSD = &n
		}
		if s.OnDemand.HasLimit {
			n := json.Number(formatMoney(s.OnDemand.LimitUSD))
			item.LimitUSD = &n
		}
		out.OnDemand = &item
	}
	if len(s.Breakdown) > 0 {
		out.Breakdown = make([]encodedShare, len(s.Breakdown))
		for i, share := range s.Breakdown {
			out.Breakdown[i] = encodedShare{
				Name:         share.Name,
				SharePercent: json.Number(FormatPercent(share.SharePercent)),
			}
		}
	}
	if len(s.History) > 0 {
		out.History = make([]encodedHistory, len(s.History))
		for i, point := range s.History {
			out.History[i] = encodedHistory{
				At:          point.At.Format(time.RFC3339),
				UsedPercent: json.Number(FormatPercent(point.UsedPercent)),
			}
		}
	}
	return out
}

type encodedSnapshot struct {
	SchemaVersion int              `json:"schemaVersion"`
	FetchedAt     string           `json:"fetchedAt"`
	Plan          string           `json:"plan,omitempty"`
	PeriodStart   string           `json:"periodStart,omitempty"`
	Included      encodedIncluded  `json:"included"`
	OnDemand      *encodedOnDemand `json:"onDemand,omitempty"`
	Breakdown     []encodedShare   `json:"breakdown,omitempty"`
	History       []encodedHistory `json:"history,omitempty"`
	Notes         string           `json:"notes,omitempty"`
}

type encodedIncluded struct {
	UsedPercent json.Number `json:"usedPercent"`
	ResetAt     string      `json:"resetAt,omitempty"`
	PeriodLabel string      `json:"periodLabel,omitempty"`
}

type encodedOnDemand struct {
	Status   string       `json:"status"`
	UsedUSD  *json.Number `json:"usedUsd,omitempty"`
	LimitUSD *json.Number `json:"limitUsd,omitempty"`
}

type encodedShare struct {
	Name         string      `json:"name"`
	SharePercent json.Number `json:"sharePercent"`
}

type encodedHistory struct {
	At          string      `json:"at"`
	UsedPercent json.Number `json:"usedPercent"`
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

func parseMoney(n *json.Number) (float64, bool, error) {
	if n == nil {
		return 0, false, nil
	}
	v, err := n.Float64()
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v < 0 || v > maxUSD {
		return 0, false, ErrField
	}
	return v, true, nil
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

func parseOnDemand(raw *rawOnDemand) (OnDemand, error) {
	if raw == nil {
		return OnDemand{}, nil
	}
	status := strings.TrimSpace(raw.Status)
	if !statusPattern.MatchString(status) {
		return OnDemand{}, ErrField
	}
	used, hasUsed, err := parseMoney(raw.UsedUSD)
	if err != nil {
		return OnDemand{}, err
	}
	limit, hasLimit, err := parseMoney(raw.LimitUSD)
	if err != nil {
		return OnDemand{}, err
	}
	return OnDemand{Present: true, Status: status, HasUsed: hasUsed, UsedUSD: used, HasLimit: hasLimit, LimitUSD: limit}, nil
}

func parseBreakdown(raw []rawShare) ([]Share, error) {
	if len(raw) > maxShares {
		return nil, ErrBreakdown
	}
	out := make([]Share, 0, len(raw))
	for _, item := range raw {
		name, err := optionalLabel(item.Name, 80)
		if err != nil || name == "" {
			return nil, ErrBreakdown
		}
		share, err := parsePercent(item.SharePercent, false)
		if err != nil {
			return nil, ErrBreakdown
		}
		out = append(out, Share{Name: name, SharePercent: share})
	}
	return out, nil
}

func parseHistory(raw []rawHistory) ([]HistoryPoint, error) {
	if len(raw) > maxHistory {
		return nil, ErrField
	}
	out := make([]HistoryPoint, 0, len(raw))
	for _, item := range raw {
		at, err := parseTime(item.At, true)
		if err != nil {
			return nil, err
		}
		used, err := parsePercent(item.UsedPercent, true)
		if err != nil {
			return nil, ErrPercent
		}
		out = append(out, HistoryPoint{At: at, UsedPercent: used})
	}
	return out, nil
}
