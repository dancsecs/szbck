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

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestSettings_Load_DoesNotExist(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	cfg, err := settings.Load("DOES_NOT_EXIST")

	chk.Nil(cfg)
	chk.Err(
		err,
		""+
			settings.ErrLoad.Error()+
			": open DOES_NOT_EXIST: no such file or directory"+
			"",
	)
}

func TestConfigBackup_Load_Valid(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	src := chk.CreateTmpSubDir("source")
	trg := chk.CreateTmpSubDir("target")

	cfgData, err := settings.Create(src, trg)
	chk.NoErr(err)

	cfgData = strings.Replace(
		cfgData,
		"#snapshotOption: --some-option",
		"snapshotOption: --some-option",
		1,
	)

	cfgData = strings.Replace(
		cfgData,
		"#restoreOption: --one-flag",
		"restoreOption: --one-flag",
		1,
	)

	cfgFile := chk.CreateTmpFileAs("", "sample.sbc", []byte(cfgData))

	cfg, err := settings.Load(cfgFile)
	chk.NoErr(err)

	chk.Str(cfg.Source, src)
	chk.Str(cfg.Target.GetPath(), trg)
	chk.StrSlice(cfg.SnapshotOptions, []string{"--some-option"})
	chk.StrSlice(cfg.RestoreOptions, []string{"--one-flag"})
}

func TestConfigBackup_LoadFromArgs_NoArgs(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg"})
	cfg, err := settings.LoadFromArgs(args)
	chk.Err(
		err,
		""+
			szargs.ErrMissing.Error()+
			": backup config filename"+
			"",
	)
	chk.Nil(cfg)
}

func TestConfigBackup_LoadFromArgs_NoneDefined(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	src := chk.CreateTmpSubDir("source")

	cfgData, err := settings.Create(src, "")
	chk.NoErr(err)

	cfgFile := chk.CreateTmpFileAs("", "sample.sbc", []byte(cfgData))

	args := szargs.New("", []string{"prg", cfgFile})
	cfg, err := settings.LoadFromArgs(args)
	chk.Err(
		err,
		""+
			settings.ErrNoTarget.Error()+
			"",
	)
	chk.Nil(cfg)
}

func TestConfigBackup_LoadFromArgs_NoDefaultOverrideOnly(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	src := chk.CreateTmpSubDir("source")
	trg := chk.CreateTmpSubDir("target")

	cfgData, err := settings.Create(src, "")
	chk.NoErr(err)

	cfgFile := chk.CreateTmpFileAs("", "sample.sbc", []byte(cfgData))

	args := szargs.New("", []string{"prg", "-t", trg, cfgFile})
	cfg, err := settings.LoadFromArgs(args)
	chk.NoErr(err)
	chk.Str(cfg.Target.GetPath(), trg)
}

func TestConfigBackup_LoadFromArgs_DefaultNoOverride(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	src := chk.CreateTmpSubDir("source")
	trg := chk.CreateTmpSubDir("target")

	cfgData, err := settings.Create(src, trg)
	chk.NoErr(err)

	cfgFile := chk.CreateTmpFileAs("", "sample.sbc", []byte(cfgData))

	args := szargs.New("", []string{"prg", cfgFile})
	cfg, err := settings.LoadFromArgs(args)
	chk.NoErr(err)
	chk.Str(cfg.Target.GetPath(), trg)
}

func TestConfigBackup_LoadFromArgs_DefaultAndOverride(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	src := chk.CreateTmpSubDir("source")
	trg1 := chk.CreateTmpSubDir("target1")
	trg2 := chk.CreateTmpSubDir("target2")

	cfgData, err := settings.Create(src, trg1)
	chk.NoErr(err)

	cfgFile := chk.CreateTmpFileAs("", "sample.sbc", []byte(cfgData))

	args := szargs.New("", []string{"prg", "-t", trg2, cfgFile})
	cfg, err := settings.LoadFromArgs(args)

	chk.NoErr(err)
	chk.Str(cfg.Target.GetPath(), trg2)
}
