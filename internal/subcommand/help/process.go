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

package help

import (
	"fmt"
	"strings"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/subcommand/create"
	"github.com/dancsecs/szbck/internal/subcommand/prune"
	"github.com/dancsecs/szbck/internal/subcommand/restore"
	"github.com/dancsecs/szbck/internal/subcommand/snapshot"
	"github.com/dancsecs/szbck/internal/subcommand/status"
	"github.com/dancsecs/szbck/internal/subcommand/trim"
	"github.com/dancsecs/szbck/internal/subcommand/vet"
)

// Process parses the remaining arguments deleting previous backups.
//
//nolint:cyclop  // Ok.
func Process(args *szargs.Args) (string, error) {
	var (
		subCommand string
		err        error
	)

	if args.HasNext() {
		subCommand = args.NextString("sub command", "")
	} else {
		subCommand = "all"
	}

	args.Done()
	err = args.Err()

	if err == nil {
		switch strings.ToLower(subCommand) {
		case "all":
			return "" +
				Usage + "\n" +
				HelpText + "\n" +
				create.HelpText + "\n" +
				snapshot.HelpText + "\n" +
				restore.HelpText + "\n" +
				prune.HelpText + "\n" +
				status.HelpText + "\n" +
				trim.HelpText + "\n" +
				vet.HelpText +
				"", nil
		case "h", "help":
			return HelpText, nil
		case "i", "init", "initialize":
			return create.HelpText, nil
		case "s", "snap", "snapshot":
			return snapshot.HelpText, nil
		case "r", "res", "restore":
			return restore.HelpText, nil
		case "p", "prune":
			return prune.HelpText, nil
		case "stat", "status":
			return status.HelpText, nil
		case "t", "trim":
			return trim.HelpText, nil
		case "v", "vet":
			return vet.HelpText, nil
		default:
			err = fmt.Errorf(
				"%w: '%s'",
				ErrUnknownSubcommand,
				subCommand,
			)
		}
	}

	return "", fmt.Errorf("%w: %w", ErrHelpError, err)
}
