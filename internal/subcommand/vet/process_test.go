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

package vet_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/subcommand/vet"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztest"
	"github.com/dancsecs/sztestlog"
)

func setupBackupConfig(chk *sztest.Chk, setSourceError bool) string {
	chk.T().Helper()

	dir := chk.CreateTmpDir()
	source := chk.CreateTmpSubDir("source")

	bckCfg, err := settings.Create(source, "")
	chk.NoErr(err)

	// Restore write permission to new snapshot directories.
	bckCfg = strings.Replace(
		bckCfg,
		"permission: 0o0500",
		"permission: 0o0700",
		1,
	)

	// Remove verbose.  Not testing rsync output just results.
	bckCfg = strings.ReplaceAll(
		bckCfg,
		"option: --verbose",
		"",
	)

	if setSourceError {
		bckCfg = strings.Replace(
			bckCfg,
			"source: "+source,
			"source: /home/DOES_NOT_EXIST",
			1,
		)
	}

	cfgFile := filepath.Join(dir, "backup.sbc")

	chk.NoErr(
		os.WriteFile(cfgFile, []byte(bckCfg), 0o0600),
	)

	return cfgFile
}

func TestVet_Process_NoArgs(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	outText, err := vet.Process(nil)
	chk.Err(
		err,
		""+
			vet.ErrVetError.Error()+
			": "+
			szargs.ErrMissing.Error()+
			": backup config filename"+
			"",
	)
	chk.Str(outText, "")
}

func TestVet_Process_SourceError(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk, true)

	outText, err := vet.Process([]string{cfgFile})
	chk.Err(
		err,
		""+
			vet.ErrVetError.Error()+
			": "+
			settings.ErrLoad.Error()+
			": "+
			settings.ErrInvalid.Error()+
			": "+
			settings.ErrConfigLine.Error()+
			"(7): "+
			settings.ErrSource.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": '/home/DOES_NOT_EXIST'"+
			"\n\tsource: /home/DOES_NOT_EXIST"+
			"",
	)
	chk.Str(outText, "")
}

func TestVet_Process_Valid(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	cfgFile := setupBackupConfig(chk, false)

	outText, err := vet.Process([]string{cfgFile})
	chk.NoErr(err)
	chk.Str(outText, "vet successful (no problems found)\n")
}
