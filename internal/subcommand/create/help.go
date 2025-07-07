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

package create

// HelpText describes the utility subcommand.
const HelpText = `{c | create} ` +
	`[-o filePath] [-p perm] [-e exclude ...] source

Create generates a new Szerszam backup configuration file.

	[-o config.sbc]
		Writes the configuration to the named Szerszam backup Configuration
      file.  Otherwise it is written to stdout.

	[-t target]
		Specifies the default target directory to create snapshots in.  If it
		not provided then a target argument will be mandatory for snapshot,
		restore, prune and status subcommands.

	source
		Specifies the root directory to back up.
`
