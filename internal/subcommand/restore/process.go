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

package restore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/directory"
	"github.com/dancsecs/szbck/internal/rsync"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/target"
)

var reFindBackupSubDir = regexp.MustCompile(
	`` +

		`2\d\d\d[0-1]\d[0-3]\d` + // Year.
		`_` +
		`[0-2]\d[0-5]\d[0-5]\d` + // Time.
		`\.` +
		`\d\d\d\d` + // Partial second.
		`\.szb`,
)

func parseArgs(args []string) (*settings.Config, string, bool, bool, error) {
	var (
		dryRun   bool
		keep     bool
		snapshot string
		cfg      *settings.Config
		err      error
	)

	dryRun, args, err = szargs.Arg("--dry-run").Is(args)

	if err == nil {
		keep, args, err = szargs.Arg("--keep").Is(args)
	}

	if err == nil {
		snapshot, _, args, err = szargs.Arg("-s").Value(args)
	}

	if err == nil {
		cfg, err = settings.LoadFromArgs(args)
	}

	return cfg, snapshot, dryRun, keep, err //nolint:wrapcheck // Ok.
}

// MakeDirs creates the target string based on the restoreFrom directory
// past the required szerszam backup directory name.
func MakeDirs(srcPath, toPath string) (string, string, error) {
	var (
		fromDir string
		toDir   string
		err     error
	)

	// Split and clean up the target into the directory and basename cleared
	// of trailing path separators.
	toDir, toBase := filepath.Split(
		strings.TrimRight(toPath, directory.PathSeparator),
	)
	toDir = strings.TrimRight(toDir, directory.PathSeparator)

	sRoot, sPath, err := target.Split(srcPath, reFindBackupSubDir)

	if errors.Is(err, target.ErrSplitNotFound) {
		// try to append latest.
		sRoot, sPath, err = target.Split(
			filepath.Join(srcPath, target.LatestDirectoryLink),
			reFindBackupSubDir,
		)
	}

	if err == nil { //nolint:nestif  // Ok.
		if sPath == "" || sPath == toBase {
			fromDir = filepath.Join(sRoot, toBase)
		} else {
			if !strings.HasPrefix(sPath, toBase+directory.PathSeparator) {
				err = fmt.Errorf(
					"%w: '%s' must start with '%s'",
					ErrInvalidSrcPath,
					sPath,
					toBase,
				)
			} else {
				fromDir = filepath.Join(sRoot, sPath)
			}
		}
	}

	if err == nil {
		toDir = filepath.Join(toDir, filepath.Dir(sPath))

		return fromDir, toDir, nil
	}

	return "", "", err
}

// Process parses the remaining arguments restoring from a szbackup snapshot.
func Process(args []string) (string, error) {
	var (
		cfg         *settings.Config
		dryRun      bool
		keep        bool
		snapshot    string
		restoreFrom string
		restoreTo   string
		err         error
	)

	cfg, snapshot, dryRun, keep, err = parseArgs(args)

	if err == nil {
		restoreFrom, restoreTo, err = MakeDirs(
			filepath.Join(cfg.Target.GetPath(), snapshot),
			cfg.Source,
		)
	}

	if err == nil {
		err = rsync.Run(
			rsync.BuildArgs(
				!keep, // Delete from target unless keep option was provided.
				dryRun,
				"", // no linkDesk for restore operations.
				cfg.Options,
				cfg.RestoreOptions,
				restoreFrom,
				restoreTo,
			),
			os.Stdout,
			os.Stderr,
		)
	}

	if err == nil {
		return "restore successful", nil
	}

	return "", fmt.Errorf("%w: %w", ErrRestoreError, err)
}
