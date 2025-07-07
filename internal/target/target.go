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

package target

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/dancsecs/szbck/internal/directory"
)

const (
	// LatestDirectoryLink names the link pointing to the latest backup set.
	LatestDirectoryLink = "latest"
	// BackupDirectoryFormat specifies how backup directories are named.
	BackupDirectoryFormat = "20060102_150405.0000"
	// BackupDirectoryExtension identifies a directory as a Szerszam backup
	// snapshot.
	BackupDirectoryExtension = ".szb"
)

// Path represent the directory containing the szerszam backup.
type Path struct {
	path string
}

// New return a new validated target path.
func New(path string) (*Path, error) {
	newPath := &Path{
		path: path,
	}

	err := newPath.Validate()

	if err == nil {
		return newPath, nil
	}

	return nil, fmt.Errorf("%w: %w", ErrNew, err)
}

// GetPath returns the target's path.
func (target Path) GetPath() string {
	return target.path
}

// Latest returns the path to the latest backup.
func (target Path) Latest() string {
	return filepath.Join(target.path, LatestDirectoryLink)
}

// HasLatest returns true if the backup has a latest symbolic link set.
func (target Path) HasLatest() (bool, error) {
	lstat, err := os.Lstat(target.Latest())
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	if err == nil && lstat.Mode().Type()&os.ModeSymlink != 0 {
		return true, nil
	}

	if err == nil {
		err = ErrInvalidLatest
	}

	return false, fmt.Errorf("%w: %w", ErrHasLatest, err)
}

// Validate insures that the target is a directory that either has a 'latest'
// symlink pointing to a backup set or the directory is empty.
func (target Path) Validate() error {
	var (
		hasLatest bool
		err       error
	)

	err = directory.Is(target.path)

	if err == nil {
		hasLatest, err = target.HasLatest()
	}

	if err == nil && !hasLatest {
		err = directory.IsEmpty(target.path)
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", ErrInvalid, err)
}

// Create a new target directory based on the provided date/time.
func (target Path) Create(tme time.Time, perm os.FileMode) (string, error) {
	var (
		newDir string
		err    error
	)

	err = target.Validate()

	if err == nil {
		newDir = filepath.Join(
			target.path,
			tme.Format(BackupDirectoryFormat)+BackupDirectoryExtension,
		)

		_, err = os.Stat(newDir)
		if err == nil {
			err = fmt.Errorf("%w: '%s'", ErrCreateAlreadyExists, newDir)
		} else if errors.Is(err, os.ErrNotExist) {
			err = nil
		}
	}

	if err == nil {
		err = os.MkdirAll(newDir, perm)
	}

	if err == nil {
		return newDir, nil
	}

	return "", fmt.Errorf("%w: %w", ErrCreateTargetFailed, err)
}

// SetLatest create a symbolic link to the supplied backup directory.
func (target Path) SetLatest(path string) error {
	err := directory.LinkRelative(path, target.Latest())

	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", ErrInvalidLatest, err)
}

// Split creates the target string based on the restoreFrom
// directory past the required szerszam backup directory name.
func Split(dir string, reSplit *regexp.Regexp) (string, string, error) {
	var pathComponents []int

	absDir, err := filepath.EvalSymlinks(dir)

	if err == nil {
		absDir, err = filepath.Abs(absDir)
	}

	if err == nil {
		pathComponents = reSplit.FindStringIndex(absDir)

		if pathComponents == nil {
			err = ErrSplitNotFound
		}
	}

	if err == nil {
		return absDir[:pathComponents[1]],
			strings.Trim(absDir[pathComponents[1]:], directory.PathSeparator),
			nil
	}

	return "", "", fmt.Errorf("%w: %w", ErrInvalidSplit, err)
}
