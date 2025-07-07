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
	"regexp"
	"strings"
)

// Parsing errors.
var (
	ErrInvalid       = errors.New("invalid config file")
	ErrConfigLine    = errors.New("line")
	ErrInvalidSyntax = errors.New("invalid key:value pair")
	ErrUnknownKey    = errors.New("unknown key")
	ErrUndefined     = errors.New("undefined")
	ErrMissing       = errors.New("missing")
	ErrSyntax        = errors.New("syntax")
	ErrRange         = errors.New("range")
	ErrDuplicate     = errors.New("duplicate")
)

func isVerbose(flag string) bool {
	return flag == "-v" || flag == "--verbose"
}

// Parse takes a the content of a configuration file, and returns a Config
// structure if there are no errors.
func Parse(txt string) (*Config, error) {
	var (
		cfg   Config
		line  string
		key   string
		value string
		found bool
		err   error
	)

	reStripComments := regexp.MustCompile(`\s*\#.*$`)

	for lineNbr, rawLine := range strings.Split(txt, "\n") {
		line = reStripComments.ReplaceAllString(rawLine, "")
		if line == "" {
			continue // skip blank and fully commented lines.
		}

		key, value, found = strings.Cut(line, ":")
		if !found {
			err = ErrInvalidSyntax
		} else {
			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			err = cfg.validateKeyValue(key, value)
		}

		if err != nil {
			err = fmt.Errorf(
				"%w(%d): %w\n\t%s", ErrConfigLine, lineNbr+1, err, rawLine,
			)

			break
		}
	}

	if err == nil {
		err = cfg.validateMandatory()
	}

	if err == nil {
		return &cfg, nil
	}

	return nil, fmt.Errorf("%w: %w", ErrInvalid, err)
}
