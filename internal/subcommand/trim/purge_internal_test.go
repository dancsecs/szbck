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

package trim

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dancsecs/szbck/internal/target"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

const (
	permNoWrite = 0o0500
	permWrite   = 0o0700
)

func fmtTS(fName string) string {
	fileTime, _ := time.ParseInLocation(
		target.BackupDirectoryFormat,
		strings.TrimSuffix(
			filepath.Base(fName),
			target.BackupDirectoryExtension,
		),
		time.Local,
	)

	return fName + ": " + fileTime.Format(time.RFC850)
}

func TestInternalTrim_ProcessPurge_RootPermissionFailure(t *testing.T) {
	chk := sztestlog.CaptureLog(t, szlog.LevelAll)
	defer chk.Release()

	// A little before now without sleeping.
	startTime := time.Now().Add(-time.Millisecond)

	rootDir := chk.CreateTmpSubDir("root")
	trg, err := target.New(rootDir)
	chk.NoErr(err)

	dirToDelete, err := trg.Create(startTime, permNoWrite)
	chk.NoErr(err)

	chk.NoErr(os.Chmod(rootDir, permNoWrite))

	defer func() {
		_ = os.Chmod(rootDir, permWrite)
	}()

	purgedCount, err := processPurge(
		[]string{dirToDelete},
		[]time.Time{startTime},
		[]bool{true},
		"",
	)

	chk.Err(
		err,
		""+
			ErrPurgeFailed.Error()+
			": unlinkat "+dirToDelete+
			": permission denied"+
			"",
	)
	chk.Int(purgedCount, 0)

	chk.Log()
}

func TestInternalTrim_ProcessPurge_FailureAfterSuccess(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	// A little before now without sleeping.
	startTime := time.Now().Add(-time.Millisecond)

	rootDir := chk.CreateTmpSubDir("root")
	trg, err := target.New(rootDir)
	chk.NoErr(err)

	dirToDelete, err := trg.Create(startTime, permNoWrite)
	chk.NoErr(err)

	purgedCount, err := processPurge(
		[]string{dirToDelete, "INVALID_DIR"},
		[]time.Time{startTime, startTime},
		[]bool{true, true},
		"",
	)

	chk.Err(
		err,
		""+
			ErrPurgeFailed.Error()+
			": chmod INVALID_DIR"+
			": no such file or directory"+
			"",
	)
	chk.Int(purgedCount, 1)

	chk.Log()
	chk.Stdout(
		"Purging snapshot: " + fmtTS(dirToDelete),
	)
}

func TestInternalTrim_ProcessPurge_Success(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	// A little before now without sleeping.
	startTime := time.Now().Add(-time.Millisecond)

	rootDir := chk.CreateTmpSubDir("root")
	trg, err := target.New(rootDir)
	chk.NoErr(err)

	dirToKeep, err := trg.Create(startTime.Add(time.Minute), permNoWrite)
	chk.NoErr(err)
	chk.NoErr(trg.SetLatest(dirToKeep))
	dirToDelete, err := trg.Create(startTime, permNoWrite)
	chk.NoErr(err)

	purgedCount, err := processPurge(
		[]string{dirToDelete, dirToKeep},
		[]time.Time{startTime, startTime.Add(time.Minute)},
		[]bool{true, false},
		"",
	)

	chk.NoErr(err)
	chk.Int(purgedCount, 1)

	chk.Log()
	chk.Stdout(
		"Purging snapshot: "+fmtTS(dirToDelete),
		"Keeping snapshot: "+fmtTS(dirToKeep),
	)
}
