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

package rsync

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/dancsecs/szbck/internal/out"
)

// Run executes rsync with the supplied arguments.
func Run(args []string, cpyOut, cpyErr *os.File) error {
	var (
		rsyncPath    string
		cmd          *exec.Cmd
		outCloser    io.ReadCloser
		outCloserErr io.ReadCloser
		err          error
	)

	rsyncPath, err = exec.LookPath("rsync")

	//nolint:nestif // Ok.
	if err == nil {
		out.Printf(
			"Running command: %s %s\n", rsyncPath, strings.Join(args, " "),
		)

		cmd = exec.Command(rsyncPath, args...) //nolint:gosec // Ok.

		if cpyOut != nil {
			outCloser, err = cmd.StdoutPipe()
			if err == nil {
				go func() {
					_, _ = io.Copy(cpyOut, outCloser)
					_ = outCloser.Close()
				}()
			}
		}

		if cpyErr != nil {
			outCloserErr, err = cmd.StderrPipe()
			if err == nil {
				go func() {
					_, _ = io.Copy(cpyErr, outCloserErr)
					_ = outCloserErr.Close()
				}()
			}
		}

		err = cmd.Run()
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", ErrRsyncError, err)
}
