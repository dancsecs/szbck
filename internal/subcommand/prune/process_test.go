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

package prune_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/subcommand/prune"
	"github.com/dancsecs/szbck/internal/target"
	"github.com/dancsecs/sztest"
	"github.com/dancsecs/sztestlog"
)

const (
	summaryUsage = "" +
		"                        Bytes" +
		"                         INodes\n" +
		" Totals:                    #" +
		"                              #\n" +
		" Before:                    # (     #%)" +
		"                    # (     #%)\n" +
		"  After:                    # (     #%)" +
		"                    # (     #%)\n" +
		"   Used:                    # (     #%)" +
		"                    # (     #%)"
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

func squashNumbers(chk *sztest.Chk) {
	chk.AddSub(`\((?:\-|\s|\d)\d\d\.\d\d\%\)`, "(     #%)")
	chk.AddSub(`\(\s(?:\-|\s|\d)\d\.\d\d\%\)`, "(     #%)")

	chk.AddSub(`\-\d\,\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d`,
		"                    #")
	chk.AddSub(`\-\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d`,
		"                   #")
	chk.AddSub(`\-\d\d\,\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d`, "                  #")
	chk.AddSub(`\-\d\,\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d`, "                 #")
	chk.AddSub(`\-\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d`, "               #")
	chk.AddSub(`\-\d\d\,\d\d\d\,\d\d\d\,\d\d\d`, "              #")
	chk.AddSub(`\-\d\,\d\d\d\,\d\d\d\,\d\d\d`, "             #")
	chk.AddSub(`\-\d\d\d\,\d\d\d\,\d\d\d`, "           #")
	chk.AddSub(`\-\d\d\,\d\d\d\,\d\d\d`, "          #")
	chk.AddSub(`\-\d\,\d\d\d\,\d\d\d`, "         #")
	chk.AddSub(`\-\d\d\d\,\d\d\d`, "       #")
	chk.AddSub(`\-\d\d\,\d\d\d`, "      #")
	chk.AddSub(`\-\d\,\d\d\d`, "     #")
	chk.AddSub(`\-\d\d\d`, "   #")
	chk.AddSub(`\-\d\d`, "  #")
	chk.AddSub(`\-\d`, " #")

	chk.AddSub(`\d\,\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d`,
		"                    #")
	chk.AddSub(`\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d`, "                  #")
	chk.AddSub(`\d\d\,\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d`, "                 #")
	chk.AddSub(`\d\,\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d`, "                #")
	chk.AddSub(`\d\d\d\,\d\d\d\,\d\d\d\,\d\d\d`, "              #")
	chk.AddSub(`\d\d\,\d\d\d\,\d\d\d\,\d\d\d`, "             #")
	chk.AddSub(`\d\,\d\d\d\,\d\d\d\,\d\d\d`, "            #")
	chk.AddSub(`\d\d\d\,\d\d\d\,\d\d\d`, "          #")
	chk.AddSub(`\d\d\,\d\d\d\,\d\d\d`, "         #")
	chk.AddSub(`\d\,\d\d\d\,\d\d\d`, "        #")
	chk.AddSub(`\d\d\d\,\d\d\d`, "      #")
	chk.AddSub(`\d\d\,\d\d\d`, "     #")
	chk.AddSub(`\d\,\d\d\d`, "    #")
	chk.AddSub(`\d`, "#")
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

func TestPrune_Process_NoArgs(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	args := szargs.New("", []string{"prg"})
	outText, err := prune.Process(args)
	chk.Err(
		err,
		""+
			prune.ErrPruneError.Error()+
			": "+
			szargs.ErrMissing.Error()+
			": backup config filename"+
			"",
	)
	chk.Str(outText, "")
}

func TestPrune_Process_BlankBackupDir(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)

	args := szargs.New("", []string{"prg", cfgFile})
	outText, err := prune.Process(args)
	chk.Err(
		err,
		""+
			prune.ErrPruneError.Error()+
			": "+
			settings.ErrNoTarget.Error()+
			"",
	)
	chk.Str(outText, "")
}

func TestPrune_Process_EmptyBackupDir(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)
	trg := chk.CreateTmpSubDir("target")

	args := szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err := prune.Process(args)
	chk.Err(
		err,
		""+
			prune.ErrPruneError.Error()+
			": "+
			prune.ErrNoBackups.Error()+
			"",
	)
	chk.Str(outText, "")
}

