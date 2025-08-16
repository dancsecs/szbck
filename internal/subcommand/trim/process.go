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

package trim

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/du"
	"github.com/dancsecs/szbck/internal/out"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/target"
)

const permToDelete = 0o0700

func parseArguments(args *szargs.Args) (*settings.Config, string, error) {
	var (
		isDryRun bool
		dryRun   string
		cfg      *settings.Config
		err      error
	)

	isDryRun = args.Is("--dry-run", "")
	if isDryRun {
		dryRun = " (DRY RUN)"
	}

	err = args.Err()

	if err == nil {
		cfg, err = settings.LoadFromArgs(args)
	}

	return cfg, dryRun, err //nolint:wrapcheck // Ok.
}

func getTimestamp(fName string) (time.Time, error) {
	fileTime, err := time.ParseInLocation(
		target.BackupDirectoryFormat,
		strings.TrimSuffix(
			filepath.Base(fName),
			target.BackupDirectoryExtension,
		),
		time.Local, //nolint:gosmopolitan // Time to match local filesystem.
	)
	if err == nil {
		return fileTime, nil
	}

	return time.Time{},
		fmt.Errorf("%w: '%s'", ErrInvalidSnapshotName, fName)
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
		}
	}

	return matchingDirs, err
}

// PurgeSnapshots removes snapshots based on the configured retention policy
// using the provide time as the root to base snapshot expiring on.  The most
// recent snapshot pointed to by the "latest" symbolic link is never deleted.
func PurgeSnapshots(
	cfg *settings.Config,
	tme time.Time, // The reference timestamp to base trim functions on.
	dryRun string,
) (int, error) {
	var (
		hourlyCutoff time.Time
		dailyCutoff  time.Time
		tms          []time.Time
		dirs         []string
		remove       []bool
		purgedCount  int
		err          error
	)

	hourlyCutoff = tme.Add(-cfg.KeepHourly)
	dailyCutoff = tme.Add((-cfg.KeepDaily))

	dirs, err = loadBackupDirs(cfg.Target.GetPath())

	if err != nil || len(dirs) < 2 {
		if len(dirs) == 0 {
			err = ErrNoBackups
		} else if len(dirs) == 1 {
			err = ErrOnlyLatest
		}
	}

	// Convert dir names to real timestamps all at once.
	if err == nil {
		tms = make([]time.Time, len(dirs))
		for i, mi := 0, len(dirs); i < mi && err == nil; i++ {
			tms[i], err = getTimestamp(dirs[i])
		}
	}

	if err == nil {
		remove = identifyRemovals(tms, hourlyCutoff, dailyCutoff)
	}

	if err == nil {
		purgedCount, err = processPurge(dirs, tms, remove, dryRun)
	}

	return purgedCount, err
}

// Process parses the remaining arguments deleting previous backups.
func Process(args *szargs.Args) (string, error) {
	var (
		dryRun      string
		cfg         *settings.Config
		purgedCount int
		beforeBytes int64
		afterBytes  int64
		err         error
	)

	cfg, dryRun, err = parseArguments(args)

	if err == nil {
		beforeBytes, err = du.Total(cfg.Target.GetPath())
	}

	if err == nil {
		purgedCount, err = PurgeSnapshots(cfg, time.Now(), dryRun)
	}

	if err == nil {
		afterBytes, err = du.Total(cfg.Target.GetPath())
	}

	if err == nil {
		out.Printf(
			"trim successful (Purged: %d)%s\n"+
				"Before: %s After: %s Total Recovered: %s bytes\n",
			purgedCount,
			dryRun,
			out.Int(beforeBytes),
			out.Int(afterBytes),
			out.Int(beforeBytes-afterBytes),
		)

		return "", nil
	}

	return "", fmt.Errorf(
		"%w (Purged: %d): %w", ErrTrimError, purgedCount, err,
	)
}
