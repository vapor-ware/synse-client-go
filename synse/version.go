package synse

// version.go defines version information of the client.

import (
	"runtime"
)

var (
	// BuildDate is the timestamp for when the binary was built.
	BuildDate string

	// GitCommit is the commit hash at which the binary was built.
	GitCommit string

	// GitTag is the git tag name at which the binary was built.
	GitTag string

	// GoVersion is the version of Go used to build the binary.
	GoVersion string

	// BinVersion is the canonical verison for the binary.
	BinVersion string
)

// Info describe the version info of the client library.
type Info struct {
	Arch       string
	BuildDate  string
	GitCommit  string
	GitTag     string
	GoVersion  string
	OS         string
	BinVersion string
}

// Get gets the version information for the client binary. The package
// variables should be set using build-time arguments.
func Get() *Info {
	return &Info{
		Arch:       runtime.GOARCH,
		OS:         runtime.GOOS,
		BuildDate:  BuildDate,
		GitCommit:  GitCommit,
		GitTag:     GitTag,
		GoVersion:  GoVersion,
		BinVersion: BinVersion,
	}
}
