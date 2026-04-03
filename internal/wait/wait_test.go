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

package wait_test

import (
	"testing"
	"time"

	"github.com/dancsecs/szbck/internal/wait"
	"github.com/dancsecs/sztestlog"
)

const clearLine = "                    "

func TestSnapshotProcess_NextHourIn(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	const minute = 22

	startTime := time.Date(2026, time.May, 15, 10, minute, 0, 0, time.Local)

	chk.Dur(
		wait.NextHourAt(
			minute,
			startTime,
		).Sub(startTime),
		time.Hour,
	)

	chk.Dur(
		wait.NextHourAt(
			minute,
			startTime.Add(time.Nanosecond),
		).Sub(startTime),
		time.Hour,
	)

	chk.Dur(
		wait.NextHourAt(
			minute,
			startTime.Add(time.Hour+time.Nanosecond*999999999),
		).Sub(startTime),
		time.Hour*2,
	)
}

func TestSnapshotProcess_WaitTillNow(t *testing.T) {
	chk := sztestlog.CaptureStdout(t)
	defer chk.Release()

	wait.Until("Timer Title Now", false, time.Now())

	chk.Stdout()
}

func TestSnapshotProcess_WaitTill(t *testing.T) {
	chk := sztestlog.CaptureStdout(t)
	defer chk.Release()

	wait.Until(
		"Timer Title (500ms)",
		false,
		time.Now().Add(time.Millisecond*500),
	)

	chk.AddSub(`\-?\d[\d\,\.]*(?:s|ms|µs|ns)?`, "#")
	chk.Stdout(
		"Starting 'Timer Title (500ms)' at ### #:#:# in: #"+
			clearLine,
		"Restarted 'Timer Title (500ms)' at: ### #:#:# TargetDelta: #"+
			clearLine,
	)
}

func TestSnapshotProcess_WaitTillMonitor(t *testing.T) {
	chk := sztestlog.CaptureStdout(t)
	defer chk.Release()

	wait.Until("Timer Title", true, time.Now().Add(time.Second*11))

	chk.AddSub(`\-?\d[\d\,\.]*(?:s|ms|µs|ns)?`, "#")
	chk.Stdout(
		"" +
			"Starting 'Timer Title' at ### #:#:# in: #" + clearLine + "\r" +
			"Starting 'Timer Title' at ### #:#:# in: #" + clearLine + "\r" +
			"Starting 'Timer Title' at ### #:#:# in: #" + clearLine + "\r" +
			"Starting 'Timer Title' at ### #:#:# in: #" + clearLine + "\r" +
			"Starting 'Timer Title' at ### #:#:# in: #" + clearLine + "\r" +
			"Starting 'Timer Title' at ### #:#:# in: #" + clearLine + "\r" +
			"Starting 'Timer Title' at ### #:#:# in: #" + clearLine + "\r" +
			"Restarted 'Timer Title' at: ### #:#:# TargetDelta: #" +
			clearLine,
	)
}
