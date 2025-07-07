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
	"strings"

	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szbck/internal/target"
)

// Create errors.
var (
	ErrCreate = errors.New("create settings error")
)

// Create creates and validates a new Backup configuration file with the
// absolute paths of the provides source and optional target.
func Create(src, trg string) (string, error) {
	var (
		absSrc  string
		absTrg  string
		cfgFile string
		err     error
	)

	absSrc, err = filepath.Abs(src)

	if err == nil {
		err = directory.Is(absSrc)
	}

	if err == nil && trg != "" {
		absTrg, err = filepath.Abs(trg)
	}

	if err == nil && absTrg != "" {
		_, err = target.New(absTrg)
	}

	if err == nil {
		cfgFile = strings.Replace(
			DefaultConfig,
			"source: /home/user",
			"source: "+absSrc,
			1,
		)

		if absTrg != "" {
			cfgFile = strings.Replace(
				cfgFile,
				"#target: /mnt/backupDir",
				"target: "+absTrg,
				1,
			)
		}

		_, err = Parse(cfgFile)
	}

	if err == nil {
		return cfgFile, nil
	}

	return "", fmt.Errorf("%w: %w", ErrCreate, err)
}
