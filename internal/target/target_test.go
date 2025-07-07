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

package target_test

import (
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szbck/internal/target"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestTarget_New_Blank(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	trg, err := target.New("")

	chk.Err(
		err,
		""+
			target.ErrNew.Error()+
			": "+
			target.ErrInvalid.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": ''"+
			"",
	)
	chk.Nil(trg)
}

func TestTarget_New_DoesNotExist(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	trg, err := target.New("DOES_NOT_EXIST")

	chk.Err(
		err,
		""+
			target.ErrNew.Error()+
			": "+
			target.ErrInvalid.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": 'DOES_NOT_EXIST'"+
			"",
	)
	chk.Nil(trg)
}

func TestTarget_HasLatest_NotALink(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	trg, err := target.New(dir)
	chk.NoErr(err)

	_ = chk.CreateTmpSubDir(target.LatestDirectoryLink)

	hasLatest, err := trg.HasLatest()
	chk.Err(
		err,
		""+
			target.ErrHasLatest.Error()+
			": "+
			target.ErrInvalidLatest.Error()+
			"",
	)
	chk.False(hasLatest)
}

func TestTarget_HasLatest_Valid_EmptyDirectory(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	trg, err := target.New(dir)
	chk.NoErr(err)

	hasLatest, err := trg.HasLatest()
	chk.NoErr(err)
	chk.False(hasLatest)

	chk.Str(trg.GetPath(), dir)
}

func TestTarget_HasLatest_Valid_HasLink(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	trg, err := target.New(dir)
	chk.NoErr(err)

	_ = chk.CreateTmpSubDir("dirToLink")

	chk.NoErr(
		directory.LinkRelative(
			filepath.Join(dir, "dirToLink"),
			filepath.Join(dir, target.LatestDirectoryLink),
		),
	)

	hasLatest, err := trg.HasLatest()
	chk.NoErr(err)
	chk.True(hasLatest)
}

func TestConfigBackup_Create_AlreadyExists(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	trg, err := target.New(dir)
	chk.NoErr(err)

	tme := time.Date(2025, time.May, 2, 3, 4, 5, 333999000, time.Local)
	tmeStr := tme.Format(target.BackupDirectoryFormat) +
		target.BackupDirectoryExtension

	dirToLink := chk.CreateTmpSubDir(tmeStr)
	chk.NoErr(
		directory.LinkRelative(
			dirToLink,
			filepath.Join(dir, target.LatestDirectoryLink),
		),
	)

	path, err := trg.Create(tme, 0o0700)
	chk.Err(
		err,
		""+
			target.ErrCreateTargetFailed.Error()+
			": "+
			target.ErrCreateAlreadyExists.Error()+
			": '"+filepath.Join(dir, tmeStr)+
			"'"+
			"",
	)
	chk.Str(path, "")
}

func TestConfigBackup_Create_Valid(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	trg, err := target.New(dir)
	chk.NoErr(err)

	tme := time.Date(2025, time.May, 2, 3, 4, 5, 333999000, time.Local)

	path, err := trg.Create(tme, 0o0700)
	chk.NoErr(err)
	chk.Str(
		path,
		filepath.Join(
			dir,
			tme.Format(target.BackupDirectoryFormat)+
				target.BackupDirectoryExtension,
		),
	)
}

func TestConfigBackup_SetLatest(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	trg, err := target.New(dir)
	chk.NoErr(err)

	backupDir := chk.CreateTmpSubDir("backup")

	chk.Err(
		trg.SetLatest(backupDir+"x"),
		""+
			target.ErrInvalidLatest.Error()+
			": "+
			directory.ErrCreateLink.Error()+
			": (from: '"+backupDir+"x"+
			"' to: '"+
			trg.Latest()+"'): "+
			directory.ErrInvalid.Error()+
			": '"+backupDir+"x"+
			"'"+
			"",
	)

	chk.NoErr(trg.SetLatest(backupDir))
}

//nolint:funlen // Ok.
func TestDirectory_Split(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	reSplit := regexp.MustCompile(`abc`)

	dir := chk.CreateTmpDir()

	pre, post, err := target.Split(dir, reSplit)
	chk.Err(
		err,
		""+
			target.ErrInvalidSplit.Error()+
			": "+
			target.ErrSplitNotFound.Error()+
			"",
	)
	chk.Str(pre, "")
	chk.Str(post, "")

	abcDir := chk.CreateTmpSubDir("abc")

	pre, post, err = target.Split(abcDir, reSplit)
	chk.NoErr(err)
	chk.Str(pre, filepath.Join(dir, "abc"))
	chk.Str(post, "")

	abcdefDir := chk.CreateTmpSubDir("abc/def")

	pre, post, err = target.Split(abcdefDir, reSplit)
	chk.NoErr(err)
	chk.Str(pre, filepath.Join(dir, "abc"))
	chk.Str(post, "def")

	abcdefghiDir := chk.CreateTmpSubDir("abc/def/ghi")

	pre, post, err = target.Split(abcdefghiDir, reSplit)
	chk.NoErr(err)
	chk.Str(pre, filepath.Join(dir, "abc"))
	chk.Str(post, "def/ghi")

	// With trailing separator.
	pre, post, err = target.Split(
		abcDir+directory.PathSeparator,
		reSplit,
	)
	chk.NoErr(err)
	chk.Str(pre, filepath.Join(dir, "abc"))
	chk.Str(post, "")

	pre, post, err = target.Split(
		abcdefDir+directory.PathSeparator,
		reSplit,
	)
	chk.NoErr(err)
	chk.Str(pre, filepath.Join(dir, "abc"))
	chk.Str(post, "def")

	pre, post, err = target.Split(
		abcdefghiDir+directory.PathSeparator,
		reSplit,
	)
	chk.NoErr(err)
	chk.Str(pre, filepath.Join(dir, "abc"))
	chk.Str(post, "def/ghi")
}
