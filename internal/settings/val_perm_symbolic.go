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
	"strings"
)

// Symbolic permission errors.
var (
	ErrSymbolicNone = errors.New("no permissions")
	ErrSymbolic     = errors.New("symbolic")
)

func validateSymbolicGroup(expGroup, perm string) (uint32, error) {
	//nolint:mnd // Ok.
	if strings.HasPrefix(perm, expGroup) {
		switch perm[2:] {
		case "rwx":
			return 0o7, nil
		case "rw":
			return 0o6, nil
		case "rx":
			return 0o5, nil
		case "r":
			return 0o4, nil
		case "wx":
			return 0o3, nil
		case "w":
			return 0o2, nil
		case "x":
			return 0o1, nil
		case "-":
			return 0o0, nil
		}
	}

	return 0, ErrSyntax
}

func validateSymbolicGroups(perms []string) (uint32, error) {
	var (
		gPerm uint32
		value uint32
		err   error
	)

	gPerm, err = validateSymbolicGroup("u", perms[0])

	if err == nil {
		value |= gPerm << 6 //nolint:mnd // Ok Shift to user bit position.
		gPerm, err = validateSymbolicGroup("g", perms[1])
	}

	if err == nil {
		value |= gPerm << 3 //nolint:mnd // Ok Shift to group bit position.
		gPerm, err = validateSymbolicGroup("o", perms[2])
	}

	if err == nil {
		value |= gPerm

		return value, nil
	}

	return 0, err
}

func validateSymbolicPermission(perm string) (uint32, error) {
	const (
		expectedPermGroups = 3
	)

	var (
		value uint32
		perms []string
		err   error
	)

	if perm == "" {
		err = ErrSymbolicNone
	}

	if err == nil {
		perms = strings.Split(perm, ";")

		if len(perms) != expectedPermGroups {
			err = ErrSyntax
		}
	}

	if err == nil {
		value, err = validateSymbolicGroups(perms)
	}

	if err == nil && value < 1 {
		err = ErrSymbolicNone
	}

	if err == nil {
		return value, nil
	}

	return 0, fmt.Errorf("%w: %w", ErrSymbolic, err)
}
