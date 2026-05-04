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

package du_test

import (
	"testing"

	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szbck/internal/du"
	"github.com/dancsecs/sztestlog"
)

func TestDu_Total_InvalidDirectory(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	size, err := du.Total("")
	chk.Err(
		err,
		""+
			du.ErrInvalid.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": ''"+
			"",
	)
	chk.Int64(size, 0)

	size, err = du.Total("DOES_NOT_EXIST")
	chk.Err(
		err,
		""+
			du.ErrInvalid.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": 'DOES_NOT_EXIST'"+
			"",
	)
	chk.Int64(size, 0)
}

func TestDu_Total_EmptyDirectory(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	dir := chk.CreateTmpDir()

	_, err := du.Total(dir)
	chk.NoErr(err)
}

func TestDu_Total_OneFile(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	_ = chk.CreateTmpFile([]byte("sample file"))

	_, err := du.Total(dir)
	chk.NoErr(err)
}

func TestDu_Total_InvalidDirectories(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	size1, size2, err := du.Totals("", "")
	chk.Err(
		err,
		""+
			du.ErrInvalid.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": ''"+
			"",
	)
	chk.Int64(size1, 0)
	chk.Int64(size2, 0)

	size1, size2, err = du.Totals("DOES_NOT_EXIST", "")
	chk.Err(
		err,
		""+
			du.ErrInvalid.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": 'DOES_NOT_EXIST'"+
			"",
	)
	chk.Int64(size1, 0)
	chk.Int64(size2, 0)

	goodDir := chk.CreateTmpDir()

	size1, size2, err = du.Totals(goodDir, "")
	chk.Err(
		err,
		""+
			du.ErrInvalid.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": ''"+
			"",
	)
	chk.Int64(size1, 0)
	chk.Int64(size2, 0)

	size1, size2, err = du.Totals(goodDir, "DOES_NOT_EXIST")
	chk.Err(
		err,
		""+
			du.ErrInvalid.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": 'DOES_NOT_EXIST'"+
			"",
	)
	chk.Int64(size1, 0)
	chk.Int64(size2, 0)
}

func TestDu_Total_EmptyDirectories(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	dir1 := chk.CreateTmpSubDir("A")
	dir2 := chk.CreateTmpSubDir("B")

	_, _, err := du.Totals(dir1, dir2)
	chk.NoErr(err)
}

func TestDu_Total_OneFile_OneFile(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	dir1 := chk.CreateTmpSubDir("A")
	_ = chk.CreateTmpFileIn(dir1, []byte("sample file"))

	dir2 := chk.CreateTmpSubDir("B")
	_ = chk.CreateTmpFileIn(dir2, []byte("sample file"))

	_, _, err := du.Totals(dir1, dir2)
	chk.NoErr(err)
}
