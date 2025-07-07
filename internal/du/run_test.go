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

package du_test

import (
	"os"
	"testing"

	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szbck/internal/du"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestDuRun_InvalidDirectory(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	out, err := du.Run(nil, os.Stderr)

	chk.Err(
		err,
		""+
			du.ErrDuError.Error()+
			": "+
			du.ErrMissing.Error()+
			": command: ' '"+
			"",
	)
	chk.Str(out, "")

	out, err = du.Run([]string{""}, os.Stderr)

	chk.Err(
		err,
		""+
			du.ErrDuError.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": ''"+
			": command: ' '"+
			"",
	)
	chk.Str(out, "")

	out, err = du.Run([]string{"INVALID_DIRECTORY"}, os.Stderr)

	chk.Err(
		err,
		""+
			du.ErrDuError.Error()+
			": "+
			directory.ErrInvalid.Error()+
			": 'INVALID_DIRECTORY'"+
			": command: ' INVALID_DIRECTORY'"+
			"",
	)
	chk.Str(out, "")

	chk.Log()
	chk.Stdout()
	chk.Stderr()
}

func TestDuRun_EmptyDirectory(t *testing.T) {
	chk := sztestlog.CaptureLogAndStderrAndStdout(t, szlog.LevelAll)
	defer chk.Release()

	dir := chk.CreateTmpDir()

	out, err := du.Run([]string{"-b", "-d", "1", dir}, os.Stderr)

	chk.NoErr(err)
	chk.Str(out, "4096\t"+dir+"\n")

	chk.Log()
	chk.Stdout()
	chk.Stderr()
}
