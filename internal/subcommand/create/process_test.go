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

package create_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/subcommand/create"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestCreate_Process_NoArgs(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg"})
	outText, err := create.Process(args)
	chk.Err(
		err,
		""+
			create.ErrInvalid.Error()+
			": "+
			szargs.ErrMissing.Error()+
			": source directory"+
			"",
	)
	chk.Str(outText, "")
}

func TestCreate_Process_InvalidSourceDirectory(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg", "/INVALID_source"})
	outText, err := create.Process(args)
	chk.Err(
		err,
		""+
			create.ErrInvalid.Error()+
			": "+
			settings.ErrCreate.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": '/INVALID_source'"+
			"",
	)
	chk.Str(outText, "")
}

func TestCreate_Process_OutFileExists(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	source := chk.CreateTmpDir()
	outFile := chk.CreateTmpFile(nil)

	args := szargs.New("", []string{"prg", "-o", outFile, source})
	outText, err := create.Process(args)

	chk.Err(
		err,
		""+
			create.ErrInvalid.Error()+
			": "+
			create.ErrOutFileExists.Error()+
			": '"+outFile+"'"+
			"",
	)
	chk.Str(outText, "")
}

func TestCreate_Process_Valid_Stdout(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	source := chk.CreateTmpSubDir("source")

	args := szargs.New("", []string{"prg", source})
	outText, err := create.Process(args)
	chk.NoErr(err)

	chkCfgText := strings.Replace(
		settings.DefaultConfig,
		"source: /home/user",
		"source: "+source,
		1,
	)

	chk.StrSlice(
		strings.Split(outText, "\n"),
		strings.Split(chkCfgText, "\n"),
	)
}

func TestCreate_Process_Valid_OutFile(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()
	source := chk.CreateTmpSubDir("source")
	toFile := filepath.Join(dir, "file.sbc")

	args := szargs.New("", []string{"prg", "-o", toFile, source})
	outText, err := create.Process(args)

	chk.NoErr(err)
	chk.Str(
		outText,
		"successfully created: "+toFile,
	)

	_, err = os.Stat(toFile)
	chk.NoErr(err)
}
