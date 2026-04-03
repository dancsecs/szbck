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

package purge_test

import (
	"os/exec"
	"testing"
	"time"

	"github.com/dancsecs/szbck/internal/purge"
	"github.com/dancsecs/szbck/internal/target"
	"github.com/dancsecs/sztestlog"
)

const (
	permNoWrite = 0o0500
)

func TestInternalTrim_ProcessPurge_InvalidDir(t *testing.T) {
	chk := sztestlog.CaptureLog(t)
	defer chk.Release()

	findPath, err := exec.LookPath("find")
	chk.NoErr(err)

	err = purge.Directory("DOES_NOT_EXIST")
	chk.Err(
		err,
		chk.ErrChain(
			purge.ErrDirRights,
			findPath,
			"‘DOES_NOT_EXIST’",
			"No such file or directory",
		),
	)

	chk.Log()
}

func TestInternalTrim_ProcessPurge_FailureAfterSuccess(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t)
	defer chk.Release()

	// A little before now without sleeping.
	startTime := time.Now().Add(-time.Millisecond)

	rootDir := chk.CreateTmpSubDir("root")
	trg, err := target.New(rootDir)
	chk.NoErr(err)

	_, err = trg.Create(startTime, permNoWrite)
	chk.NoErr(err)

	chk.NoErr(err)

	chk.NoErr(purge.Directory(rootDir))
	chk.Log()
	chk.Stdout()
}
