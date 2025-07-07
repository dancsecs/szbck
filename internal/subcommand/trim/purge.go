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
	"os"
	"time"

	"github.com/dancsecs/szbck/internal/out"
)

func processPurge(
	dirs []string, tms []time.Time, remove []bool, dryRun string,
) (int, error) {
	var (
		purgedCount int
		err         error
	)

	for i, dir := range dirs {
		if dryRun == "" && remove[i] {
			err = os.Chmod(dir, permToDelete)
			if err == nil {
				err = os.RemoveAll(dir)
			}

			if err == nil {
				purgedCount++
			}
		}

		if err == nil {
			if remove[i] {
				out.Print("Purging snapshot"+dryRun+": "+dir+": "+
					tms[i].Format(time.RFC850), "\n",
				)
			} else {
				out.Print("Keeping snapshot"+dryRun+": "+dir+": "+
					tms[i].Format(time.RFC850), "\n",
				)
			}
		}
	}

	if err == nil {
		return purgedCount, nil
	}

	// Note some may have already been purged before error.
	return purgedCount, fmt.Errorf("%w: %w", ErrPurgeFailed, err)
}
