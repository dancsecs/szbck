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

package snapshot

import (
	"errors"
)

// Snapshot errors.
var (
	ErrSnapshotError = errors.New("snapshot error")
	ErrAtRange       = errors.New(
		"daemon at range must be between 0 and 59 inclusively",
	)
	ErrAtUsage      = errors.New("--at specified without --daemon")
	ErrMonitorUsage = errors.New("--monitor specified without --daemon")

	ErrTrimNotImplement = errors.New("trim retention not yet implemented")
)
