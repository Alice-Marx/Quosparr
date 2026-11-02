// Package cursorusage reads Cursor product usage from a local JSON snapshot.
//
// Cursor IDE, Cloud Agents, and the rest of the Cursor product have no public
// third-party billing API. This package does not call the network, scrape
// cookies, read dashboard HTML, or read the Cursor console. The only input is
// the snapshot file under the Quosparr data directory (see UsagePath), written
// by a manual import or by a later routine. This package does not write that
// file except through Import, which copies a document the user chose.
//
// The snapshot is non-secret usage stats. Tokens, API keys, and other unknown
// fields are ignored on read and are not written back. Grok Bot assistant
// usage is a different file and a different package.
//
// includedPools, when set, splits the included allowance into Cursor Models
// and Other Models. Each pool object that is present must carry usedPercent
// from 0 to 100; label and note are optional. When includedPools is absent,
// callers keep included.usedPercent and the optional breakdown. grokBotWeeklyRef
// only records that the weekly assistant quota lives in the Grok Bot snapshot.
// This file does not store that weekly percent.
package cursorusage
