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

package create

import (
	"errors"
	"fmt"
	"os"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/settings"
)

// // WritePermWarningMsg informs that write access should not be granted.
// const WritePermWarningMsg = "WARNING:  " +
// 	"Backup permissions should not permit writes!\n"

// func writePermissionWarning(perm os.FileMode) string {
// 	const (
// 		usrWrite = 0o0200
// 		grpWrite = 0o0020
// 		othWrite = 0o0002
// 	)

// 	if perm&usrWrite+perm&grpWrite+perm&othWrite > 0 {
// 		return WritePermWarningMsg
// 	}

// 	return ""
// }

func parseArgs(args *szargs.Args) (string, string, string, error) {
	var (
		outFile string
		source  string
		trg     string
	)

	outFile, _ = args.ValueString("-o", "")
	trg, _ = args.ValueString("-t", "")
	source = args.NextString("source directory", "")
	args.Done()

	return outFile, source, trg, args.Err() //nolint:wrapcheck // Ok.
}

func writeFile(outFile, cfg string) error {
	const defaultPerm = 0o0666

	_, err := os.Stat(outFile)
	if err == nil {
		err = fmt.Errorf("%w: '%s'", ErrOutFileExists, outFile)
	} else if errors.Is(err, os.ErrNotExist) {
		err = nil
	}

	if err == nil {
		err = os.WriteFile(outFile, []byte(cfg), defaultPerm)
	}

	return err //nolint:wrapcheck // Ok.
}

// Process parses the remaining arguments creating a szbackup configuration
// file.
func Process(args *szargs.Args) (string, error) {
	var (
		cfgTxt  string
		source  string
		trg     string
		outFile string
		err     error
	)

	outFile, source, trg, err = parseArgs(args)

	if err == nil {
		cfgTxt, err = settings.Create(source, trg)
	}

	if err == nil && outFile == "" {
		return cfgTxt, nil
	}

	if err == nil {
		err = writeFile(outFile, cfgTxt)
	}

	if err == nil {
		return "successfully created: " + outFile, nil
	}

	return "", fmt.Errorf("%w: %w", ErrInvalid, err)
}
