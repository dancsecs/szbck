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

package trim

// HelpText describes the overall operation of the utility.
const HelpText = `{t | trim} ` +
	`[--dry-run] [-t target] config.szb

Implements the specified retention policy as defined in the backup
configuration file deleting backups as appropriate. The most recent snapshot
pointed to by the "latest" symbolic link is never deleted.

   [--dry-run]
      Identifies all of the actions the utility would take without making any
      changes to the backup source.

   [-t target]
      Specifies the backup set to prune.  It is optional if the backup config
      file specifies a target and mandatory if not specified in the backup
      config file.

   config.sbc
      the backup configuration file defining the backup.
`
