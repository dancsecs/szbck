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

func TestSettingsInternal_Perm_InvalidBlank(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	err := cfg.validatePermission("")
	chk.Err(
		err,
		""+
			ErrPermission.Error()+
			": "+
			ErrMissing.Error()+
			"",
	)
	chk.Uint32(uint32(cfg.Permission), 0)
}

func TestSettingsInternal_Perm_InvalidDuplicate(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	cfg.Permission = 0o0700

	err := cfg.validatePermission("0o0777")
	chk.Err(
		err,
		""+
			ErrPermission.Error()+
			": "+
			ErrDuplicate.Error()+
			"",
	)
	chk.Uint32(uint32(cfg.Permission), 0o0700)
}

func TestSettingsInternal_Perm_ValidOctal(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	err := cfg.validatePermission("0o0777")
	chk.NoErr(err)
	chk.Uint32(uint32(cfg.Permission), 0o0777)
}

func TestSettingsInternal_Perm_ValidSymbolic(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	err := cfg.validatePermission("u:rwx;g:rx;o:-")
	chk.NoErr(err)
	chk.Uint32(uint32(cfg.Permission), 0o0750)
}
