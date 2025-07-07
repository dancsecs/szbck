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

package status

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/dancsecs/szbck/internal/du"
	"github.com/dancsecs/szbck/internal/out"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/target"
)

func parseArguments(args []string) (*settings.Config, error) {
	var (
		cfg *settings.Config
		err error
	)

	cfg, err = settings.LoadFromArgs(args)

	return cfg, err //nolint:wrapcheck // Ok.
}

func loadBackupDirs(trg string) ([]string, error) {
	matchingDirs, err := filepath.Glob(
		filepath.Join(trg, "*"+target.BackupDirectoryExtension),
	)

	if err == nil {
		if len(matchingDirs) == 0 {
			err = ErrNoBackups
		} else {
			// sort list and remove the newest
			slices.Sort(matchingDirs)
		}
	}

	if err == nil {
		return matchingDirs, nil
	}

	return nil, err
}

func buildReport(trg string) (string, error) {
	var (
		outStr    string
		dirs      []string
		totalSize int64
		dirSize   int64
		err       error
	)

	dirs, err = loadBackupDirs(trg)

	if err == nil {
		totalSize, err = du.Total(trg)
	}

	if err == nil {
		outStr = fmt.Sprintf(
			"Backup Sets: %s\n"+
				"Total Bytes: %s\n\n",
			out.Int(int64(len(dirs))),
			out.Int(totalSize),
		)
	}

	for i, mi := 0, len(dirs); i < mi && err == nil; i++ {
		dirSize, err = du.Total(dirs[i])
		if err == nil {
			outStr += filepath.Base(dirs[i]) + ": " + out.Int(dirSize) + "\n"
		}
	}

	if err == nil {
		return outStr, nil
	}

	return "", fmt.Errorf("%w: %w", ErrReportFailed, err)
}

// Process parses the remaining arguments deleting previous backups.
func Process(args []string) (string, error) {
	var (
		cfg    *settings.Config
		report string
		err    error
	)

	cfg, err = parseArguments(args)

	if err == nil {
		report, err = buildReport(cfg.Target.GetPath())
	}

	if err == nil {
		return "status successful\n\n" + report, nil
	}

	return "", fmt.Errorf("%w: %w", ErrStatusError, err)
}
