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
package cursorusage
