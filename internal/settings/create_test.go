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

package settings_test

import (
	"strings"
	"testing"

	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/target"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestSettings_Create_InvalidSource(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	fDir := chk.CreateTmpDir()
	fName := chk.CreateTmpFileAs(fDir, "NOT_A_DIRECTORY", nil)

	cfg, err := settings.Create(fName, "")
	chk.Err(
		err,
		""+
			settings.ErrCreate.Error()+
			": "+
			directory.ErrNotADirectory.Error()+
			": '"+fName+"'"+
			"",
	)
	chk.Str(cfg, "")
}

func TestSettings_Create_InvalidTarget(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	fDir := chk.CreateTmpDir()
	fName := chk.CreateTmpFileAs(fDir, "NOT_A_DIRECTORY", nil)

	cfg, err := settings.Create(fDir, fName)
	chk.Err(
		err,
		""+
			settings.ErrCreate.Error()+
			": "+
			target.ErrNew.Error()+
			": "+
			target.ErrInvalid.Error()+
			": "+
			directory.ErrNotADirectory.Error()+
			": '"+fName+"'"+
			"",
	)
	chk.Str(cfg, "")
}

func TestConfigBackup_Create_Valid(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	src := chk.CreateTmpSubDir("source")
	trg := chk.CreateTmpSubDir("target")

	cfgData := strings.Replace(
		settings.DefaultConfig,
		"source: /home/user",
		"source: "+src,
		1,
	)

	cfgData = strings.Replace(
		cfgData,
		"#target: /mnt/backupDir",
		"target: "+trg,
		1,
	)

	cfg, err := settings.Create(src, trg)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(cfg, "\n"),
		strings.Split(cfgData, "\n"),
	)
}
