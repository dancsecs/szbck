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

// Usage describes the overall operation of the utility.
const Usage = "" +
	`Szerszam backup utility takes Apple Time machine like snapshots.  It
requires the underlying system to have the utility rsync installed which
will perform the actual snapshots.  One of the following sub commands must
be provided as follows:
`

// HelpText describes the overall operation of the help subcommand.
const HelpText = `{h | help} ` +
	`[subcommand]

Help information is displayed.

   [subcommand]
      Limits the help displayed to the specified subcommand.
`
