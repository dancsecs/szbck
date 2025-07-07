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
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestSettings_Parse_InvalidSource(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	cfgData := strings.Replace(
		settings.DefaultConfig,
		"source: /home/user",
		"source: INVALID",
		1,
	)

	cfg, err := settings.Parse(cfgData)
	chk.Err(
		err,
		""+
			settings.ErrInvalid.Error()+
			": "+
			settings.ErrConfigLine.Error()+
			"(7): "+
			settings.ErrSource.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": 'INVALID'"+
			"\n\tsource: INVALID"+
			"",
	)
	chk.Nil(cfg)
}

func TestSettings_Parse_UnknownKey(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	cfgData := "unknownKey: unknownValue"

	cfg, err := settings.Parse(cfgData)
	chk.Err(
		err,
		""+
			settings.ErrInvalid.Error()+
			": "+
			settings.ErrConfigLine.Error()+
			"(1): "+
			settings.ErrUnknownKey.Error()+
			": 'unknownKey'"+
			"\n\tunknownKey: unknownValue"+
			"",
	)
	chk.Nil(cfg)
}

func TestSettings_Parse_InvalidEntry(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	cfgData := "noValue"

	cfg, err := settings.Parse(cfgData)
	chk.Err(
		err,
		""+
			settings.ErrInvalid.Error()+
			": "+
			settings.ErrConfigLine.Error()+
			"(1): "+
			settings.ErrInvalidSyntax.Error()+
			"\n\tnoValue"+
			"",
	)
	chk.Nil(cfg)
}

func TestSettings_Parse_InvalidTarget(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	srcDir := chk.CreateTmpSubDir("source")

	cfgData := strings.Replace(
		settings.DefaultConfig,
		"source: /home/user",
		"source: "+srcDir,
		1,
	)

	cfgData = strings.Replace(
		cfgData,
		"#target: /mnt/backupDir",
		"target: INVALID",
		1,
	)

	cfg, err := settings.Parse(cfgData)
	chk.Err(
		err,
		""+
			settings.ErrInvalid.Error()+
			": "+
			settings.ErrConfigLine.Error()+
			"(12): "+
			settings.ErrTarget.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": 'INVALID'"+
			"\n\ttarget: INVALID"+
			"",
	)
	chk.Nil(cfg)
}
