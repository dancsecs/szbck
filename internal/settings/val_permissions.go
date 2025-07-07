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
	"os"
	"strings"
)

// Permission errors.
var (
	ErrPermission         = errors.New("permission error")
	ErrPermissionSymbolic = errors.New("invalid symbolic")
	ErrPermissionRange    = errors.New("range error")
)

func (cfg *Config) validatePermission(perm string) error {
	var (
		validPerm uint32
		err       error
	)

	if cfg.Permission != 0 {
		err = ErrDuplicate
	}

	if err == nil && perm == "" {
		err = ErrMissing
	}

	if err == nil {
		if strings.HasPrefix(perm, "0o") {
			validPerm, err = validateOctalPermission(perm)
		} else {
			validPerm, err = validateSymbolicPermission(perm)
		}
	}

	if err == nil {
		cfg.Permission = os.FileMode(validPerm)

		return nil
	}

	return fmt.Errorf("%w: %w", ErrPermission, err)
}
