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

package help_test

import (
	"strings"
	"testing"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/subcommand/create"
	"github.com/dancsecs/szbck/internal/subcommand/help"
	"github.com/dancsecs/szbck/internal/subcommand/prune"
	"github.com/dancsecs/szbck/internal/subcommand/restore"
	"github.com/dancsecs/szbck/internal/subcommand/snapshot"
	"github.com/dancsecs/szbck/internal/subcommand/status"
	"github.com/dancsecs/szbck/internal/subcommand/trim"
	"github.com/dancsecs/szbck/internal/subcommand/vet"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestHelpProcess_UnknownSubcommand(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg", "unknown"})
	helpText, err := help.Process(args)

	chk.Str(helpText, "")

	chk.Err(
		err,
		""+
			help.ErrHelpError.Error()+
			": "+
			help.ErrUnknownSubcommand.Error()+
			": 'unknown'"+
			"",
	)
}

func TestHelpProcess_NoSubcommand(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg"})
	helpText, err := help.Process(args)
	chk.NoErr(err)

	wantTxt := strings.Split(help.Usage, "\n")
	wantTxt = append(wantTxt, strings.Split(help.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(create.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(snapshot.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(restore.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(prune.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(status.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(trim.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(vet.HelpText, "\n")...)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		wantTxt,
	)
}

func TestHelpProcess_All(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg", "ALL"})
	helpText, err := help.Process(args)
	chk.NoErr(err)

	wantTxt := strings.Split(help.Usage, "\n")
	wantTxt = append(wantTxt, strings.Split(help.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(create.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(snapshot.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(restore.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(prune.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(status.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(trim.HelpText, "\n")...)
	wantTxt = append(wantTxt, strings.Split(vet.HelpText, "\n")...)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		wantTxt,
	)
}

func TestHelpProcess_Help(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg", "H"})
	helpText, err := help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(help.HelpText, "\n"),
	)

	args = szargs.New("", []string{"prg", "HELP"})
	helpText, err = help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(help.HelpText, "\n"),
	)
}

func TestHelpProcess_Initialize(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg", "I"})
	helpText, err := help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(create.HelpText, "\n"),
	)

	args = szargs.New("", []string{"prg", "INIT"})
	helpText, err = help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(create.HelpText, "\n"),
	)

	args = szargs.New("", []string{"prg", "INITIALIZE"})
	helpText, err = help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(create.HelpText, "\n"),
	)
}

func TestHelpProcess_Snapshot(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg", "S"})
	helpText, err := help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(snapshot.HelpText, "\n"),
	)

	args = szargs.New("", []string{"prg", "SNAP"})
	helpText, err = help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(snapshot.HelpText, "\n"),
	)

	args = szargs.New("", []string{"prg", "SNAPSHOT"})
	helpText, err = help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(snapshot.HelpText, "\n"),
	)
}

func TestHelpProcess_Restore(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg", "R"})
	helpText, err := help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(restore.HelpText, "\n"),
	)

	args = szargs.New("", []string{"prg", "RES"})
	helpText, err = help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(restore.HelpText, "\n"),
	)

	args = szargs.New("", []string{"prg", "RESTORE"})
	helpText, err = help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(restore.HelpText, "\n"),
	)
}

func TestHelpProcess_Prune(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg", "P"})
	helpText, err := help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(prune.HelpText, "\n"),
	)

	args = szargs.New("", []string{"prg", "PRUNE"})
	helpText, err = help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(prune.HelpText, "\n"),
	)
}

func TestHelpProcess_Status(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg", "STAT"})
	helpText, err := help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(status.HelpText, "\n"),
	)

	args = szargs.New("", []string{"prg", "STATUS"})
	helpText, err = help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(status.HelpText, "\n"),
	)
}

func TestHelpProcess_Trim(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg", "T"})
	helpText, err := help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(trim.HelpText, "\n"),
	)

	args = szargs.New("", []string{"prg", "TRIM"})
	helpText, err = help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(trim.HelpText, "\n"),
	)
}

func TestHelpProcess_Vet(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	args := szargs.New("", []string{"prg", "V"})
	helpText, err := help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(vet.HelpText, "\n"),
	)

	args = szargs.New("", []string{"prg", "VET"})
	helpText, err = help.Process(args)
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(vet.HelpText, "\n"),
	)
}
