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

package status_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/subcommand/status"
	"github.com/dancsecs/szbck/internal/target"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztest"
	"github.com/dancsecs/sztestlog"
)

//nolint:goCheckNoGlobals // Ok.
var (
	rsyncCmd string
	rootTime time.Time
)

//nolint:goCheckNoInits // Ok.
func init() {
	var err error

	rsyncCmd, err = exec.LookPath("rsync")
	if err != nil {
		panic(err)
	}

	rootTime = time.Date(2025, time.May, 2, 3, 4, 5, 678900000, time.Local)
}

func setupBackupConfig(chk *sztest.Chk) string {
	chk.T().Helper()

	dir := chk.CreateTmpDir()
	source := chk.CreateTmpSubDir("source")

	bckCfg, err := settings.Create(source, "")
	chk.NoErr(err)

	// Restore write permission to new snapshot directories.
	bckCfg = strings.Replace(
		bckCfg,
		"permission: 0o0500",
		"permission: 0o0700",
		1,
	)

	// Remove verbose.  Not testing rsync output just results.
	bckCfg = strings.Replace(
		bckCfg,
		"option: --verbose",
		"",
		2,
	)

	cfgFile := filepath.Join(dir, "backup.sbc")

	chk.NoErr(
		os.WriteFile(cfgFile, []byte(bckCfg), 0o0600),
	)

	return cfgFile
}

func makeSnapshotDir(chk *sztest.Chk, dir string, delta int) string {
	chk.T().Helper()

	trg, err := target.New(dir)
	chk.NoErr(err)

	trgBk, err := trg.Create(
		rootTime.Add(time.Duration(delta)*time.Minute),
		0o0700,
	)
	chk.NoErr(err)

	chk.NoErr(trg.SetLatest(trgBk))

	return trgBk
}

func TestStatus_Process_NoArgs(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	outText, err := status.Process(nil)
	chk.Err(
		err,
		""+
			status.ErrStatusError.Error()+
			": "+
			szargs.ErrMissing.Error()+
			": backup config filename"+
			"",
	)
	chk.Str(outText, "")
}

func TestStatus_Process_InvalidConfigFileDir(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()

	outText, err := status.Process([]string{dir})
	chk.Err(
		err,
		""+
			status.ErrStatusError.Error()+
			": "+
			settings.ErrLoad.Error()+
			": read "+dir+
			": is a directory"+
			"",
	)
	chk.Str(outText, "")
}

func TestStatus_Process_BlankBackupDir(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)

	outText, err := status.Process([]string{cfgFile})
	chk.Err(
		err,
		""+
			status.ErrStatusError.Error()+
			": "+
			settings.ErrNoTarget.Error()+
			"",
	)
	chk.Str(outText, "")
}

func TestStatus_Process_InvalidBackupDir(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	cfgFile := setupBackupConfig(chk)

	outText, err := status.Process([]string{"-t", dir, cfgFile})
	chk.Err(
		err,
		""+
			status.ErrStatusError.Error()+
			": "+
			target.ErrInvalid.Error()+
			": "+
			target.ErrNew.Error()+
			": "+
			target.ErrInvalid.Error()+
			": "+
			directory.ErrNewNotEmpty.Error()+
			"",
	)
	chk.Str(outText, "")
}

func TestStatus_Process_InvalidBackupDirContent(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	cfgFile := setupBackupConfig(chk)
	tme := time.Date(2025, time.May, 2, 3, 4, 5, 333999000, time.Local)

	_ = chk.CreateTmpFileAs(
		dir,
		""+
			tme.Format(target.BackupDirectoryFormat)+
			target.BackupDirectoryExtension,
		nil,
	)

	outText, err := status.Process([]string{"-t", dir, cfgFile})
	chk.Err(
		err,
		""+
			status.ErrStatusError.Error()+
			": "+
			target.ErrInvalid.Error()+
			": "+
			target.ErrNew.Error()+
			": "+
			target.ErrInvalid.Error()+
			": "+
			directory.ErrNewNotEmpty.Error()+
			"",
	)
	chk.Str(outText, "")
}

func TestStatus_Process_EmptyBackupDir(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	outText, err := status.Process([]string{"-t", trgDir, cfgFile})
	chk.Err(
		err,
		""+
			status.ErrStatusError.Error()+
			": "+
			status.ErrReportFailed.Error()+
			": "+
			status.ErrNoBackups.Error()+
			"",
	)
	chk.Str(outText, "")
}

func TestStatus_Process_OneEmptyBackupDir(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	bkDir := makeSnapshotDir(chk, trgDir, 0)

	outText, err := status.Process([]string{"-t", trgDir, cfgFile})
	chk.NoErr(err)
	chk.StrSlice(
		strings.Split(outText, "\n"),
		[]string{
			"status successful",
			"",
			"Backup Sets: 1",
			"Total Bytes: 8,216",
			"",
			filepath.Base(bkDir) + ": 4,096",
			"",
		},
	)
}

func TestStatus_Process_TwoBackupDirsWithOneFile(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	bkDir1 := makeSnapshotDir(chk, trgDir, 0)
	bkDir2 := makeSnapshotDir(chk, trgDir, 30)

	_ = chk.CreateTmpFileIn(bkDir2, []byte("This is a file in dir 2"))

	outText, err := status.Process([]string{"-t", trgDir, cfgFile})
	chk.NoErr(err)
	chk.StrSlice(
		strings.Split(outText, "\n"),
		[]string{
			"status successful",
			"",
			"Backup Sets: 2",
			"Total Bytes: 12,335",
			"",
			filepath.Base(bkDir1) + ": 4,096",
			filepath.Base(bkDir2) + ": 4,119",
			"",
		},
	)
}
