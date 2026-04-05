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

const clearLine = "                    "

// NextHourAt returns the time adjusted to the atMin time after the hour.
func NextHourAt(atMin int, from time.Time) time.Time {
	next := time.Date(
		from.Year(), from.Month(), from.Day(),
		from.Hour(), atMin, 0, 0,
		from.Location(),
	)

	for !next.After(from) {
		next = next.Add(time.Hour)
	}

	return next
}

// NextMinuteIn returns the duration until the start of the next minute.
func nextMinuteIn(from time.Time) time.Duration {
	next := time.Date(
		from.Year(), from.Month(), from.Day(),
		from.Hour(), from.Minute(), 0, 0,
		from.Location(),
	)

	for !next.After(from) {
		next = next.Add(time.Minute)
	}

	return next.Sub(from)
}

// NextSecondIn returns the duration until the start of the next second.
func nextSecondIn(from time.Time) time.Duration {
	next := time.Date(
		from.Year(), from.Month(), from.Day(),
		from.Hour(), from.Minute(), from.Second(), 0,
		from.Location(),
	)

	for !next.After(from) {
		next = next.Add(time.Second)
	}

	return next.Sub(from)
}

func chkIn(from time.Time, maxDur time.Duration) time.Duration {
	switch {
	case maxDur <= time.Nanosecond:
		return time.Nanosecond
	case maxDur <= time.Second:
		return maxDur
	case maxDur <= time.Minute+time.Minute:
		return nextSecondIn(from)
	default:
		return nextMinuteIn(from)
	}
}

// Until waits (sleeps) until the specified time displaying an updated
// countdown if monitor is true.
func Until(title string, monitor bool, targetTime time.Time) {
	targetTimeStr := targetTime.Format("2006-01-02 15:04:05.999")

	now := time.Now()
	maxSleep := targetTime.Sub(now)

	if maxSleep <= 0 {
		return
	}

	szlog.Say0f(
		"Starting '%s' at %s in: %v%s\r",
		title,
		targetTimeStr,
		maxSleep.Truncate(time.Second),
		clearLine,
	)

	for maxSleep > 0 {
		if monitor {
			szlog.Say0f(
				"Starting '%s' at %s in: %v%s\r",
				title,
				targetTimeStr,
				maxSleep.Truncate(time.Second),
				clearLine,
			)
		}

		time.Sleep(chkIn(now, maxSleep))
		now = time.Now()
		maxSleep = targetTime.Sub(now)
	}

	now = time.Now()
	if monitor {
		szlog.Say0f(
			"\nRestarted '%s' at: %s TargetDelta: %v%s\n",
			title,
			now.Format("2006-01-02 15:04:05.999"),
			now.Sub(targetTime),
			clearLine,
		)
	} else {
		szlog.Say1f(
			"Restarted '%s' at: %s TargetDelta: %v%s\n",
			title,
			now.Format("2006-01-02 15:04:05.999"),
			now.Sub(targetTime),
			clearLine,
		)
	}
}
