/*
   Golang rsync backup utility wrapper: szbck.
   Copyright (C) 2025-2026 Leslie Dancsecs

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
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dancsecs/szargs"
	"github.com/dancsecs/szbck/internal/fstat"
	"github.com/dancsecs/szbck/internal/out"
	"github.com/dancsecs/szbck/internal/rsync"
	"github.com/dancsecs/szbck/internal/settings"
	"github.com/dancsecs/szbck/internal/subcommand/trim"
	"github.com/dancsecs/szlog"
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

func parseArgs(
	args *szargs.Args,
) (*settings.Config, string, bool, bool, error) {
	var (
		cfg       *settings.Config
		isDryRun  bool
		dryRun    string
		trimAfter bool
		daemon    bool
		err       error
	)

	isDryRun = args.Is("--dry-run", "")
	if isDryRun {
		dryRun = " (DRY RUN)"
	}

	trimAfter = args.Is("--trim", "")
	daemon = args.Is("--daemon", "")
	err = args.Err()

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
//nolint:cyclop,funlen,gocognit // Ok.
func Process(args *szargs.Args) (string, error) {
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
		fsStat         *fstat.StatFS
		err            error
		startTime      time.Time
	)

	cfg, dryRunMsg, trimAfter, daemon, err = parseArgs(args)

	if err == nil {
		fsStat, err = fstat.New(cfg.Target.GetPath())
	}

	runOnce := true
	sleepBetweenRuns := time.Nanosecond

	for (runOnce || daemon) && err == nil {
		szlog.Say1f("Starting in: %v\n", sleepBetweenRuns)
		time.Sleep(sleepBetweenRuns)
		startTime = time.Now()
		runOnce = false

		newDir, err = cfg.Target.Create(startTime, initialBackupDirPerm)

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
			if errors.Is(err, trim.ErrNoBackups) ||
				errors.Is(err, trim.ErrOnlyLatest) {
				err = nil
			}

			purgedMsg = " (Purged: " + out.Int(int64(purgedCount)) + ")"
			totalPurged += purgedCount
			totalPurgedMsg = " (Total Purged: " +
				out.Int(int64(totalPurged)) + ")"
		}

		//nolint:forbidigo // Ok.
		if err == nil {
			fmt.Printf("snapshot successful%s%s\nSyncing...\n",
				purgedMsg,
				dryRunMsg,
			)
			fmt.Println(fsStat.Delta())
			fsStat, err = fstat.New(cfg.Target.GetPath())
		}

		sleepBetweenRuns = NextHourIn(time.Since(startTime))
	}

	if err == nil {
		return "", nil
	}

	return "", fmt.Errorf("%w%s: %w", ErrSnapshotError, totalPurgedMsg, err)
}
