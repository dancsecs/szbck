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

package prune

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/du"
	"github.com/dancsecs/szbck/internal/out"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/target"
)

const permToDelete = 0o0700

func parseArguments(args []string) (*settings.Config, string, string, error) {
	var (
		isDryRun bool
		dryRun   string
		numToDel string
		found    bool
		cfg      *settings.Config
		err      error
	)

	isDryRun, args, err = szargs.Arg("--dry-run").Is(args)
	if err == nil {
		if isDryRun {
			dryRun = " (DRY RUN)"
		}
	}

	if err == nil {
		numToDel, found, args, err = szargs.Arg("-n").Value(args)
		if err == nil && !found {
			numToDel = "1"
		}
	}

	if err == nil {
		cfg, err = settings.LoadFromArgs(args)
	}

	return cfg, dryRun, numToDel, err //nolint:wrapcheck // Ok.
}

func loadBackupDirs(trg string) ([]string, error) {
	matchingDirs, err := filepath.Glob(
		filepath.Join(trg, "*"+target.BackupDirectoryExtension),
	)

	if err == nil {
		switch len(matchingDirs) {
		case 0:
			err = ErrNoBackups
		case 1:
			err = ErrOnlyLatest
		default:
			// sort list and remove the newest
			slices.Sort(matchingDirs)
			matchingDirs = matchingDirs[:len(matchingDirs)-1]
		}
	}

	return matchingDirs, err
}

func validateNumberToDelete(rawNum string, maxNum int) (int, error) {
	var (
		tmpNum   int64
		numToDel int
		err      error
	)

	if rawNum == "all" {
		numToDel = maxNum
	} else {
		tmpNum, err = strconv.ParseInt(rawNum, 10, 0)
		if err != nil || tmpNum < 1 {
			err = fmt.Errorf("%w: '%s'", ErrInvalidNum, rawNum)
		}

		if err == nil {
			numToDel = min(int(tmpNum), maxNum)
		}
	}

	return numToDel, err
}

func pruneDirectories(num int, dirs []string, dryRun string) error {
	var err error

	if num == 1 {
		out.Print("Purging oldest backup" + dryRun + "\n\n")
	} else {
		out.Printf("Purging %d oldest backups"+dryRun+"\n\n", num)
	}

	for i := 0; i < num && err == nil; i++ {
		out.Print("Purging backup: " + dirs[i] + "\n")

		if dryRun == "" {
			err = os.Chmod(dirs[i], permToDelete)
			if err == nil {
				err = os.RemoveAll(dirs[i])
			}
		}
	}

	return err //nolint:wrapcheck // ok.
}

// Process parses the remaining arguments deleting previous backups.
func Process(args []string) (string, error) {
	var (
		rawNumToDel  string
		dryRun       string
		numToDel     int
		cfg          *settings.Config
		matchingDirs []string
		beforeBytes  int64
		afterBytes   int64
		err          error
	)

	cfg, dryRun, rawNumToDel, err = parseArguments(args)

	if err == nil {
		matchingDirs, err = loadBackupDirs(cfg.Target.GetPath())
	}

	if err == nil {
		beforeBytes, err = du.Total(cfg.Target.GetPath())
	}

	if err == nil {
		numToDel, err = validateNumberToDelete(
			rawNumToDel,
			len(matchingDirs),
		)
	}

	if err == nil {
		err = pruneDirectories(numToDel, matchingDirs, dryRun)
	}

	if err == nil {
		afterBytes, err = du.Total(cfg.Target.GetPath())
	}

	if err == nil {
		out.Printf(
			"purge successful%s\n"+
				"Before: %s After: %s Total Recovered: %s bytes\n",
			dryRun,
			out.Int(beforeBytes),
			out.Int(afterBytes),
			out.Int(beforeBytes-afterBytes),
		)

		return "", nil
	}

	return "", fmt.Errorf("%w: %w", ErrPruneError, err)
}
