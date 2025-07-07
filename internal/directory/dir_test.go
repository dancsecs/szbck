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

package directory_test

import (
	"path/filepath"
	"testing"

	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestDirectory_Is_Missing(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	err := directory.Is("")
	chk.Err(
		err,
		""+
			directory.ErrInvalid.Error()+
			": ''"+
			"",
	)

	err = directory.Is("DOES_NOT_EXIST")
	chk.Err(
		err,
		""+
			directory.ErrInvalid.Error()+
			": 'DOES_NOT_EXIST'"+
			"",
	)
}

func TestDirectory_Is_Invalid(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	fDir := chk.CreateTmpDir()
	fName := chk.CreateTmpFileAs(fDir, "NOT_A_DIRECTORY", nil)

	err := directory.Is(fName)
	chk.Err(
		err,
		""+
			directory.ErrNotADirectory.Error()+
			": '"+fName+"'"+
			"",
	)
}

func TestDirectory_Is_Valid(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpSubDir("source")

	err := directory.Is(dir)

	chk.NoErr(err)
}

func TestDirectory_IsEmpty(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpSubDir("target")

	err := directory.IsEmpty(dir)
	chk.NoErr(err)

	_ = chk.CreateTmpFileIn(dir, nil)

	err = directory.IsEmpty(dir)

	chk.Err(
		err,
		directory.ErrNewNotEmpty.Error(),
	)
}

func TestDirectory_Link(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()

	source := chk.CreateTmpSubDir("source")

	link := filepath.Join(dir, "linked")

	err := directory.LinkRelative(source, link)
	chk.NoErr(err)

	notADir := chk.CreateTmpFile(nil)

	err = directory.LinkRelative(notADir, "badLink")

	chk.Err(
		err,
		""+
			directory.ErrCreateLink.Error()+
			": (from: '"+notADir+"' to: 'badLink'): "+
			directory.ErrNotADirectory.Error()+
			": '"+notADir+"'"+
			"",
	)
}
