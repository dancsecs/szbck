<!---             *****  AUTO GENERATED:  DO NOT MODIFY  ***** -->
<!---                   MODIFY TEMPLATE: '.README.gtm.md' -->
<!---               See: 'https://github.com/dancsecs/gotomd' -->

<!---
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
-->

# Szerszam Backup Utility

```go
package main
```

Package szbck wraps the system's rsync utility to create time machine like
backups. It is driven by a backup configuration file which identifies the root
directory to backup, a target directory to store snapshots, as well as
enabling and disabling various rsync options, defines the default permissions
for the root of each backup (snapshot) and lists items to exclude from backup
and restores.

It can be used to backup a single machine and/or to keep several machines in
sync.  The utility defines several subcommands implementing the various backup
functions as follows:

    SubCommand  Description
    ==========  =============================================================
    Help        Displays help on the utility and all of the subcommands.

    Create      Creates a backup configuration file.

    Snapshot    Creates a new backup snapshot hard linking unchanged items to
                previous snapshots.  It will implement a retention policy if
                and only if the --trim option is provided.  It can be run in
                --daemon mode where it will run each hour.  Further options
                --at and --monitor can be provided in --daemon mode to specify
                the minute after the hour to run (defaults to the start time
                minute) and if a countdown should be displayed.

    Restore     Restores/removes/replaces files in the source directory
                identified by the backup configuration file honoring all
                exclusions identified by the config file from a backup
                snapshot. The backup snapshot need not be from the same machine
                enabling syncing between machines.

    Prune       Removes oldest backup snapshots.  NOTE: This does not
                implement a retention policy but simply purges the
                specified number of the oldest snapshots.  A retention
                policy has not yet been implemented.

    Trim        Manually implements the retention policy as specified in the
                configuration file.  Can be invoked by using the --trim option
                on the snapshot subcommand.

    Status        Reports on the number of backup snapshots and the space used
                overall and by each snapshot.

    Vet         Parses a backup configuration file identifying any errors
                or problems without making any attempts at any operations.

    szbck
    Szerszam backup utility takes Apple Time machine like snapshots.  It
    requires the underlying system to have the utility rsync installed which
    will perform the actual snapshots.  One of the following sub commands must
    be provided as follows:

    {h | help} [subcommand]

    Help information is displayed.

       [subcommand]
          Limits the help displayed to the specified subcommand.

    {c | create} [-o filePath] [-p perm] [-e exclude ...] source

    Create generates a new Szerszam backup configuration file.

        [-o config.sbc]
            Writes the configuration to the named Szerszam backup Configuration
          file.  Otherwise it is written to stdout.

        [-t target]
            Specifies the default target directory to create snapshots in.  If it
            not provided then a target argument will be mandatory for snapshot,
            restore, prune and status subcommands.

        source
            Specifies the root directory to back up.

    {s | snap | snapshot} [--dry-run] [--daemon [--at minute] [--monitor]] [--trim] [-t target] config.szb

    Create a new snapshot of the source listed in the configuration file located
    in the target directory.

       [--dry-run]
          Identifies all of the actions the utility would take without making any
          changes to the backup source.

       [--daemon]
          Runs in a loop creating a snapshot every hour and sleeping between runs.

       [--at minute]
          If daemon mode is enabled this specifies the minute after the hour the
          snapshot should start otherwise the current time's minute value will be
          used.  Valid values are between 0-59.  An unexpected argument error will
          occur if specified without --daemon being specified.

       [--monitor]
          If daemon mode is enabled then a countdown until the next backup is
          displayed.  By minute, then by second for the last minute.

       [--trim]
          Executes the retention policy as specified in the configuration file
          after the snapshot has been successfully completed.

       [-t target]
          Specifies the backup set to create the new snapshot in.  It is optional
          if the backup config file specifies a target and mandatory if not
          specified in the backup config file.

       config.sbc
          The backup configuration file defining the backup.

    {r | rest | restore} [--dry-run] [--keep] [-s snapshot] [-t target] config.szb

    Restores the specified file or directory tree from the backup.

       [--dry-run]
          Identifies all of the actions the utility would take without making any
          changes to the backup source.

       [--keep]
          Blocks the restore from deleting source files missing from the target
          backup.

       [-s snapshot]
          Specifies the specif snapshot in the target directory to use.  It will
          default to the symbolic link 'latest' is not provided.

       [-t target]
          Specifies the backup set to restore from.  It is optional if the backup
          config file specifies a target and mandatory if not specified in the
          backup config file.

       config.sbc
          the backup configuration file defining the backup.

    {p | prune} [--dry-run] [-n {number | all}] [-t target] config.szb

    Deletes the oldest backups.  Defaults to 1.
    NOTE:  The latest backup will not be purged.

       [--dry-run]
          Identifies all of the actions the utility would take without making any
          changes to the backup source.

       [-n {number | all}]
          Specify the number of older backups to purge or all previous backups.

       [-t target]
          Specifies the backup set to prune.  It is optional if the backup config
          file specifies a target and mandatory if not specified in the backup
          config file.

       config.sbc
          the backup configuration file defining the backup.

    {stat | status} [-t target] config.szb

    Reports the status on the specified backup set.

       [-t target]
          Specifies the backup set to create the new snapshot in.  It is optional
          if the backup config file specifies a target and mandatory if not
          specified in the backup config file.

       config.sbc
          The backup configuration file defining the backup.

    {t | trim} [--dry-run] [-t target] config.szb

    Implements the specified retention policy as defined in the backup
    configuration file deleting backups as appropriate. The most recent snapshot
    pointed to by the "latest" symbolic link is never deleted.

       [--dry-run]
          Identifies all of the actions the utility would take without making any
          changes to the backup source.

       [-t target]
          Specifies the backup set to prune.  It is optional if the backup config
          file specifies a target and mandatory if not specified in the backup
          config file.

       config.sbc
          the backup configuration file defining the backup.

    {v | vet} config.szb

    Loads and parses the named configuration files reporting any issues.

       config.sbc
          the backup configuration file defining the backup.

# Examples:

    // Display help on the utility and all sub commands.
        szbck help

    // Creates a new backup configuration file.
        szbck create -o config.szb -t backupTo /home/myDirectory

    // Create a new snapshot.
        szbck snapshot config.szb

    // Restore updated/missing files and purge extra files (unless the --keep
    // option is specified).
        szbck restore config.szb

    // Purge the oldest 5 snapshots from a backup set.
        szbck prune -n 5 config.szb

    // Vet changes made to a config.szb file.
        szbck vet config.szb

# Dedication

This project is dedicated to Reem.
Your brilliance, courage, and quiet strength continue to inspire me.
Every line is written in gratitude for the light and hope you brought into my
life.

NOTE: Documentation reviewed and polished with the assistance of ChatGPT from
OpenAI.
