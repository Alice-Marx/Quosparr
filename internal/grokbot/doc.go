// Package grokbot reads the Grok Bot assistant weekly quota from a local
// JSON snapshot.
//
// Grok Bot assistant has no public third-party billing API. This package
// does not call the network, scrape cookies, read Grok Bot settings pages,
// or read Cursor dashboards. The only input is the snapshot file under the
// Quosparr data directory (see UsagePath), written by a manual import or by
// a later routine. The snapshot is non-secret usage stats: tokens, API keys
// and other unknown fields are ignored on read and are not written back.
package grokbot
