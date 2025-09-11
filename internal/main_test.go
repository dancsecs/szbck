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
	"github.com/dancsecs/sztestlog"
)

// Main is the actual mainline for the puzzle application classically
// returning an int to be returned when exiting.
func TestBackupMain_MissingSubCommand(t *testing.T) {
	chk := sztestlog.CaptureLog(t)
	defer chk.Release()

	args := []string{"programName"}

	chk.Int(
		internal.Main(args),
		1,
	)

	chk.Log(
		"F:programName - missing argument: sub command",
	)
}

func TestBackupMain_UnknownSubCommand(t *testing.T) {
	chk := sztestlog.CaptureLog(t)
	defer chk.Release()

	args := []string{"programName", "UnknownSubCommand"}

	chk.Int(
		internal.Main(args),
		1,
	)

	chk.Log(
		"F:programName - " +
			internal.ErrUnknownSubcommand.Error() +
			": 'UnknownSubCommand'",
	)
}

func TestBackupMain_Help(t *testing.T) {
	chk := sztestlog.CaptureStdout(t)
	defer chk.Release()

	args := []string{"programName", "h"}

	chk.Int(
		internal.Main(args),
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
	)
}

func TestBackupMain_Create(t *testing.T) {
	chk := sztestlog.CaptureLog(t)
	defer chk.Release()

	args := []string{"programName", "c"}

	chk.Int(
		internal.Main(args),
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
	chk := sztestlog.CaptureLog(t)
	defer chk.Release()

	args := []string{"programName", "s"}

	chk.Int(
		internal.Main(args),
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
	chk := sztestlog.CaptureLog(t)
	defer chk.Release()

	args := []string{"programName", "r"}

	chk.Int(
		internal.Main(args),
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
	chk := sztestlog.CaptureLog(t)
	defer chk.Release()

	args := []string{"programName", "p"}

	chk.Int(
		internal.Main(args),
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
	chk := sztestlog.CaptureLog(t)
	defer chk.Release()

	args := []string{"programName", "stat"}

	chk.Int(
		internal.Main(args),
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
	chk := sztestlog.CaptureLog(t)
	defer chk.Release()

	args := []string{"programName", "trim"}

	chk.Int(
		internal.Main(args),
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
	chk := sztestlog.CaptureLog(t)
	defer chk.Release()

	args := []string{"programName", "vet"}

	chk.Int(
		internal.Main(args),
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

func TestArgUsage_Dedication(t *testing.T) {
	chk := sztestlog.CaptureLogAndStdout(t)
	defer chk.Release()

	args := []string{
		"noProgName",
		"--Reem",
	}

	chk.Int(
		internal.Main(args),
		1,
	)

	chk.Stdout(`
*****************************************************************************
**                                                                         **
** This project is dedicated to Reem.                                      **
** Your brilliance, courage, and quiet strength continue to inspire me.    **
** Every line is written in gratitude for the light and hope you brought   **
** into my life.                                                           **
**                                                                         **
*****************************************************************************
`)

	chk.Log("F:noProgName - missing argument: sub command")
}
