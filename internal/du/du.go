/*
   Golang rsync backup utility wrapper: szbck.
   Copyright (C) 2025-2026 Leslie Dancsecs

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

package du

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dancsecs/szbck/internal/directory"
)

// Total returns the total number of bytes used by the directory tree.
func Total(dir string) (int64, error) {
	const (
		base = 10
		bits = 64
	)

	var (
		tmpStr string
		total  int64
		err    error
	)

	err = directory.Is(dir)

	if err == nil {
		tmpStr, err = Run([]string{"-s", "-b", dir}, os.Stderr)
	}

	if err == nil {
		total, err = strconv.ParseInt(
			strings.Split(tmpStr, "\t")[0],
			base,
			bits,
		)
	}

	if err == nil {
		return total, nil
	}

	return 0, fmt.Errorf("%w: %w", ErrInvalid, err)
}

// Totals returns the total number of bytes used by the directory trees with
// dir2 size accounting for hard links to dir1.
func Totals(dir1, dir2 string) (int64, int64, error) {
	const (
		base = 10
		bits = 64
	)

	var (
		tmpStr  string
		results []string
		total1  int64
		total2  int64
		err     error
	)

	err = directory.Is(dir1)

	if err == nil {
		err = directory.Is(dir2)
	}

	if err == nil {
		tmpStr, err = Run([]string{"-s", "-b", dir1, dir2}, os.Stderr)
	}

	if err == nil {
		results = strings.Split(tmpStr, "\n")
	}

	if err == nil {
		total1, err = strconv.ParseInt(
			strings.Split(results[0], "\t")[0],
			base,
			bits,
		)
	}

	if err == nil {
		total2, err = strconv.ParseInt(
			strings.Split(results[1], "\t")[0],
			base,
			bits,
		)
	}

	if err == nil {
		return total1, total2, nil
	}

	return 0, 0, fmt.Errorf("%w: %w", ErrInvalid, err)
}
