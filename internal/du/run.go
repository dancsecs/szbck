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

package du

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/dancsecs/szbck/internal/directory"
)

// Run executes rsync with the supplied arguments.
func Run(args []string, cpyErr *os.File) (string, error) {
	var (
		duPath string
		cmd    *exec.Cmd
		out    []byte
		outErr io.ReadCloser
		err    error
	)

	if len(args) == 0 {
		err = ErrMissing
	}

	if err == nil {
		err = directory.Is(args[len(args)-1])
	}

	if err == nil {
		duPath, err = exec.LookPath("du")
	}

	if err == nil {
		cmd = exec.Command(duPath, args...) //nolint:gosec // Ok.

		if cpyErr != nil {
			outErr, err = cmd.StderrPipe()
			if err == nil {
				go func() {
					_, _ = io.Copy(cpyErr, outErr)
					_ = outErr.Close()
				}()
			}
		}

		out, err = cmd.Output()
	}

	if err == nil {
		return string(out), nil
	}

	return "", fmt.Errorf(
		"%w: %w: command: '%s %s'",
		ErrDuError,
		err,
		duPath,
		strings.Join(args, " "),
	)
}