func TestPrune_Process_OnlyOneBackupDir(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	_ = makeSnapshotDir(chk, trgDir, 0)

	args := szargs.New("", []string{"prg", "-t", trgDir, cfgFile})
	outText, err := prune.Process(args)
	chk.Err(
		err,
		""+
			prune.ErrPruneError.Error()+
			": "+
			prune.ErrOnlyLatest.Error()+
			"",
	)
	chk.Str(outText, "")
}

func TestPrune_Process_TwoBackupDirs_DryRun(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	dirToDelete := makeSnapshotDir(chk, trgDir, 0)
	_ = makeSnapshotDir(chk, trgDir, 30)

	args := szargs.New(
		"",
		[]string{"prg", "--dry-run", "-t", trgDir, cfgFile},
	)
	outText, err := prune.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	squashNumbers(chk)
	chk.Log()
	chk.Stdout(
		"Purging oldest backup (DRY RUN)",
		"",
		"Purging backup: "+dirToDelete,
		"prune successful (DRY RUN)",
		"Syncing...",
		summaryUsage,
	)
	chk.Stderr()
}

func TestPrune_Process_TwoBackupDirs_DefaultOne(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	dirToDelete := makeSnapshotDir(chk, trgDir, 0)
	_ = makeSnapshotDir(chk, trgDir, 30)

	args := szargs.New("", []string{"prg", "-t", trgDir, cfgFile})
	outText, err := prune.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	squashNumbers(chk)
	chk.Log()
	chk.Stdout(
		"Purging oldest backup",
		"",
		"Purging backup: "+dirToDelete,
		"prune successful",
		"Syncing...",
		summaryUsage,
	)
	chk.Stderr()
}

func TestPrune_Process_ThreeBackupDirs_DefaultOne(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	dirToDelete := makeSnapshotDir(chk, trgDir, 0)
	_ = makeSnapshotDir(chk, trgDir, 30)
	_ = makeSnapshotDir(chk, trgDir, 60)

	args := szargs.New("", []string{"prg", "-t", trgDir, cfgFile})
	outText, err := prune.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	squashNumbers(chk)
	chk.Log()
	chk.Stdout(
		"Purging oldest backup",
		"",
		"Purging backup: "+dirToDelete,
		"prune successful",
		"Syncing...",
		summaryUsage,
	)
	chk.Stderr()
}

func TestPrune_Process_TwoBackupDirs_All(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	dirToDelete1 := makeSnapshotDir(chk, trgDir, 0)
	dirToDelete2 := makeSnapshotDir(chk, trgDir, 30)
	_ = makeSnapshotDir(chk, trgDir, 60)

	args := szargs.New(
		"",
		[]string{"prg", "-n", "all", "-t", trgDir, cfgFile},
	)
	outText, err := prune.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	squashNumbers(chk)
	chk.Log()
	chk.Stdout(
		"Purging 2 oldest backups",
		"",
		"Purging backup: "+dirToDelete1,
		"Purging backup: "+dirToDelete2,
		"prune successful",
		"Syncing...",
		summaryUsage,
	)
	chk.Stderr()
}

func TestPrune_Process_TwoBackupDirs_InvalidNum(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	_ = makeSnapshotDir(chk, trgDir, 0)
	_ = makeSnapshotDir(chk, trgDir, 30)

	args := szargs.New(
		"",
		[]string{"prg", "-n", "-1", "-t", trgDir, cfgFile},
	)
	outText, err := prune.Process(args)
	chk.Err(
		err,
		""+
			prune.ErrPruneError.Error()+
			": "+
			prune.ErrInvalidNum.Error()+
			": '-1'"+
			"",
	)
	chk.Str(outText, "")

	chk.Log()
	chk.Stdout()
	chk.Stderr()
}

func TestPrune_Process_TwoBackupDirs_TooMany(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk)
	trgDir := chk.CreateTmpSubDir("target")

	dirToDelete1 := makeSnapshotDir(chk, trgDir, 0)
	dirToDelete2 := makeSnapshotDir(chk, trgDir, 30)
	_ = makeSnapshotDir(chk, trgDir, 60)

	args := szargs.New(
		"",
		[]string{"prg", "-n", "1000", "-t", trgDir, cfgFile},
	)
	outText, err := prune.Process(args)

	chk.NoErr(err)
	chk.Str(outText, "")

	squashNumbers(chk)
	chk.Log()
	chk.Stdout(
		"Purging 2 oldest backups",
		"",
		"Purging backup: "+dirToDelete1,
		"Purging backup: "+dirToDelete2,
		"prune successful",
		"Syncing...",
		summaryUsage,
	)
	chk.Stderr()
}
