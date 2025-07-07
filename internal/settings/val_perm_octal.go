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

package settings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Octal permission errors.
var (
	ErrOctal = errors.New("octal")
)

func validateOctalPermission(perm string) (uint32, error) {
	const (
		base = 8
		bits = 32
	)

	var (
		value uint64
		err   error
	)

	if !strings.HasPrefix(perm, "0o0") {
		err = ErrSyntax
	}

	if err == nil {
		value, err = strconv.ParseUint(perm[3:], base, bits)
	}

	if errors.Is(err, strconv.ErrSyntax) {
		err = ErrSyntax
	}

	if errors.Is(err, strconv.ErrRange) {
		err = ErrRange
	}

	if err == nil && (value < 1 || value > 0o0777) {
		err = ErrRange
	}

	if err == nil {
		return uint32(value), nil //nolint:gosec // Ok range already checked.
	}

	return 0, fmt.Errorf("%w: %w", ErrOctal, err)
}
