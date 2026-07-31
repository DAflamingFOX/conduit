package version

import (
	"fmt"
	"runtime/debug"
)

// Version is the application version string.
// It can be overridden at build time using ldflags:
//
//	go build -ldflags "-X github.com/daflamingfox/conduit/internal/version.Version=v1.0.0"
var Version = "dev"

// Commit is the git commit hash string (optional ldflags override).
var Commit = ""

// BuildTime is the ISO build timestamp string (optional ldflags override).
var BuildTime = ""

// Get returns a clean, deterministic version string.
// If Version was set via ldflags (e.g. git tag), it returns Version.
// Otherwise, it inspects runtime debug build info embedded by the Go compiler (git commit hash, modified status).
func Get() string {
	if Version != "" && Version != "dev" {
		return Version
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return Version
	}

	var revision string
	var modified bool

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}

	if revision != "" {
		// Truncate commit hash to 7 characters if long
		shortCommit := revision
		if len(shortCommit) > 7 {
			shortCommit = shortCommit[:7]
		}
		if modified {
			return fmt.Sprintf("dev-%s-dirty", shortCommit)
		}
		return fmt.Sprintf("dev-%s", shortCommit)
	}

	return Version
}
