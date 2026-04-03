/*
   Golang rsync backup utility wrapper: szbck.
   Copyright (C) 2026 Leslie Dancsecs

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
	"time"

	"github.com/dancsecs/szlog"
)

// NextHourAt returns the time adjusted to the atMin time after the hour.
func NextHourAt(atMin int, from time.Time) time.Time {
	nextStartAt := time.Date(
		from.Year(), from.Month(), from.Day(),
		from.Hour(), atMin, 0, 0,
		from.Location(),
	)

	for !nextStartAt.After(from) {
		nextStartAt = nextStartAt.Add(time.Hour)
	}

	return nextStartAt
}

// Until waits (sleeps) until the specified time displaying an updated
// countdown if monitor is true.
func Until(monitor bool, targetTime time.Time) {
	targetTimeStr := targetTime.Format("2006-01-02 15:04:05.999")

	now := time.Now()
	maxSleep := targetTime.Sub(now)

	if maxSleep <= 0 {
		return
	}

	for maxSleep > 0 {
		if !monitor {
			szlog.Say1f("Next Backup at %s in: %v\n", targetTimeStr, maxSleep)
			time.Sleep(maxSleep)
		} else {
			szlog.Say0f(
				"Next Backup at %s in: %v\r", targetTimeStr, maxSleep,
			)

			switch {
			case maxSleep < time.Second:
				time.Sleep(maxSleep)
			case maxSleep < time.Second*10:
				time.Sleep(time.Second)
			default:
				time.Sleep(time.Second * 5) //nolint:mnd // 5 second default.
			}
		}

		now = time.Now()
		maxSleep = targetTime.Sub(now)
	}

	now = time.Now()
	if monitor {
		szlog.Say0f(
			"Restarting at: %s TargetDelta: %v                         \n",
			now.Format("2006-01-02 15:04:05.999"),
			now.Sub(targetTime),
		)
	} else {
		szlog.Say1f(
			"Restarting at: %s TargetDelta: %v\n",
			now.Format("2006-01-02 15:04:05.999"),
			now.Sub(targetTime),
		)
	}
}
