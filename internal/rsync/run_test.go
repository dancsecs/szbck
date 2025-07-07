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
	"os"
	"testing"

	"github.com/dancsecs/szbck/internal/rsync"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestRsyncRun_NoArgs(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	err := rsync.Run(nil, os.Stdout, nil)

	chk.Err(
		err,
		rsync.ErrRsyncError.Error()+
			": exit status 1"+
			"",
	)

	chk.AddSub(
		`Running\scommand\:\s.*rsync\s`,
		"Running command: RsyncCommand",
	)
	chk.Log()
	chk.Stdout(
		"Running command: RsyncCommand",
	)
	chk.Stderr()
}

func TestRsyncRun_SimpleFiles(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	source := chk.CreateTmpSubDir("source")
	target := chk.CreateTmpSubDir("target")

	_ = chk.CreateTmpFileIn(source, []byte("file1"))
	_ = chk.CreateTmpFileIn(source, []byte("file2"))

	err := rsync.Run([]string{"-av", source, target}, os.Stdout, os.Stderr)

	chk.NoErr(err)

	chk.AddSub(
		`Running\scommand\:\s.*rsync`,
		"Running command: RsyncCommand",
	)
	chk.AddSub(`\d+`, "#")
	chk.Log()
	chk.Stdout(
		"Running command: "+
			"RsyncCommand -av "+source+" "+target,
		"sending incremental file list",
		"source/",
		"source/tmpFile0.tmp",
		"source/tmpFile1.tmp",
		"",
		"sent # bytes  received # bytes  #.# bytes/sec",
		"total size is #  speedup is #.#",
	)
	chk.Stderr()
}
