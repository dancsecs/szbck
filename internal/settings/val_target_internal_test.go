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

	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestInternalSettings_ValidateTarget_InvalidBlank(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.Err(
		cfg.validateTarget(""),
		""+
			ErrTarget.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": ''"+
			"",
	)
}

func TestInternalSettings_ValidateTarget_InvalidDirectory(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	fDir := chk.CreateTmpDir()
	fName := chk.CreateTmpFileAs(fDir, "NOT_A_DIRECTORY", nil)

	chk.Err(
		cfg.validateTarget(fName),
		""+
			ErrTarget.Error()+
			": "+
			directory.ErrNotADirectory.Error()+
			": '"+fName+"'"+
			"",
	)
}

func TestInternalSettings_ValidateTarget_Valid(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	fDir := chk.CreateTmpDir()
	chk.NoErr(cfg.validateTarget(fDir))
}

func TestInternalSettings_ValidateTarget_InvalidDuplicate(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	fDir := chk.CreateTmpDir()

	chk.NoErr(cfg.validateTarget(fDir))

	chk.Err(
		cfg.validateTarget(fDir),
		""+
			ErrTarget.Error()+
			": "+
			ErrDuplicate.Error()+
			"",
	)
}
