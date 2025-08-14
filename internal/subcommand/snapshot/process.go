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

package snapshot

import (
	"fmt"
	"os"
	"time"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/du"
	"github.com/dancsecs/szbck/internal/out"
	"github.com/dancsecs/szbck/internal/rsync"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/subcommand/trim"
)

const initialBackupDirPerm = 0o0700

// NextHourIn returns the time to wait for the start of the next hour.  If
// elapsed is >= 31 minutes than an extra hour is added to the delay.
func NextHourIn(elapsed time.Duration) time.Duration {
	timeToNextTime := time.Hour - elapsed

	for timeToNextTime < time.Minute*31 {
		timeToNextTime += time.Hour
	}

	return timeToNextTime
}

func parseArgs(args []string) (*settings.Config, string, bool, bool, error) {
	var (
		cfg       *settings.Config
		isDryRun  bool
		dryRun    string
		trimAfter bool
		daemon    bool
		err       error
	)

	isDryRun, args, err = szargs.Arg("--dry-run").Is(args)

	if err == nil && isDryRun {
		dryRun = " (DRY RUN)"
	}

	if err == nil {
		trimAfter, args, err = szargs.Arg("--trim").Is(args)
	}

	if err == nil {
		daemon, args, err = szargs.Arg("--daemon").Is(args)
	}

	if err == nil {
		cfg, err = settings.LoadFromArgs(args)
	}

	return cfg, dryRun, trimAfter, daemon, err //nolint:wrapcheck // Ok.
}

func run(dryRun bool, linkDest, newDir string, cfg *settings.Config) error {
	return rsync.Run( //nolint:wrapcheck // Ok.
		rsync.BuildArgs(
			true, // Delete from target
			dryRun,
			linkDest,
			cfg.Options,
			cfg.SnapshotOptions,
			cfg.Source,
			newDir,
		),
		os.Stdout,
		os.Stderr,
	)
}

// Process parses the remaining arguments creating a szbackup snapshot.
//
//nolint:cyclop,funlen // Ok.
func Process(args []string) (string, error) {
	var (
		cfg            *settings.Config
		dryRunMsg      string
		trimAfter      bool
		daemon         bool
		purgedCount    int
		purgedMsg      string
		totalPurged    int
		totalPurgedMsg string
		hasLatest      bool
		linkDest       string
		newDir         string
		beforeBytes    int64
		afterBytes     int64
		err            error
		startTime      time.Time
	)

	cfg, dryRunMsg, trimAfter, daemon, err = parseArgs(args)

	if err == nil {
		beforeBytes, err = du.Total(cfg.Target.GetPath())
	}

	runOnce := true
	sleepBetweenRuns := time.Nanosecond

	for (runOnce || daemon) && err == nil {
		time.Sleep(sleepBetweenRuns)
		startTime = time.Now()
		runOnce = false

		newDir, err = cfg.Target.Create(time.Now(), initialBackupDirPerm)

		if err == nil {
			hasLatest, err = cfg.Target.HasLatest()
			if err == nil && hasLatest {
				linkDest = cfg.Target.Latest()
			}
		}

		if err == nil {
			err = run(dryRunMsg != "", linkDest, newDir, cfg)
		}

		if err == nil && dryRunMsg == "" {
			err = os.Chmod(newDir, cfg.Permission)
		}

		if err == nil && dryRunMsg == "" {
			err = cfg.Target.SetLatest(newDir)
		}

		if err == nil && dryRunMsg != "" {
			err = os.RemoveAll(newDir)
		}

		if err == nil && trimAfter {
			purgedCount, err = trim.PurgeSnapshots(cfg, time.Now(), dryRunMsg)
			purgedMsg = " (Purged: " + out.Int(int64(purgedCount)) + ")"
			totalPurged += purgedCount
			totalPurgedMsg = " (Total Purged: " +
				out.Int(int64(totalPurged)) + ")"
		}

		if err == nil {
			afterBytes, err = du.Total(cfg.Target.GetPath())
		}

		if err == nil {
			out.Printf(
				"snapshot successful%s%s\n"+
					"Before: %s After: %s Used: %s bytes\n",
				purgedMsg,
				dryRunMsg,
				out.Int(beforeBytes),
				out.Int(afterBytes),
				out.Int(afterBytes-beforeBytes),
			)

			beforeBytes = afterBytes
		}

		sleepBetweenRuns = NextHourIn(time.Since(startTime))
	}

	if err == nil {
		return "", nil
	}

	return "", fmt.Errorf("%w%s: %w", ErrSnapshotError, totalPurgedMsg, err)
}
