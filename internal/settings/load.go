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

	"github.com/dancsecs/szargs"
)

// Load errors.
var (
	ErrLoad     = errors.New("load error")
	ErrNoTarget = errors.New("no target configured or overridden")
)

// Load reads and validates the config file in the named directory.
func Load(fPath string) (*Config, error) {
	var (
		cfg      *Config
		fileData []byte
		err      error
	)

	//nolint:gosec // Ok.
	fileData, err = os.ReadFile(fPath)
	if err == nil {
		cfg, err = Parse(string(fileData))
	}

	if err == nil {
		return cfg, nil
	}

	return nil, fmt.Errorf("%w: %w", ErrLoad, err)
}

// LoadFromArgs loads the specified configuration file present as the last
// argument with the backup target optionally overridden by a "-t" filename
// argument. An error if the replacement is invalid or if both the replacement
// and the configured target are not defined.
func LoadFromArgs(args *szargs.Args) (*Config, error) {
	var (
		trgOverride string
		cfgFilename string
		cfg         *Config
		err         error
	)

	trgOverride, _ = args.ValueString("-t", "")
	cfgFilename = args.NextString("backup config filename", "")
	args.Done()
	err = args.Err()

	if err == nil {
		cfg, err = Load(cfgFilename)
	}

	if err == nil && cfg.Target == nil && trgOverride == "" {
		err = ErrNoTarget
	}

	if err == nil && trgOverride != "" {
		cfg.Target = nil
		err = cfg.validateTarget(trgOverride)
	}

	if err == nil {
		return cfg, nil
	}

	return nil, err
}
