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

package directory

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// PathSeparator provides for a single project reference for a string casted
// instance of the operating system's path separator.
const PathSeparator = string(os.PathSeparator)

// Is confirms the named directory exists.
func Is(dir string) error {
	var (
		stat os.FileInfo
		err  error
	)

	stat, err = os.Stat(dir)
	if err != nil {
		err = ErrInvalid
	}

	if err == nil && !stat.IsDir() {
		err = ErrNotADirectory
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: '%s'", err, dir)
}

// IsEmpty confirms that a "new" target directory is empty.
func IsEmpty(dir string) error {
	var (
		items *os.File
		err   error
	)

	items, err = os.Open(dir) //nolint:gosec // Ok.

	if err == nil {
		defer func() {
			_ = items.Close()
		}()

		_, err = items.Readdirnames(1)
	}

	if err == nil {
		err = ErrNewNotEmpty
	} else if errors.Is(err, io.EOF) {
		return nil
	}

	return err
}

// LinkRelative create the provided symbolic link path to the supplied
// directory path.
func LinkRelative(fromDir, toLink string) error {
	err := Is(fromDir)

	if err == nil {
		err = os.Remove(toLink)
		if errors.Is(err, os.ErrNotExist) {
			err = nil
		}
	}

	if err == nil {
		err = os.Symlink(filepath.Base(fromDir), toLink)
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf(
		"%w: (from: '%s' to: '%s'): %w",
		ErrCreateLink,
		fromDir,
		toLink,
		err,
	)
}
