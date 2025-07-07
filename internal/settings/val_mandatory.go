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

// Mandatory errors.
var (
	ErrSourceMissing     = errors.New("missing source")
	ErrPermissionMissing = errors.New("missing permission")
	ErrKeepHourlyMissing = errors.New("missing keep hourly retention")
	ErrKeepDailyMissing  = errors.New("missing keep daily retention")
	ErrNoSnapshotOptions = errors.New("no snapshot options defined")
	ErrNoRestoreOptions  = errors.New("no restore options defined")
)

func (cfg *Config) validateMandatory() error {
	var err error

	addError := func(add bool, newErr error) {
		if add {
			if err == nil {
				err = newErr
			} else {
				err = fmt.Errorf("%w: %w", err, newErr)
			}
		}
	}

	addError(cfg.Source == "", ErrSourceMissing)
	addError(cfg.Permission == 0, ErrPermissionMissing)
	addError(cfg.KeepHourly == 0, ErrKeepHourlyMissing)
	addError(cfg.KeepDaily == 0, ErrKeepDailyMissing)
	addError(cfg.KeepDaily <= cfg.KeepHourly, ErrRetentionDailyMin)

	if len(cfg.Options) == 0 {
		addError(len(cfg.SnapshotOptions) == 0, ErrNoSnapshotOptions)
		addError(len(cfg.RestoreOptions) == 0, ErrNoRestoreOptions)
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", ErrUndefined, err)
}
