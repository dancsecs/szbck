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
	"fmt"
)

func (cfg *Config) validateKeyValue(key, value string) error {
	switch key {
	case "source":
		return cfg.validateSource(value)
	case "target":
		return cfg.validateTarget(value)
	case "permission":
		return cfg.validatePermission(value)
	case "option":
		return cfg.validateOption(value)
	case "snapshotOption":
		return cfg.validateSnapshotOption(value)
	case "restoreOption":
		return cfg.validateRestoreOption(value)
	case "keepHourly":
		return cfg.validateKeepHourly(value)
	case "keepDaily":
		return cfg.validateKeepDaily(value)
	default:
		return fmt.Errorf("%w: '%s'", ErrUnknownKey, key)
	}
}
