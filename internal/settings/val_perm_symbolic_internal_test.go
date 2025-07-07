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

func TestSettingsInternal_PermSymbolic_InvalidBlank(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateSymbolicPermission("")
	chk.Err(
		err,
		""+
			ErrSymbolic.Error()+
			": "+
			ErrSymbolicNone.Error()+
			"",
	)
	chk.Uint32(perm, 0)
}

func TestSettingsInternal_PermSymbolic_InvalidGroupsMin(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateSymbolicPermission(";")
	chk.Err(
		err,
		""+
			ErrSymbolic.Error()+
			": "+
			ErrSyntax.Error()+
			"",
	)
	chk.Uint32(perm, 0)
}

func TestSettingsInternal_PermSymbolic_InvalidGroupsMax(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateSymbolicPermission(";;;")
	chk.Err(
		err,
		""+
			ErrSymbolic.Error()+
			": "+
			ErrSyntax.Error()+
			"",
	)
	chk.Uint32(perm, 0)
}

func TestSettingsInternal_PermSymbolic_InvalidGroupsIdentifier(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateSymbolicPermission("q:rwx;g:rx;o:-")
	chk.Err(
		err,
		""+
			ErrSymbolic.Error()+
			": "+
			ErrSyntax.Error()+
			"",
	)
	chk.Uint32(perm, 0)

	perm, err = validateSymbolicPermission("u:rwx;q:rx;o:-")
	chk.Err(
		err,
		""+
			ErrSymbolic.Error()+
			": "+
			ErrSyntax.Error()+
			"",
	)
	chk.Uint32(perm, 0)

	perm, err = validateSymbolicPermission("u:rwx;g:rx;q:-")
	chk.Err(
		err,
		""+
			ErrSymbolic.Error()+
			": "+
			ErrSyntax.Error()+
			"",
	)
	chk.Uint32(perm, 0)
}

func TestSettingsInternal_PermSymbolic_InvalidPermissionsNone(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateSymbolicPermission("u:-;g:-;o:-")
	chk.Err(
		err,
		""+
			ErrSymbolic.Error()+
			": "+
			ErrSymbolicNone.Error()+
			"",
	)
	chk.Uint32(perm, 0)
}

func TestSettingsInternal_PermSymbolic_ValidUser(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateSymbolicPermission("u:rwx;g:-;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0700)

	perm, err = validateSymbolicPermission("u:rw;g:-;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0600)

	perm, err = validateSymbolicPermission("u:rx;g:-;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0500)

	perm, err = validateSymbolicPermission("u:r;g:-;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0400)

	perm, err = validateSymbolicPermission("u:wx;g:-;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0300)

	perm, err = validateSymbolicPermission("u:w;g:-;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0200)

	perm, err = validateSymbolicPermission("u:x;g:-;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0100)
}

func TestSettingsInternal_PermSymbolic_ValidGroup(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateSymbolicPermission("u:-;g:rwx;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0070)

	perm, err = validateSymbolicPermission("u:-;g:rw;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0060)

	perm, err = validateSymbolicPermission("u:-;g:rx;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0050)

	perm, err = validateSymbolicPermission("u:-;g:r;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0040)

	perm, err = validateSymbolicPermission("u:-;g:wx;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0030)

	perm, err = validateSymbolicPermission("u:-;g:w;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0020)

	perm, err = validateSymbolicPermission("u:-;g:x;o:-")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0010)
}

func TestSettingsInternal_PermSymbolic_ValidOther(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateSymbolicPermission("u:-;g:-;o:rwx")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0007)

	perm, err = validateSymbolicPermission("u:-;g:-;o:rw")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0006)

	perm, err = validateSymbolicPermission("u:-;g:-;o:rx")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0005)

	perm, err = validateSymbolicPermission("u:-;g:-;o:r")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0004)

	perm, err = validateSymbolicPermission("u:-;g:-;o:wx")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0003)

	perm, err = validateSymbolicPermission("u:-;g:-;o:w")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0002)

	perm, err = validateSymbolicPermission("u:-;g:-;o:x")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0001)
}

func TestSettingsInternal_PermSymbolic_ValidAll(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	perm, err := validateSymbolicPermission("u:rwx;g:rwx;o:rwx")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0777)

	perm, err = validateSymbolicPermission("u:rw;g:rw;o:rw")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0666)

	perm, err = validateSymbolicPermission("u:rx;g:rx;o:rx")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0555)

	perm, err = validateSymbolicPermission("u:r;g:r;o:r")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0444)

	perm, err = validateSymbolicPermission("u:wx;g:wx;o:wx")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0333)

	perm, err = validateSymbolicPermission("u:w;g:w;o:w")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0222)

	perm, err = validateSymbolicPermission("u:x;g:x;o:x")
	chk.NoErr(err)
	chk.Uint32(perm, 0o0111)
}
