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

package purge

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Directory enables write access to all entries in the dir hierarchy
// and purges the whole directory.
func Directory(dir string) error {
	var (
		out   []byte
		lines []string
	)

	findPath, err := exec.LookPath("find")

	if err == nil {
		//nolint:gosec // Ok.
		cmd := exec.Command(
			findPath,
			dir, "-type", "d",
			"!", "-perm", "-u=w",
			"-print",
		)
		out, err = cmd.CombinedOutput()
		lines = strings.Split(strings.Trim(string(out), " \t\n"), "\n")
	}

	if err != nil && len(lines) == 1 && strings.HasPrefix(lines[0], findPath) {
		err = fmt.Errorf("%w: %s", ErrDirRights, lines[0])
	}

	for i := 0; err == nil && i < len(lines); i++ {
		const fullDirPermissions = 0o0700

		if lines[i] != "" {
			err = os.Chmod(lines[i], fullDirPermissions)
		}
	}

	if err == nil {
		err = os.RemoveAll(dir)
	}

	return err //nolint:wrapCheck // Ok.
}
