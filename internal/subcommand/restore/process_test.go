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

package restore_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/rsync"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/subcommand/restore"
	"github.com/dancsecs/szbck/internal/subcommand/snapshot"
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

func TestRestore_MakeDirs(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	const backupRootDate = "20251212_151617.1234.szb"

	// Create source directory chain as
	// /dir/trgRoot/user/abc/def/ghi
	srcRoot := chk.CreateTmpSubDir(backupRootDate)

	srcUser := chk.CreateTmpSubDir(filepath.Join(srcRoot, "user"))

	srcAbc := chk.CreateTmpSubDir(filepath.Join(srcUser, "abc"))
	srcDef := chk.CreateTmpSubDir(filepath.Join(srcAbc, "def"))
	srcGhi := chk.CreateTmpSubDir(filepath.Join(srcDef, "ghi"))

	//nolint:wrapcheck,err113 // Ok.
	tst := func(fromPath, toPath, eFrom, eTo string) error {
		fromDir, toDir, err := restore.MakeDirs(fromPath, toPath)
		if err != nil {
			return err
		}

		if fromDir != eFrom {
			return fmt.Errorf("From: %s -> expected: %s", fromDir, eFrom)
		}

		if toDir != eTo {
			return fmt.Errorf("To: %s -> expected: %s", toDir, eTo)
		}

		return err
	}

	validTests := [][4]string{
		{srcRoot, "/home/user", srcUser, "/home"},
		{srcUser, "/home/user", srcUser, "/home"},
		{srcAbc, "/home/user", srcAbc, "/home/user"},
		{srcDef, "/home/user", srcDef, "/home/user/abc"},
		{srcGhi, "/home/user", srcGhi, "/home/user/abc/def"},
	}

	for i, vt := range validTests {
		chk.NoErr(tst(vt[0], vt[1], vt[2], vt[3]),
			fmt.Sprintf("Valid Test: %d", i),
		)
	}

	chk.Err(
		tst(srcGhi, "/home/user2", srcGhi, "/home/user2/abc/def"),
		""+
			restore.ErrInvalidSrcPath.Error()+
			": 'user/abc/def/ghi' must start with 'user2'"+
			"",
	)
}

func TestRestoreProcess_noArgs(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg"})
	textOut, err := restore.Process(args)

	chk.Err(
		err,
		""+
			restore.ErrRestoreError.Error()+
			": "+
			szargs.ErrMissing.Error()+
			": backup config filename"+
			"",
	)
	chk.Str(textOut, "")
}

func TestRestoreProcess_InvalidSpecificDirectory(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	_ = chk.CreateTmpDir()
	_, cfgFile := setupBackupConfig(chk)
	trg := chk.CreateTmpSubDir("target")

	args := szargs.New("", []string{"prg", "-s", ".", "-t", trg, cfgFile})
	textOut, err := restore.Process(args)

	chk.Err(
		err,
		""+
			restore.ErrRestoreError.Error()+
			": "+
			target.ErrInvalidSplit.Error()+
			": lstat "+filepath.Join(trg, "latest")+
			": no such file or directory"+
			"",
	)
	chk.Str(textOut, "")
}

//nolint:funlen // Ok.
func TestRestoreProcess_FromLatestTargetDirectory(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	source, cfgFile := setupBackupConfig(chk)
	trg := chk.CreateTmpSubDir("target")

	file := chk.CreateTmpFileIn(source, []byte("file1"))

	args := szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err := snapshot.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	args = szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err = restore.Process(args)

	chk.NoErr(err)
	chk.Str(outText, "restore successful")

	chk.NoErr(os.Remove(file))

	args = szargs.New(
		"",
		[]string{
			"prg",
			"-s",
			target.LatestDirectoryLink,
			"-t",
			trg,
			cfgFile,
		},
	)
	outText, err = restore.Process(args)

	chk.NoErr(err)
	chk.Str(outText, "restore successful")

	chk.AddSub(`\d+`, "#")
	chk.Log()
	chk.Stdout(
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+source+
			" "+filepath.Join(trg, "#_#.#.szb")+
			"",
		"snapshot successful\n"+
			"Before: 8,192 After: 12,317 Used: 4,125 bytes",
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+filepath.Join(trg, "#_#.#.szb", "source")+
			" "+dir+
			"",
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+filepath.Join(trg, "#_#.#.szb", "source")+
			" "+dir+
			"",
		"",
	)
}

