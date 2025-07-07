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
	"time"
)

func getProximity(fTme, tTme time.Time) (bool, bool) {
	var (
		sameWeek bool
		sameDay  bool
	)

	sameDay = fTme.Year() == tTme.Year() && fTme.YearDay() == tTme.YearDay()

	fISOYear, fISOWeek := fTme.ISOWeek()
	tISOYear, tISOWeek := tTme.ISOWeek()
	sameWeek = fISOYear == tISOYear && fISOWeek == tISOWeek

	return sameWeek, sameDay
}

func identifyRemovals(tms []time.Time, dayCut, weekCut time.Time) []bool {
	var (
		prevIndex int
		currIndex int
		remove    []bool
	)

	remove = make([]bool, len(tms))

	currIndex = len(tms) - 1
	prevIndex = currIndex - 1

	for ; prevIndex >= 0; prevIndex-- {
		if dayCut.Before(tms[prevIndex]) {
			currIndex = prevIndex

			continue
		}

		sameWeek, sameDay := getProximity(tms[currIndex], tms[prevIndex])

		if weekCut.Before(tms[prevIndex]) { //nolint:nestif // Ok.
			if !sameDay {
				currIndex = prevIndex
			} else {
				remove[prevIndex] = true
			}
		} else {
			if !sameWeek {
				currIndex = prevIndex
			} else {
				remove[prevIndex] = true
			}
		}
	}

	return remove
}
