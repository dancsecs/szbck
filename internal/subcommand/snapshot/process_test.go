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

package snapshot_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/rsync"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/subcommand/snapshot"
	"github.com/dancsecs/szbck/internal/subcommand/trim"
	"github.com/dancsecs/szbck/internal/target"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztest"
	"github.com/dancsecs/sztestlog"
)

const basicOptions = "" +
	" " + "--archive" +
	" " + "--quiet" +
	" " + "--human-readable" +
	" " + "--acls" +
	" " + "--xattrs" +
	" " + "--atimes" +
	" " + "--hard-links" +
	" " + "--fsync" +
	" " + "--exclude=.cache" +
	" " + "--exclude=go/pkg" +
	""

//nolint:goCheckNoGlobals // Ok.
var rsyncCmd string

//nolint:goCheckNoInits // Ok.
func init() {
	var err error

	rsyncCmd, err = exec.LookPath("rsync")
	if err != nil {
		panic(err)
	}
}

func setupBackupConfig(chk *sztest.Chk) (string, string) {
	chk.T().Helper()

	dir := chk.CreateTmpDir()
	source := chk.CreateTmpSubDir("source")

	bckCfg, err := settings.Create(source, "")
	chk.NoErr(err)

	bckCfg = strings.Replace(
		bckCfg,
		"permission: 0o0500",
		"permission: 0o0700",
		1,
	)

	bckCfg = strings.Replace(
		bckCfg,
		"#option: --verbose",
		"option: --quiet",
		1,
	)

	cfgFile := filepath.Join(dir, "backup.sbc")

	chk.NoErr(
		os.WriteFile(cfgFile, []byte(bckCfg), 0o0600),
	)

	return source, cfgFile
}

func TestSnapshotProcess_NextHourIn(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	chk.Dur(snapshot.NextHourIn(0), time.Hour)
	chk.Dur(snapshot.NextHourIn(time.Minute*31), time.Hour+time.Minute*29)
	chk.Dur(snapshot.NextHourIn(time.Minute*29), time.Minute*31)
}

func TestSnapshotProcess_MissingArgs(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg"})
	outText, err := snapshot.Process(args)
	chk.Err(
		err,
		""+
			snapshot.ErrSnapshotError.Error()+
			": "+
			szargs.ErrMissing.Error()+
			": backup config filename"+
			"",
	)
	chk.Str(outText, "")
}

func TestSnapshotProcess_NoFiles(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	source, cfgFile := setupBackupConfig(chk)
	trg := chk.CreateTmpSubDir("target")

	args := szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err := snapshot.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	chk.AddSub(`\d+`, "#")
	chk.Log()

	chk.Stdout(
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+"--delete"+
			" "+source+
			" "+filepath.Join(trg, "#_#.#"+target.BackupDirectoryExtension),
		"snapshot successful\n"+
			"Before: 8,192 After: 12,312 Used: 4,120 bytes\n",
	)
}

func TestSnapshotProcess_DryRun(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	source, cfgFile := setupBackupConfig(chk)
	trg := chk.CreateTmpSubDir("target")

	_ = chk.CreateTmpFileIn(source, []byte("file"))

	args := szargs.New(
		"",
		[]string{"prg", "--dry-run", "-t", trg, cfgFile},
	)
	outText, err := snapshot.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	chk.AddSub(`\d+`, "#")
	chk.Log()
	chk.Stdout(
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+rsync.FlgDryRun+
			" "+source+
			" "+filepath.Join(trg, "#_#.#"+target.BackupDirectoryExtension),
		"snapshot successful (DRY RUN)\n"+
			"Before: #,# After: #,# Used: # bytes",
	)
}

func TestSnapshotProcess_OneFile(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	source, cfgFile := setupBackupConfig(chk)
	trg := chk.CreateTmpSubDir("target")

	_ = chk.CreateTmpFileIn(source, []byte("file"))

	args := szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err := snapshot.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	args = szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err = snapshot.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	chk.AddSub(`\d+`, "#")
	chk.Log()
	chk.Stdout(
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+source+
			" "+filepath.Join(trg, "#_#.#"+target.BackupDirectoryExtension),
		"snapshot successful\n"+
			"Before: #,# After: #,# Used: #,# bytes",
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+rsync.FlgLinkDest+
			filepath.Join(trg, target.LatestDirectoryLink)+
			" "+source+
			" "+filepath.Join(trg, "#_#.#"+target.BackupDirectoryExtension),
		"snapshot successful\n"+
			"Before: #,# After: #,# Used: #,# bytes",
	)
}

func TestSnapshotProcess_TwoFiles(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	source, cfgFile := setupBackupConfig(chk)
	trg := chk.CreateTmpSubDir("target")

	_ = chk.CreateTmpFileIn(source, []byte("file1"))

	args := szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err := snapshot.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	_ = chk.CreateTmpFileIn(source, []byte("file1"))

	args = szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err = snapshot.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	chk.AddSub(`\d+`, "#")
	chk.Log()
	chk.Stdout(
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+source+
			" "+filepath.Join(trg, "#_#.#"+target.BackupDirectoryExtension),
		"snapshot successful\n"+
			"Before: #,# After: #,# Used: #,# bytes",
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+rsync.FlgLinkDest+
			filepath.Join(trg, target.LatestDirectoryLink)+
			" "+source+
			" "+filepath.Join(trg, "#_#.#"+target.BackupDirectoryExtension),
		"snapshot successful\n"+
			"Before: #,# After: #,# Used: #,# bytes",
	)
}

func TestSnapshotProcess_Trim(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	source, cfgFile := setupBackupConfig(chk)
	trg := chk.CreateTmpSubDir("target")

	args := szargs.New(
		"",
		[]string{"prg", "--trim", "-t", trg, cfgFile},
	)
	outText, err := snapshot.Process(args)

	chk.Err(
		err,
		""+
			snapshot.ErrSnapshotError.Error()+
			" (Total Purged: 0): "+
			trim.ErrOnlyLatest.Error()+
			"",
	)
	chk.Str(outText, "")

	chk.AddSub(`\d+`, "#")
	chk.Log()
	chk.Stdout(
		"Running command: " +
			rsyncCmd + basicOptions +
			" " + rsync.FlgDelete +
			" " + source +
			" " + filepath.Join(trg, "#_#.#"+target.BackupDirectoryExtension),
	)
}
