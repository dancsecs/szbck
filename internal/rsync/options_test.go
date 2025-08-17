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

package rsync_test

import (
	"testing"

	"github.com/dancsecs/szbck/internal/rsync"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestRsync_BuildArgs_NoneNoVerbose(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelError)
	defer chk.Release()

	chk.StrSlice(
		rsync.BuildArgs(
			false,          // deleteFromTarget
			false,          // dryRun
			"",             // latest
			[]string{"-a"}, // options
			nil,            // additionalOptions
			"from",         // fromPath
			"to",           // toPath
		),
		[]string{
			"-a",
			"from",
			"to",
		},
	)
}

func TestRsync_BuildArgs_NoneOneVerbose(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelInfo)
	defer chk.Release()

	chk.StrSlice(
		rsync.BuildArgs(
			false,          // deleteFromTarget
			false,          // dryRun
			"",             // latest
			[]string{"-a"}, // options
			nil,            // additionalOptions
			"from",         // fromPath
			"to",           // toPath
		),
		[]string{
			"--verbose",
			"-a",
			"from",
			"to",
		},
	)
}

func TestRsync_BuildArgs_NoneTwoVerbose(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelDebug)
	defer chk.Release()

	chk.StrSlice(
		rsync.BuildArgs(
			false,          // deleteFromTarget
			false,          // dryRun
			"",             // latest
			[]string{"-a"}, // options
			nil,            // additionalOptions
			"from",         // fromPath
			"to",           // toPath
		),
		[]string{
			"--verbose",
			"--verbose",
			"-a",
			"from",
			"to",
		},
	)
}

func TestRsync_BuildArgs_NoneQuiet(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	chk.StrSlice(
		rsync.BuildArgs(
			false, // deleteFromTarget
			false, // dryRun
			"",    // latest
			[]string{
				"-a",
				"--quiet",
			}, // options
			nil,    // additionalOptions
			"from", // fromPath
			"to",   // toPath
		),
		[]string{
			"-a",
			"--quiet",
			"from",
			"to",
		},
	)
}

func TestRsync_BuildArgs_WithDeleteFromTarget(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	chk.StrSlice(
		rsync.BuildArgs(
			true,           // deleteFromTarget
			false,          // dryRun
			"",             // latest
			[]string{"-a"}, // options
			nil,            // additionalOptions
			"from",         // fromPath
			"to",           // toPath
		),
		[]string{
			"--verbose",
			"--verbose",
			"-a",
			rsync.FlgDelete,
			"from",
			"to",
		},
	)
}

func TestRsync_BuildArgs_WithDryRun(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	chk.StrSlice(
		rsync.BuildArgs(
			false,          // deleteFromTarget
			true,           // dryRun
			"",             // latest
			[]string{"-a"}, // options
			nil,            // additionalOptions
			"from",         // fromPath
			"to",           // toPath
		),
		[]string{
			"--verbose",
			"--verbose",
			"-a",
			rsync.FlgDryRun,
			"from",
			"to",
		},
	)
}

func TestRsync_BuildArgs_WithLinkDest(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	chk.StrSlice(
		rsync.BuildArgs(
			false,          // deleteFromTarget
			false,          // dryRun
			"linkDest",     // latest
			[]string{"-a"}, // options
			nil,            // additionalOptions
			"from",         // fromPath
			"to",           // toPath
		),
		[]string{
			"--verbose",
			"--verbose",
			"-a",
			rsync.FlgLinkDest + "linkDest",
			"from",
			"to",
		},
	)
}

func TestRsync_BuildArgs_WithAdditional(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	chk.StrSlice(
		rsync.BuildArgs(
			false,          // deleteFromTarget
			false,          // dryRun
			"",             // latest
			[]string{"-a"}, // options
			[]string{"-i"}, // additionalOptions
			"from",         // fromPath
			"to",           // toPath
		),
		[]string{
			"--verbose",
			"--verbose",
			"-a",
			"-i",
			"from",
			"to",
		},
	)
}

func TestRsync_BuildArgs_WithOnlyAdditional(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	chk.StrSlice(
		rsync.BuildArgs(
			false,          // deleteFromTarget
			false,          // dryRun
			"",             // latest
			nil,            // options
			[]string{"-i"}, // additionalOptions
			"from",         // fromPath
			"to",           // toPath
		),
		[]string{
			"--verbose",
			"--verbose",
			"-i",
			"from",
			"to",
		},
	)
}
