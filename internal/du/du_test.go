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

package du_test

import (
	"testing"

	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szbck/internal/du"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestDu_Total_InvalidDirectory(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
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
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()

	size, err := du.Total(dir)
	chk.NoErr(err)
	chk.Int64(size, 4096)
}

func TestDu_Total_OneFile(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	_ = chk.CreateTmpFile([]byte("sample file"))

	size, err := du.Total(dir)
	chk.NoErr(err)
	chk.Int64(size, 4107)
}
