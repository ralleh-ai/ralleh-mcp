package runtime

// Build metadata is populated by release builds through -ldflags. Defaults are
// intentionally safe and deterministic for local tests.
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)
