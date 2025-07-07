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

func TestInternalSettings_ValidateKeyValue_Unknown(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.Err(
		cfg.validateKeyValue("UNKNOWN", ""),
		""+
			ErrUnknownKey.Error()+
			": 'UNKNOWN'"+
			"",
	)
}

func TestInternalSettings_ValidateKeyValue_SnapshotOption(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.NoErr(cfg.validateKeyValue("snapshotOption", "-v"))
	chk.StrSlice(cfg.SnapshotOptions, []string{"-v"})
}

func TestInternalSettings_ValidateKeyValue_RestoreOption(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.NoErr(cfg.validateKeyValue("restoreOption", "-v"))
	chk.StrSlice(cfg.RestoreOptions, []string{"-v"})
}
