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
)

// Source errors.
var (
	ErrSource    = errors.New("invalid source")
	ErrSourceAbs = errors.New("source must be absolute path")
)

func (cfg *Config) validateSource(source string) error {
	var (
		absSrc string
		err    error
	)

	if cfg.Source != "" {
		err = ErrDuplicate
	}

	if err == nil {
		err = directory.Is(source)
	}

	if err == nil {
		absSrc, err = filepath.Abs(source)
	}

	if err == nil {
		cfg.Source = absSrc

		return nil
	}

	return fmt.Errorf("%w: %w", ErrSource, err)
}
