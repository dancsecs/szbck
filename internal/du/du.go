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
