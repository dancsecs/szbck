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

package settings

import (
	"os"
	"time"

	"github.com/dancsecs/szbck/internal/target"
)

// Config defines required parameter to run a szerszam backup.
type Config struct {
	// Source defines the directory to be backed up.
	Source string
	// Target defines the directory to store backup snapshots.
	Target *target.Path
	// Default permissions for new backup directory.
	Permission os.FileMode
	// Options for both snapshots and restores.
	Options []string
	// Addition snapshot options.
	SnapshotOptions []string
	// Addition restore options.
	RestoreOptions []string
	// Retention parameters
	KeepHourly time.Duration
	KeepDaily  time.Duration
}
