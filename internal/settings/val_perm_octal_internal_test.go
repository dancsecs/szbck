/*
   Golang rsync backup utility  wrapper: szbck.
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

package settings

import (
	"testing"

	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestSettingsInternal_PermOctal_InvalidBlank(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateOctalPermission("")
	chk.Err(
		err,
		""+
			ErrOctal.Error()+
			": "+
			ErrSyntax.Error()+
			"",
	)
	chk.Uint32(perm, 0)
}

func TestSettingsInternal_PermOctal_InvalidSyntax(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateOctalPermission("0o0abc")
	chk.Err(
		err,
		""+
			ErrOctal.Error()+
			": "+
			ErrSyntax.Error()+
			"",
	)
	chk.Uint32(perm, 0)
}

func TestSettingsInternal_PermOctal_InvalidPrefix(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateOctalPermission("9")
	chk.Err(
		err,
		""+
			ErrOctal.Error()+
			": "+
			ErrSyntax.Error()+
			"",
	)
	chk.Uint32(perm, 0)
}

func TestSettingsInternal_PermOctal_InvalidParseRange(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateOctalPermission("0o077777777777777777777777777777777")
	chk.Err(
		err,
		""+
			ErrOctal.Error()+
			": "+
			ErrRange.Error()+
			"",
	)
	chk.Uint32(perm, 0)
}

func TestSettingsInternal_PermOctal_InvalidPermissionRangeMin(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateOctalPermission("0o00")
	chk.Err(
		err,
		""+
			ErrOctal.Error()+
			": "+
			ErrRange.Error()+
			"",
	)
	chk.Uint32(perm, 0)
}

func TestSettingsInternal_PermOctal_InvalidPermissionRangeMax(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateOctalPermission("0o07770")
	chk.Err(
		err,
		""+
			ErrOctal.Error()+
			": "+
			ErrRange.Error()+
			"",
	)
	chk.Uint32(perm, 0)
}

func TestSettingsInternal_PermOctal_Valid(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateOctalPermission("0o023")
	chk.NoErr(err)
	chk.Uint32(perm, 0o23)
}
