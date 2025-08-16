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

/*
Szbck wraps the system's rsync utility to create time machine like backups. It
is driven by a backup configuration file which identifies the root directory
to backup, a target directory to store snapshots, as well as enabling and
disabling various rsync options, defines the default permissions for the root
of each backup (snapshot) and lists items to exclude from backup and restores.

It can be used to backup a single machine and/or to keep several machines in
sync.  The utility defines several subcommands implementing the various backup
functions as follows:

	SubCommand  Description
	==========  =============================================================
	Help        Displays help on the utility and all of the subcommands.

	Create      Creates a backup configuration file.

	Snapshot    Creates a new backup snapshot hard linking unchanged items to
				previous snapshots.  It will implement a retention policy if
				and only if the --trim option is provided.

	Restore     Restores/removes/replaces files in the source directory
				identified by the backup configuration file honoring all
				exclusions identified by the config file from a backup
				snapshot. The backup snapshot need not be from the same machine
				enabling syncing between machines.

	Prune       Removes oldest backup snapshots.  NOTE: This does not
				implement a retention policy but simply purges the
				specified number of the oldest snapshots.  A retention
				policy has not yet been implemented.

	Trim		Manually implements the retention policy as specified in the
				configuration file.  Can be invoked by using the --trim option
				on the snapshot subcommand.

	Status		Reports on the number of backup snapshots and the space used
				overall and by each snapshot.

	Vet         Parses a backup configuration file identifying any errors
				or problems without making any attempts at any operations.

Examples:

	// Display help on the utility and all sub commands.
		szbck help

	// Creates a new backup configuration file.
		szbck create -o config.szb -t backupTo /home/myDirectory

	// Create a new snapshot.
		szbck snapshot config.szb

	// Restore updated/missing files and purge extra files (unless the --keep
	// option is specified).
		szbck restore config.szb

	// Purge the oldest 5 snapshots from a backup set.
		szbck prune -n 5 config.szb

	// Vet changes made to a config.szb file.
		szbck vet config.szb
*/
package main

import (
	"os"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal"
)

/*
Simply invokes the internal version of main which returns an int as a classic
type main function.  This is returned to the operating system via the os.Exit
function which cannot be tested.  Therefore this wrapper is the only function
in this utility that is not tested.
*/
func main() {
	args := szargs.New("", os.Args)

	returnValue := internal.Main(args)

	os.Exit(returnValue)
}
