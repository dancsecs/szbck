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

package trim_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/subcommand/trim"
	"github.com/dancsecs/szbck/internal/target"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztest"
	"github.com/dancsecs/sztestlog"
)

//nolint:goCheckNoGlobals // Ok.
var (
	rsyncCmd      string
	rootTime      time.Time
	monDec31_2018 time.Time
	tueDec31_2024 time.Time
	wedDec31_2025 time.Time
	thuDec31_2020 time.Time
	friDec31_2021 time.Time
	satDec31_2022 time.Time
	sunDec31_2023 time.Time
)

//nolint:goCheckNoInits // Ok.
func init() {
	var err error

	rsyncCmd, err = exec.LookPath("rsync")
	if err != nil {
		panic(err)
	}

	// Tuesday May 6, 2025 07:08:09.1234
	rootTime = time.Date(2025, time.May, 6, 7, 8, 9, 123400000, time.Local)

	monDec31_2018 = time.Date(
		2018, time.December, 31, 12, 0, 0, 0, time.Local,
	)
	tueDec31_2024 = time.Date(
		2024, time.December, 31, 12, 0, 0, 0, time.Local,
	)
	wedDec31_2025 = time.Date(
		2025, time.December, 31, 12, 0, 0, 0, time.Local,
	)
	thuDec31_2020 = time.Date(
		2020, time.December, 31, 12, 0, 0, 0, time.Local,
	)
	friDec31_2021 = time.Date(
		2021, time.December, 31, 12, 0, 0, 0, time.Local,
	)
	satDec31_2022 = time.Date(
		2022, time.December, 31, 12, 0, 0, 0, time.Local,
	)
	sunDec31_2023 = time.Date(
		2023, time.December, 31, 12, 0, 0, 0, time.Local,
	)
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
		"#option: --verbose",
		"option: --quiet",
		1,
	)

	cfgFile := filepath.Join(dir, "backup.sbc")

	chk.NoErr(
		os.WriteFile(cfgFile, []byte(bckCfg), 0o0600),
	)

	return cfgFile
}

func makeSnapshotDir(chk *sztest.Chk, dir string, tme time.Time) string {
	chk.T().Helper()

	trg, err := target.New(dir)
	chk.NoErr(err)

	trgBk, err := trg.Create(
		tme,
		0o0700,
	)
	chk.NoErr(err)

	chk.NoErr(trg.SetLatest(trgBk))

	return trgBk
}

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

func TestTrim_Process_NoArgs(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg"})
	outText, err := trim.Process(args)
	chk.Err(
		err,
		""+
			trim.ErrTrimError.Error()+
			" (Purged: 0): "+
			szargs.ErrMissing.Error()+
			": backup config filename"+
			"",
	)
	chk.Str(outText, "")
}

func TestTrim_Process_BlankBackupDir(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)

	args := szargs.New("", []string{"prg", cfgFile})
	outText, err := trim.Process(args)
	chk.Err(
		err,
		""+
			trim.ErrTrimError.Error()+
			" (Purged: 0): "+
			settings.ErrNoTarget.Error()+
			"",
	)
	chk.Str(outText, "")
}

func TestTrim_Process_EmptyBackupDir(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)
	trg := chk.CreateTmpSubDir("target")

	args := szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err := trim.Process(args)
	chk.Err(
		err,
		""+
			trim.ErrTrimError.Error()+
			" (Purged: 0): "+
			trim.ErrNoBackups.Error()+
			"",
	)
	chk.Str(outText, "")
}

func TestTrim_Process_OnlyOneBackupDir(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	// A little before now without sleeping.
	startTime := time.Now().Add(-time.Millisecond)

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	_ = makeSnapshotDir(chk, trgDir, startTime)

	args := szargs.New("", []string{"prg", "-t", trgDir, cfgFile})
	outText, err := trim.Process(args)
	chk.Err(
		err,
		""+
			trim.ErrTrimError.Error()+
			" (Purged: 0): "+
			trim.ErrOnlyLatest.Error()+
			"",
	)
	chk.Str(outText, "")
}

