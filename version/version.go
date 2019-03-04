package version

// Default build-time variables
// These values are overridden via ldflags
var (
	PlatformName = ""
	Version      = "unknown-version"
	GitCommit    = "unknown-commit"
	BuildTime    = "unknown-buildtime"
)
