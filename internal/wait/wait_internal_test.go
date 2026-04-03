/*
   Golang rsync backup utility wrapper: szbck.
   Copyright (C) 2025-2026 Leslie Dancsecs

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

package wait

import (
	"testing"
	"time"

	"github.com/dancsecs/sztestlog"
)

func TestSnapshotProcess_NextHourIn(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	const minute = 22

	startTime := time.Date(2026, time.May, 15, 10, minute, 0, 0, time.Local)

	chk.Dur(
		nextHourAt(
			minute,
			startTime,
		).Sub(startTime),
		time.Hour,
	)

	chk.Dur(
		nextHourAt(
			minute,
			startTime.Add(time.Nanosecond),
		).Sub(startTime),
		time.Hour,
	)

	chk.Dur(
		nextHourAt(
			minute,
			startTime.Add(time.Hour+time.Nanosecond*999999999),
		).Sub(startTime),
		time.Hour*2,
	)
}
