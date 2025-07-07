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

	helpText, err := help.Process([]string{"unknown"})

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

	helpText, err := help.Process(nil)
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

	helpText, err := help.Process([]string{"ALL"})
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

	helpText, err := help.Process([]string{"H"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(help.HelpText, "\n"),
	)

	helpText, err = help.Process([]string{"HELP"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(help.HelpText, "\n"),
	)
}

func TestHelpProcess_Initialize(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	helpText, err := help.Process([]string{"I"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(create.HelpText, "\n"),
	)

	helpText, err = help.Process([]string{"INIT"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(create.HelpText, "\n"),
	)

	helpText, err = help.Process([]string{"INITIALIZE"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(create.HelpText, "\n"),
	)
}

func TestHelpProcess_Snapshot(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	helpText, err := help.Process([]string{"S"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(snapshot.HelpText, "\n"),
	)

	helpText, err = help.Process([]string{"SNAP"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(snapshot.HelpText, "\n"),
	)

	helpText, err = help.Process([]string{"SNAPSHOT"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(snapshot.HelpText, "\n"),
	)
}

func TestHelpProcess_Restore(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	helpText, err := help.Process([]string{"R"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(restore.HelpText, "\n"),
	)

	helpText, err = help.Process([]string{"RES"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(restore.HelpText, "\n"),
	)

	helpText, err = help.Process([]string{"RESTORE"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(restore.HelpText, "\n"),
	)
}

func TestHelpProcess_Prune(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	helpText, err := help.Process([]string{"P"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(prune.HelpText, "\n"),
	)

	helpText, err = help.Process([]string{"PRUNE"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(prune.HelpText, "\n"),
	)
}

func TestHelpProcess_Status(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	helpText, err := help.Process([]string{"STAT"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(status.HelpText, "\n"),
	)

	helpText, err = help.Process([]string{"STATUS"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(status.HelpText, "\n"),
	)
}

func TestHelpProcess_Trim(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	helpText, err := help.Process([]string{"T"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(trim.HelpText, "\n"),
	)

	helpText, err = help.Process([]string{"TRIM"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(trim.HelpText, "\n"),
	)
}

func TestHelpProcess_Vet(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	helpText, err := help.Process([]string{"V"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(vet.HelpText, "\n"),
	)

	helpText, err = help.Process([]string{"VET"})
	chk.NoErr(err)

	chk.StrSlice(
		strings.Split(helpText, "\n"),
		strings.Split(vet.HelpText, "\n"),
	)
}
