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

package internal

import (
	"fmt"
	"strings"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/subcommand/create"
	"github.com/dancsecs/szbck/internal/subcommand/help"
	"github.com/dancsecs/szbck/internal/subcommand/prune"
	"github.com/dancsecs/szbck/internal/subcommand/restore"
	"github.com/dancsecs/szbck/internal/subcommand/snapshot"
	"github.com/dancsecs/szbck/internal/subcommand/status"
	"github.com/dancsecs/szbck/internal/subcommand/trim"
	"github.com/dancsecs/szbck/internal/subcommand/vet"
	"github.com/dancsecs/szlog"
)

func parseGlobalArgs(args []string) (string, []string, error) {
	var (
		verbose       int
		doubleVerbose int
		subCommand    string
		err           error
	)

	verbose, args = szargs.Arg("-v").Count(args)

	doubleVerbose, args = szargs.Arg("-vv").Count(args)

	verbose += doubleVerbose + doubleVerbose

	for verbose > 0 {
		verbose--

		szlog.IncLevel()
	}

	subCommand, args, err = szargs.Next("sub command", args)

	if err == nil {
		return subCommand, args, nil
	}

	return "", nil, err //nolint:wrapcheck // Ok.
}

// Main is the actual mainline for the puzzle application classically returning
// an int to be returned when exiting.
//
//nolint:cyclop   // Ok.
func Main(programName string, args []string) int {
	var (
		subCommand  string
		outText     string
		err         error
		returnValue int
	)

	subCommand, args, err = parseGlobalArgs(args)
	if err == nil {
		switch strings.ToLower(subCommand) {
		case "h", "help":
			outText, err = help.Process(args)
			if err == nil {
				outText = programName + "\n" + outText
			}
		case "c", "create":
			outText, err = create.Process(args)
		case "s", "snap", "snapshot":
			outText, err = snapshot.Process(args)
		case "r", "res", "restore":
			outText, err = restore.Process(args)
		case "p", "prune":
			outText, err = prune.Process(args)
		case "stat", "status":
			outText, err = status.Process(args)
		case "t", "trim":
			outText, err = trim.Process(args)
		case "v", "vet":
			outText, err = vet.Process(args)
		default:
			err = fmt.Errorf(
				"%w: '%s'",
				ErrUnknownSubcommand,
				subCommand,
			)
		}
	}

	if err != nil {
		szlog.Fatalf("%s - %v\n", programName, err)

		returnValue = 1
	}

	fmt.Print(outText) //nolint:forbidigo // Ok.

	return returnValue
}