func TestTrim_Process_OnlyOneInvalidBackupDir(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	// A little before now without sleeping.
	startTime := time.Now().Add(-time.Millisecond)

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	tme := startTime
	snapshotPath := makeSnapshotDir(chk, trgDir, tme)
	tme = tme.Add(time.Minute * 30)
	_ = makeSnapshotDir(chk, trgDir, tme)

	chk.NoErr(
		os.Rename(
			snapshotPath,
			snapshotPath+target.BackupDirectoryExtension,
		),
	)

	args := szargs.New("", []string{"prg", "-t", trgDir, cfgFile})
	outText, err := trim.Process(args)
	chk.Err(
		err,
		""+
			trim.ErrTrimError.Error()+
			" (Purged: 0): "+
			trim.ErrInvalidSnapshotName.Error()+
			": '"+snapshotPath+target.BackupDirectoryExtension+"'"+
			"",
	)
	chk.Str(outText, "")
}

func TestTrim_Process_TwoBackupDirs_DryRun(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	// A little before now without sleeping.
	startTime := time.Now().Add(-time.Millisecond)

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	tme := startTime.Add(-time.Minute * 30)
	snap1 := makeSnapshotDir(chk, trgDir, tme) // Thirty minutes ago.
	snap2 := makeSnapshotDir(chk, trgDir, startTime)

	args := szargs.New(
		"",
		[]string{"prg", "--dry-run", "-t", trgDir, cfgFile},
	)
	outText, err := trim.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	chk.Log()
	chk.Stdout(
		"Keeping snapshot (DRY RUN): "+fmtTS(snap1),
		"Keeping snapshot (DRY RUN): "+fmtTS(snap2),
		"trim successful (Purged: 0) (DRY RUN)\n"+
			"Before: 12,312 After: 12,312 Total Recovered: 0 bytes",
	)
	chk.Stderr()
}

func TestTrim_Process_TwoBackupDirs_PurgeNoneDryRun(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	// A little before now without sleeping.
	startTime := time.Now().Add(-time.Millisecond)

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	tme := startTime.Add(-time.Hour * 24 * 356)
	snap1 := makeSnapshotDir(chk, trgDir, tme) // Last Year.
	snap2 := makeSnapshotDir(chk, trgDir, startTime)

	args := szargs.New(
		"",
		[]string{"prg", "--dry-run", "-t", trgDir, cfgFile},
	)
	outText, err := trim.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	chk.Log()
	chk.Stdout(
		"Keeping snapshot (DRY RUN): "+fmtTS(snap1),
		"Keeping snapshot (DRY RUN): "+fmtTS(snap2),
		"trim successful (Purged: 0) (DRY RUN)\n"+
			"Before: 12,312 After: 12,312 Total Recovered: 0 bytes\n",
	)
	chk.Stderr()
}

func TestTrim_Process_TwoBackupDirs_PurgeNone(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	// A little before now without sleeping.
	startTime := time.Now().Add(-time.Millisecond)

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	tme := startTime.Add(-time.Hour * 24 * 356)
	snap1 := makeSnapshotDir(chk, trgDir, tme) // Last Year.
	snap2 := makeSnapshotDir(chk, trgDir, startTime)

	args := szargs.New("", []string{"prg", "-t", trgDir, cfgFile})
	outText, err := trim.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	chk.Log()
	chk.Stdout(
		"Keeping snapshot: "+fmtTS(snap1),
		"Keeping snapshot: "+fmtTS(snap2),
		"trim successful (Purged: 0)\n"+
			"Before: 12,312 After: 12,312 Total Recovered: 0 bytes",
	)
	chk.Stderr()
}

func TestTrim_Process_TwoBackupDirs_PurgeDaily(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	// A little before now without sleeping.
	startTime := time.Now().Add(-time.Millisecond)

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	// Four days ago at noon
	tme := startTime.Add(-time.Hour * 24 * 4).Truncate(time.Hour * 12)
	purge1 := makeSnapshotDir(chk, trgDir, tme) // Last Year.
	tme = tme.Add(time.Hour)
	purge2 := makeSnapshotDir(chk, trgDir, tme)
	tme = tme.Add(time.Hour)
	keep3 := makeSnapshotDir(chk, trgDir, tme)

	keepRoot := makeSnapshotDir(chk, trgDir, startTime)

	args := szargs.New("", []string{"prg", "-t", trgDir, cfgFile})
	outText, err := trim.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	chk.Log()
	chk.Stdout(
		"Purging snapshot: "+fmtTS(purge1),
		"Purging snapshot: "+fmtTS(purge2),
		"Keeping snapshot: "+fmtTS(keep3),
		"Keeping snapshot: "+fmtTS(keepRoot),
		"trim successful (Purged: 2)\n"+
			"Before: 20,504 After: 12,312 Total Recovered: 8,192 bytes\n",
	)
	chk.Stderr()
}
