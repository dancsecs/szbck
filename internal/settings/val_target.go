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
	"path/filepath"

	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szbck/internal/target"
)

// Target errors.
var (
	ErrTarget    = errors.New("invalid target")
	ErrTargetAbs = errors.New("source must be absolute path")
)

func (cfg *Config) validateTarget(trg string) error {
	var (
		newTrg *target.Path
		absTrg string
		err    error
	)

	if cfg.Target != nil {
		err = ErrDuplicate
	}

	if err == nil {
		err = directory.Is(trg)
	}

	if err == nil {
		absTrg, err = filepath.Abs(trg)
	}

	if err == nil {
		newTrg, err = target.New(absTrg)
	}

	if err == nil {
		cfg.Target = newTrg

		return nil
	}

	return fmt.Errorf("%w: %w", ErrTarget, err)
}
