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
)

// RestoreOption  errors.
var (
	ErrRestoreOption = errors.New("invalid restore option")
)

func (cfg *Config) validateRestoreOption(value string) error {
	var err error

	for i, mi := 0, len(cfg.RestoreOptions); i < mi && err == nil; i++ {
		if cfg.RestoreOptions[i] == value && !isVerbose(value) {
			err = ErrDuplicate
		}
	}

	if err == nil && value == "" {
		err = ErrMissing
	}

	if err == nil {
		cfg.RestoreOptions = append(cfg.RestoreOptions, value)

		return nil
	}

	return fmt.Errorf("%w: %w", ErrRestoreOption, err)
}
