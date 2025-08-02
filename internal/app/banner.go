// Package app is a global Application object.
//
// It sets up the Application and provides access to the Application.
package app

import (
	"fmt"
)

const asciiArt = `
   ██████╗  ██████╗ ██╗  ██╗ █████╗ ██████╗ ██╗████████╗
  ██╔════╝ ██╔═══██╗██║  ██║██╔══██╗██╔══██╗██║╚══██╔══╝
  ██║  ███╗██║   ██║███████║███████║██████╔╝██║   ██║
  ██║   ██║██║   ██║██╔══██║██╔══██║██╔══██╗██║   ██║
  ╚██████╔╝╚██████╔╝██║  ██║██║  ██║██████╔╝██║   ██║
   ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ ╚═╝   ╚═╝
  `

// Banner returns the ascii art banner
func Banner(showBanner bool) string {
	if showBanner {
		info := GetInfo()
		return asciiArt + "\n" + fmt.Sprintf("Version: %s, Tag: %s, Ref: %s, Commit: %s, BuildDate: %s, Compiler: %s\n",
			info.Version, info.GitTag, info.GitRef, info.GitCommit, info.BuildDate, info.CompilerVersion)
	}

	return ""
}
