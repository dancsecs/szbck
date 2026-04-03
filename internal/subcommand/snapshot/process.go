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
	"github.com/dancsecs/szbck/internal/wait"
)

const initialBackupDirPerm = 0o0700

//nolint:nestif,cyclop,funlen // Ok.
func parseArgs(
	args *szargs.Args,
	startTime time.Time,
) (*settings.Config, string, bool, bool, int, bool, error) {
	const maxMinute = 59

	var (
		cfg       *settings.Config
		isDryRun  bool
		dryRun    string
		trimAfter bool
		daemon    bool
		runAtMin  uint8
		foundAt   bool
		monitor   bool
		err       error
	)

	isDryRun = args.Is("--dry-run", "")
	if isDryRun {
		dryRun = " (DRY RUN)"
	}

	trimAfter = args.Is("--trim", "")

	daemon = args.Is("--daemon", "")

	monitor = args.Is("--monitor", "")

	runAtMin, foundAt = args.ValueUint8("--at", "")

	if !args.HasErr() {
		if !daemon {
			if foundAt {
				runAtMin = 0

				args.PushErr(ErrAtUsage)
			}

			if monitor {
				monitor = false

				args.PushErr(ErrMonitorUsage)
			}
		} else {
			if errors.Is(args.Err(), szargs.ErrRange) || runAtMin > maxMinute {
				args.PushErr(ErrAtRange)

				runAtMin = 0
			}

			if !args.HasErr() && !foundAt {
				//nolint:gosec // Ok. Daemon defaulting to current minute.
				runAtMin = min(uint8(startTime.Minute()), maxMinute)
			}
		}
	}

	err = args.Err()

	if err == nil {
		cfg, err = settings.LoadFromArgs(args)
	}

	return cfg, dryRun, trimAfter, daemon, int(runAtMin), monitor, err
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
		runAtMin       int
		monitor        bool
		purgedCount    int
		purgedMsg      string
		totalPurged    int
		totalPurgedMsg string
		hasLatest      bool
		linkDest       string
		newDir         string
		fsStat         *fstat.StatFS
		err            error
	)

	cfg, dryRunMsg, trimAfter, daemon, runAtMin, monitor, err = parseArgs(
		args,
		time.Now(),
	)

	if err == nil {
		fsStat, err = fstat.New(cfg.Target.GetPath())
	}

	runOnce := true
	targetRunTime := time.Now()

	for (runOnce || daemon) && err == nil {
		wait.Until("Next Backup", monitor, targetRunTime)

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

		targetRunTime = wait.NextHourAt(runAtMin, time.Now())
	}

	if err == nil {
		return "", nil
	}

	return "", fmt.Errorf("%w%s: %w", ErrSnapshotError, totalPurgedMsg, err)
}
