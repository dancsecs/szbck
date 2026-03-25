/*
   Golang rsync backup utility wrapper: szbck.
   Copyright (C) 2026 Leslie Dancsecs

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

package fstat_test

import (
	"strings"
	"testing"

	"github.com/dancsecs/szbck/internal/fstat"
	"github.com/dancsecs/sztestlog"
)

func TestStatfs_InvalidDirectory(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	statfs, err := fstat.New("INVALID_DIRECTORY")

	chk.Nil(statfs)
	chk.Err(
		err,
		chk.ErrChain(
			"statfs failed",
			"no such file or directory",
		),
	)
}

func TestStatfs_Tmp(t *testing.T) {
	chk := sztestlog.CaptureNothing(t)
	defer chk.Release()

	dir := chk.CreateTmpDir()

	statfs, err := fstat.New(dir)

	chk.NoErr(err)
	chk.NotNil(statfs)

	chk.AddSub(`\-?\d[\d\,]*`, "#")

	chk.Str(
		statfs.Status(),
		"Total: # Avail: # (#.#%) INodes: # Avail: # (#.#%)",
	)

	_ = chk.CreateTmpFile([]byte("string"))

	chk.StrSlice(
		strings.Split(statfs.Delta(), "\n"),
		[]string{
			"Before: Total: # Avail: # (#.#%) INodes: # Avail: # (#.#%)",
			" After: Total: # Avail: # (#.#%) INodes: # Avail: # (#.#%)",
			"Deltas: Bytes: # (#.#%) INodes: # (#.#%)",
		},
	)
}
