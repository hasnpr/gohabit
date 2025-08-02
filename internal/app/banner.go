// Package app is a global Application object.
//
// It sets up the Application and provides access to the Application.
package app

import (
	"fmt"
)

var (
	// GitCommit is SHA1 ref of current build
	GitCommit string

	// GitRef is branch name of current build
	GitRef string

	// GitTag is branch name of current build
	GitTag string

	// BuildDate is the timestamp of build
	BuildDate string

	// CompilerVersion is the version of go compiler
	CompilerVersion string
)

// Name is the name of the Application
const (
	Name    = "GoHabit"
	Version = 0
)

const asciiArt = `  _     ____       _      ____     ___
 (_)   |_| \_\ /_/   \_\ |____/   \___/   \___/  |_____|
`

// Banner returns the ascii art banner
func Banner(noBanner bool) string {
	if !noBanner {
		return asciiArt + "\n" + fmt.Sprintf("Tag: %s, Ref: %s, GitCommit: %s, BuildDate: %s, Compiler: %s\n", GitTag, GitRef, GitCommit, BuildDate, CompilerVersion)
	}

	return ""
}
