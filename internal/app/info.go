package app

import "runtime"

var (
	// Build-time variables (set via ldflags)
	Version   = "dev"
	GitCommit = "unknown"
	GitRef    = "unknown"
	GitTag    = "unknown"
	BuildDate = "unknown"
)

// CompilerVersion returns the Go compiler version
func CompilerVersion() string {
	return runtime.Version()
}

// Info represents build information
type Info struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	GitCommit       string `json:"git_commit"`
	GitRef          string `json:"git_ref"`
	GitTag          string `json:"git_tag"`
	BuildDate       string `json:"build_date"`
	CompilerVersion string `json:"compiler_version"`
}

// GetInfo returns build information
func GetInfo() Info {
	return Info{
		Name:            Name,
		Version:         Version,
		GitCommit:       GitCommit,
		GitRef:          GitRef,
		GitTag:          GitTag,
		BuildDate:       BuildDate,
		CompilerVersion: CompilerVersion(),
	}
}
