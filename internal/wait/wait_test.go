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

func TestSnapshotProcess_WaitTillNow(t *testing.T) {
	chk := sztestlog.CaptureStdout(t)
	defer chk.Release()

	wait.Until(false, time.Now())

	chk.Stdout()
}

func TestSnapshotProcess_WaitTill(t *testing.T) {
	chk := sztestlog.CaptureStdout(t)
	defer chk.Release()

	wait.Until(false, time.Now().Add(time.Millisecond*500))

	chk.AddSub(`\-?\d[\d\,\.]*(?:s|ms|µs|ns)?`, "#")
	chk.Stdout(
		"Next Backup at ### #:#:# in: #",
		"Restarting at: ### #:#:# TargetDelta: #",
	)
}

func TestSnapshotProcess_WaitTillMonitor(t *testing.T) {
	chk := sztestlog.CaptureStdout(t)
	defer chk.Release()

	wait.Until(true, time.Now().Add(time.Second*11))

	chk.AddSub(`\-?\d[\d\,\.]*(?:s|ms|µs|ns)?`, "#")
	chk.Stdout(
		"" +
			"Next Backup at ### #:#:# in: #\r" +
			"Next Backup at ### #:#:# in: #\r" +
			"Next Backup at ### #:#:# in: #\r" +
			"Next Backup at ### #:#:# in: #\r" +
			"Next Backup at ### #:#:# in: #\r" +
			"Next Backup at ### #:#:# in: #\r" +
			"Next Backup at ### #:#:# in: #\r" +
			"Restarting at: ### #:#:# TargetDelta: #                         ",
	)
}
