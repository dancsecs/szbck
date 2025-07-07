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

package out_test

import (
	"testing"

	"github.com/dancsecs/szbck/internal/out"
	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestFormat_Int(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	chk.Str(
		out.Int(0),
		"0",
	)

	chk.Str(
		out.Int(3567823256),
		"3,567,823,256",
	)

	chk.Str(
		out.Int(-43567823256),
		"-43,567,823,256",
	)
}

func TestOut_Print(t *testing.T) {
	chk := sztestlog.CaptureStdout(t, szlog.LevelAll)
	defer chk.Release()

	szlog.SetLevel(szlog.LevelError)

	out.Print("This line will not be displayed")
	out.Printf("This %s will not be displayed", "formatted line")

	szlog.SetLevel(szlog.LevelAll)

	out.Print("This line will be displayed\n")
	out.Printf("This %s will be displayed\n", "formatted line")

	chk.Stdout(
		"This line will be displayed",
		"This formatted line will be displayed",
		"",
	)
}
