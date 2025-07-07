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

package internal_test

import (
	"testing"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal"
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

// Main is the actual mainline for the puzzle application classically
// returning an int to be returned when exiting.
func TestBackupMain_MissingSubCommand(t *testing.T) {
	chk := sztestlog.CaptureLog(t, szlog.LevelAll)
	defer chk.Release()

	chk.Int(
		internal.Main("programName", nil),
		1,
	)

	chk.Int(
		internal.Main("programName", []string{}),
		1,
	)

	chk.Log(
		"F:programName - missing argument: sub command",
		"F:programName - missing argument: sub command",
	)
}

func TestBackupMain_UnknownSubCommand(t *testing.T) {
	chk := sztestlog.CaptureLog(t, szlog.LevelAll)
	defer chk.Release()

	chk.Int(
		internal.Main("programName", []string{"UnknownSubCommand"}),
		1,
	)

	chk.Log(
		"F:programName - " +
			internal.ErrUnknownSubcommand.Error() +
			": 'UnknownSubCommand'",
	)
}

func TestBackupMain_Verbose(t *testing.T) {
	chk := sztestlog.CaptureLog(t, szlog.LevelFatal)
	defer chk.Release()

	chk.Int(
		internal.Main(
			"programName", []string{"-v", "-v", "-vv", "UnknownSubCommand"},
		),
		1,
	)

	chk.Int(
		int(szlog.Level()),
		int(szlog.LevelDebug),
	)

	chk.Log(
		"F:programName - " +
			internal.ErrUnknownSubcommand.Error() +
			": 'UnknownSubCommand'",
	)
}

func TestBackupMain_Help(t *testing.T) {
	chk := sztestlog.CaptureStdout(t, szlog.LevelAll)
	defer chk.Release()

	chk.Int(
		internal.Main("programName", []string{"h"}),
		0,
	)

	chk.AddSub(`^[\t\s]+`, "")
	chk.Stdout(
		"programName",
		help.Usage,
		help.HelpText,
		create.HelpText,
		snapshot.HelpText,
		restore.HelpText,
		prune.HelpText,
		status.HelpText,
		trim.HelpText,
		vet.HelpText,
		"",
	)
}

func TestBackupMain_Create(t *testing.T) {
	chk := sztestlog.CaptureLog(t, szlog.LevelAll)
	defer chk.Release()

	chk.Int(
		internal.Main("programName", []string{"c"}),
		1,
	)

	chk.Log(
		"" +
			"F:programName - " +
			create.ErrInvalid.Error() +
			": " +
			szargs.ErrMissing.Error() +
			": source directory",
	)
}

func TestBackupMain_Snapshot(t *testing.T) {
	chk := sztestlog.CaptureLog(t, szlog.LevelAll)
	defer chk.Release()

	chk.Int(
		internal.Main("programName", []string{"s"}),
		1,
	)

	chk.Log(
		"" +
			"F:programName - " +
			snapshot.ErrSnapshotError.Error() +
			": " +
			szargs.ErrMissing.Error() +
			": backup config filename" +
			"",
	)
}

func TestBackupMain_Restore(t *testing.T) {
	chk := sztestlog.CaptureLog(t, szlog.LevelAll)
	defer chk.Release()

	chk.Int(
		internal.Main("programName", []string{"r"}),
		1,
	)

	chk.Log(
		"" +
			"F:programName - " +
			restore.ErrRestoreError.Error() +
			": " +
			szargs.ErrMissing.Error() +
			": backup config filename" +
			"",
	)
}

func TestBackupMain_Prune(t *testing.T) {
	chk := sztestlog.CaptureLog(t, szlog.LevelAll)
	defer chk.Release()

	chk.Int(
		internal.Main("programName", []string{"p"}),
		1,
	)

	chk.Log(
		"" +
			"F:programName - " +
			prune.ErrPruneError.Error() +
			": " +
			szargs.ErrMissing.Error() +
			": backup config filename" +
			"",
	)
}

func TestBackupMain_Status(t *testing.T) {
	chk := sztestlog.CaptureLog(t, szlog.LevelAll)
	defer chk.Release()

	chk.Int(
		internal.Main("programName", []string{"stat"}),
		1,
	)

	chk.Log(
		"" +
			"F:programName - " +
			status.ErrStatusError.Error() +
			": " +
			szargs.ErrMissing.Error() +
			": backup config filename" +
			"",
	)
}

func TestBackupMain_Trim(t *testing.T) {
	chk := sztestlog.CaptureLog(t, szlog.LevelAll)
	defer chk.Release()

	chk.Int(
		internal.Main("programName", []string{"trim"}),
		1,
	)

	chk.Log(
		"" +
			"F:programName - " +
			trim.ErrTrimError.Error() +
			" (Purged: 0): " +
			szargs.ErrMissing.Error() +
			": backup config filename" +
			"",
	)
}

func TestBackupMain_Vet(t *testing.T) {
	chk := sztestlog.CaptureLog(t, szlog.LevelAll)
	defer chk.Release()

	chk.Int(
		internal.Main("programName", []string{"vet"}),
		1,
	)

	chk.Log(
		"" +
			"F:programName - " +
			vet.ErrVetError.Error() +
			": " +
			szargs.ErrMissing.Error() +
			": backup config filename" +
			"",
	)
}
