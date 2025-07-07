/*
   Golang rsync backup utility wrapper: szbck.
   Copyright (C) 2025 Leslie Dancsecs

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package rsync

import (
	"strings"

	"github.com/dancsecs/szlog"
)

// Rsync flags.
const (
	FlgDryRun    = "--dry-run"
	FlgDelete    = "--delete"
	FlgLinkDest  = "--link-dest="
	extraOptions = 5 // Flags plus source and destination.
)

func outConfigured(options []string) bool {
	for _, option := range options {
		hasOutput := option == "-v" || option == "--verbose" ||
			option == "-q" || option == "--quiet" ||
			strings.HasPrefix(option, "--info=") ||
			strings.HasPrefix(option, "--debug=") ||
			strings.HasPrefix(option, "--stderr=")
		if hasOutput {
			return true
		}
	}

	return false
}

func addOption(options []string, add bool, opt string) []string {
	if add {
		options = append(options, opt)
	}

	return options
}

// BuildArgs options returns a slice representing thr rsync arguments.
func BuildArgs(
	deleteFromTarget bool,
	dryRun bool,
	linkDest string,
	basicOptions []string,
	additionalOptions []string,
	fromPath string,
	toPath string,
) []string {
	var verboseOptions []string

	if !outConfigured(basicOptions) && !outConfigured(additionalOptions) {
		// Rsync verbose not configured.  Set according to application
		// verbose level.
		if szlog.Level() >= szlog.LevelInfo {
			verboseOptions = []string{"--verbose", "--verbose"}
		} else if szlog.Level() >= szlog.LevelWarn {
			verboseOptions = []string{"--verbose"}
		}
	}

	options := make(
		[]string, 0, 0+
			len(verboseOptions)+
			len(basicOptions)+
			len(additionalOptions)+
			extraOptions,
	)

	options = append(options, verboseOptions...)
	options = append(options, basicOptions...)

	options = addOption(options, deleteFromTarget, FlgDelete)
	options = addOption(options, dryRun, FlgDryRun)
	options = addOption(options, linkDest != "", FlgLinkDest+linkDest)

	options = append(options, additionalOptions...)

	options = append(options, fromPath, toPath)

	return options
}
