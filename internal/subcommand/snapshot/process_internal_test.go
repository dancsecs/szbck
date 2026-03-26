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

package snapshot

import (
	"testing"
	"time"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/sztestlog"
)

func TestSnapshotProcess_NextHourIn(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	startTime := time.Date(2026, time.May, 15, 10, 22, 0, 0, time.Local)

	chk.Dur(
		nextHourIn(
			22,
			startTime,
		),
		time.Hour,
	)

	chk.Dur(
		nextHourIn(
			22,
			startTime.Add(-time.Nanosecond),
		),
		time.Nanosecond,
	)

	chk.Dur(
		nextHourIn(
			22,
			startTime.Add(time.Nanosecond),
		),
		time.Hour-time.Nanosecond,
	)

	chk.Dur(
		nextHourIn(
			22,
			startTime.Add(time.Hour+time.Nanosecond*999999999),
		),
		time.Minute*59+time.Second*59+time.Nanosecond,
	)
}

//nolint:dogsled,funlen // Ok.
func TestSnapshotProcess_ParseArgDaemonAt(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	startTime := time.Date(2026, time.May, 15, 10, 22, 0, 0, time.Local)

	_, _, _, daemon, runAtMin, err := parseArgs(
		szargs.New(
			"programDesc",
			[]string{
				"programName", "--at", "MISSING_CONFIG_FILE",
			}),
		startTime,
	)

	chk.False(daemon)
	chk.Int(runAtMin, 0)
	chk.Err(
		err,
		chk.ErrChain(
			szargs.ErrUnexpected,
			"[MISSING_CONFIG_FILE]",
		),
	)

	_, _, _, daemon, runAtMin, err = parseArgs(
		szargs.New(
			"programDesc",
			[]string{
				"programName", "--daemon", "MISSING_CONFIG_FILE",
			}),
		startTime,
	)

	chk.True(daemon)
	chk.Int(runAtMin, 22)
	chk.Err(
		err,
		chk.ErrChain(
			settings.ErrLoad,
			"open MISSING_CONFIG_FILE",
			"no such file or directory",
		),
	)

	_, _, _, daemon, runAtMin, err = parseArgs(
		szargs.New(
			"programDesc",
			[]string{
				"programName",
				"--daemon",
				"--at", "200",
				"MISSING_CONFIG_FILE",
			}),
		startTime,
	)

	chk.True(daemon)
	chk.Int(runAtMin, 0)
	chk.Err(
		err,
		chk.ErrChain(
			ErrAtRange,
		),
	)

	_, _, _, daemon, runAtMin, err = parseArgs(
		szargs.New(
			"programDesc",
			[]string{
				"programName",
				"--daemon",
				"--at", "55",
				"MISSING_CONFIG_FILE",
			}),
		startTime,
	)

	chk.True(daemon)
	chk.Int(runAtMin, 55)
	chk.Err(
		err,
		chk.ErrChain(
			settings.ErrLoad,
			"open MISSING_CONFIG_FILE",
			"no such file or directory",
		),
	)
}