func TestRestoreProcess_FromBaseTargetDirectory(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	source, cfgFile := setupBackupConfig(chk)
	trg := chk.CreateTmpSubDir("target")

	file := chk.CreateTmpFileIn(source, []byte("file1"))

	args := szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err := snapshot.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	args = szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err = restore.Process(args)

	chk.NoErr(err)
	chk.Str(outText, "restore successful")

	chk.NoErr(os.Remove(file))

	args = szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err = restore.Process(args)

	chk.NoErr(err)
	chk.Str(outText, "restore successful")

	chk.AddSub(`\d+`, "#")
	chk.Log()
	chk.Stdout(
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+source+
			" "+filepath.Join(trg, "#_#.#.szb")+
			"",
		"snapshot successful\n"+
			"Before: 8,192 After: 12,317 Used: 4,125 bytes",
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+filepath.Join(trg, "#_#.#.szb", "source")+
			" "+dir+
			"",
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+filepath.Join(trg, "#_#.#.szb", "source")+
			" "+dir+
			"",
		"",
	)
}

func TestRestoreProcess_DryRun_And_Keep(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	source, cfgFile := setupBackupConfig(chk)
	trg := chk.CreateTmpSubDir("target")

	file := chk.CreateTmpFileIn(source, []byte("file1"))

	args := szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err := snapshot.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	args = szargs.New("", []string{"prg", "--dry-run", "-t", trg, cfgFile})
	outText, err = restore.Process(args)

	chk.NoErr(err)
	chk.Str(outText, "restore successful")

	chk.NoErr(os.Remove(file))

	args = szargs.New("", []string{"prg", "--keep", "-t", trg, cfgFile})
	outText, err = restore.Process(args)

	chk.NoErr(err)
	chk.Str(outText, "restore successful")

	chk.AddSub(`\d+`, "#")
	chk.Log()
	chk.Stdout(
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+source+
			" "+filepath.Join(trg, "#_#.#.szb")+
			"",
		"snapshot successful\n"+
			"Before: 8,192 After: 12,317 Used: 4,125 bytes",
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+rsync.FlgDryRun+ // --dry-run
			" "+filepath.Join(trg, "#_#.#.szb", "source")+
			" "+dir+
			"",
		"Running command: "+
			rsyncCmd+basicOptions+
			// " "+rsync.FlgDelete+  --keep
			" "+filepath.Join(trg, "#_#.#.szb", "source")+
			" "+dir+
			"",
		"",
	)
}

//nolint:funlen // Ok.
func TestRestoreProcess_SpecificDirectoryTree(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	// dir := chk.CreateTmpDir()
	source, cfgFile := setupBackupConfig(chk)
	trg := chk.CreateTmpSubDir("target")

	srcSubDir1 := chk.CreateTmpSubDir("source", "subDir1")

	srcSubDir2 := chk.CreateTmpSubDir("source", "subDir2")

	fileSub1 := chk.CreateTmpFileIn(srcSubDir1, []byte("fileSub1"))

	fileSub2 := chk.CreateTmpFileIn(srcSubDir2, []byte("fileSub2"))

	args := szargs.New("", []string{"prg", "-t", trg, cfgFile})
	outText, err := snapshot.Process(args)
	chk.NoErr(err)
	chk.Str(outText, "")

	args = szargs.New(
		"",
		[]string{
			"prg",
			"-s",
			filepath.Join(target.LatestDirectoryLink, "source", "subDir2"),
			"-t",
			trg,
			cfgFile,
		},
	)
	outText, err = restore.Process(args)

	chk.NoErr(err)
	chk.Str(outText, "restore successful")

	chk.NoErr(os.Remove(fileSub1))
	chk.NoErr(os.Remove(fileSub2))

	args = szargs.New(
		"",
		[]string{
			"prg",
			"-s",
			filepath.Join(target.LatestDirectoryLink, "source", "subDir2"),
			"-t",
			trg,
			cfgFile,
		},
	)
	outText, err = restore.Process(args)

	chk.NoErr(err)
	chk.Str(outText, "restore successful")

	_, err = os.Stat(fileSub1) // deleted and not restored
	chk.Err(
		err,
		""+
			"stat "+fileSub1+": no such file or directory"+
			"",
	)

	_, err = os.Stat(fileSub2) // deleted and restored
	chk.NoErr(err)

	chk.AddSub(`\d+`, "#")

	chk.Log()
	chk.Stdout(
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+source+
			" "+filepath.Join(trg, "#_#.#.szb")+
			"",
		"snapshot successful\n"+
			"Before: 8,192 After: 20,520 Used: 12,328 bytes",
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+filepath.Join(trg, "#_#.#.szb", "source", "subDir2")+
			" "+source+
			"",
		"Running command: "+
			rsyncCmd+basicOptions+
			" "+rsync.FlgDelete+
			" "+filepath.Join(trg, "#_#.#.szb", "source", "subDir2")+
			" "+source+
			"",
		"",
	)
}
