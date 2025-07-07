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

package vet

import (
	"fmt"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/settings"
)

func parseArguments(args []string) error {
	var (
		configFileName string
		err            error
	)

	configFileName, err = szargs.Last("backup config filename", args)

	if err == nil {
		_, err = settings.Load(configFileName)
	}

	return err //nolint:wrapcheck // Ok.
}

// Process parses the remaining arguments deleting previous backups.
func Process(args []string) (string, error) {
	err := parseArguments(args)

	if err == nil {
		return "vet successful (no problems found)\n", nil
	}

	return "", fmt.Errorf("%w: %w", ErrVetError, err)
}
