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
	"testing"
	"time"

	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

type tmeEntry struct {
	tme          time.Time
	removeWeekly bool
	removeDaily  bool
}

func mkTme(y, m, d, h int, delWeekly, delDaily bool) tmeEntry {
	return tmeEntry{
		tme: time.Date(
			y, time.Month(m), d, h, 1, 2, 3456000000, time.Local,
		),
		removeWeekly: delWeekly,
		removeDaily:  delDaily,
	}
}

func mkTms(t []tmeEntry) []time.Time {
	tms := make([]time.Time, len(t))
	for i, tm := range t {
		tms[i] = tm.tme
	}

	return tms
}

func mkRemovedWeekly(t []tmeEntry) []bool {
	remove := make([]bool, len(t))
	for i, tm := range t {
		remove[i] = tm.removeWeekly
	}

	return remove
}

func mkRemovedDaily(t []tmeEntry) []bool {
	remove := make([]bool, len(t))
	for i, tm := range t {
		remove[i] = tm.removeDaily
	}

	return remove
}

func TestInternalTrim_Identify_Dec31_Sunday(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	tme := []tmeEntry{
		mkTme(2023, 12, 23, 10, true, true),   //  0 - Sat.
		mkTme(2023, 12, 23, 20, true, false),  //  1 - Sat.
		mkTme(2023, 12, 24, 10, true, true),   //  2 - Sun.
		mkTme(2023, 12, 24, 20, false, false), //  3 - Sun.
		mkTme(2023, 12, 25, 10, true, true),   //  4 - Mon.
		mkTme(2023, 12, 25, 20, true, false),  //  5 - Mon.
		mkTme(2023, 12, 26, 10, true, true),   //  6 - Tue.
		mkTme(2023, 12, 26, 20, true, false),  //  7 - Tue.
		mkTme(2023, 12, 27, 10, true, true),   //  8 - Wed.
		mkTme(2023, 12, 27, 20, true, false),  //  9 - Wed.
		mkTme(2023, 12, 28, 10, true, true),   // 10 - Thu.
		mkTme(2023, 12, 28, 20, true, false),  // 11 - Thu.
		mkTme(2023, 12, 29, 10, true, true),   // 12 - Fri.
		mkTme(2023, 12, 29, 20, true, false),  // 13 - Fri.
		mkTme(2023, 12, 30, 10, true, true),   // 14 - Sat.
		mkTme(2023, 12, 30, 20, true, false),  // 15 - Sat.
		mkTme(2023, 12, 31, 10, true, true),   // 16 - Sun.
		mkTme(2023, 12, 31, 20, false, false), // 17 - Sun.
		mkTme(2024, 1, 1, 10, true, true),     // 18 - Mon.
		mkTme(2024, 1, 1, 20, true, false),    // 19 - Mon.
		mkTme(2024, 1, 2, 10, true, true),     // 20 - Tue.
		mkTme(2024, 1, 2, 20, true, false),    // 21 - Tue.
		mkTme(2024, 1, 3, 10, true, true),     // 22 - Wed.
		mkTme(2024, 1, 3, 20, true, false),    // 23 - Wed.
		mkTme(2024, 1, 4, 10, true, true),     // 24 - Thu.
		mkTme(2024, 1, 4, 20, true, false),    // 25 - Thu.
		mkTme(2024, 1, 5, 10, true, true),     // 26 - Fri.
		mkTme(2024, 1, 5, 20, true, false),    // 27 - Fri.
		mkTme(2024, 1, 6, 10, true, true),     // 28 - Sat.
		mkTme(2024, 1, 6, 20, true, false),    // 29 - Sat.
		mkTme(2024, 1, 7, 10, true, true),     // 30 - Sun.
		mkTme(2024, 1, 7, 20, false, false),   // 31 - Sun.
		mkTme(2024, 1, 8, 10, true, true),     // 32 - Mon.
		mkTme(2024, 1, 8, 20, false, false),   // 33 - Mon.
	}

	chk.BoolSlice(
		identifyRemovals(mkTms(tme), time.Now(), time.Now()),
		mkRemovedWeekly(tme),
	)

	chk.BoolSlice(
		identifyRemovals(mkTms(tme), time.Now(), tme[0].tme),
		mkRemovedDaily(tme),
	)
}

func TestInternalTrim_Identify_Dec31_Monday(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	tme := []tmeEntry{
		mkTme(2018, 12, 22, 10, true, true),   //  0 - Sat.
		mkTme(2018, 12, 22, 20, true, false),  //  1 - Sat.
		mkTme(2018, 12, 23, 10, true, true),   //  2 - Sun.
		mkTme(2018, 12, 23, 20, false, false), //  3 - Sun.
		mkTme(2018, 12, 24, 10, true, true),   //  4 - Mon.
		mkTme(2018, 12, 24, 20, true, false),  //  5 - Mon.
		mkTme(2018, 12, 25, 10, true, true),   //  6 - Tue.
		mkTme(2018, 12, 25, 20, true, false),  //  7 - Tue.
		mkTme(2018, 12, 26, 10, true, true),   //  8 - Wed.
		mkTme(2018, 12, 26, 20, true, false),  //  9 - Wed.
		mkTme(2018, 12, 27, 10, true, true),   // 10 - Thu.
		mkTme(2018, 12, 27, 20, true, false),  // 11 - Thu.
		mkTme(2018, 12, 28, 10, true, true),   // 12 - Fri.
		mkTme(2018, 12, 28, 20, true, false),  // 13 - Fri.
		mkTme(2018, 12, 29, 10, true, true),   // 14 - Sat.
		mkTme(2018, 12, 29, 20, true, false),  // 15 - Sat.
		mkTme(2018, 12, 30, 10, true, true),   // 16 - Sun.
		mkTme(2018, 12, 30, 20, false, false), // 17 - Sun.
		mkTme(2018, 12, 31, 10, true, true),   // 18 - Mon.
		mkTme(2018, 12, 31, 20, true, false),  // 19 - Mon.
		mkTme(2019, 1, 1, 10, true, true),     // 20 - Tue.
		mkTme(2019, 1, 1, 20, true, false),    // 21 - Tue.
		mkTme(2019, 1, 2, 10, true, true),     // 22 - Wed.
		mkTme(2019, 1, 2, 20, true, false),    // 23 - Wed.
		mkTme(2019, 1, 3, 10, true, true),     // 24 - Thu.
		mkTme(2019, 1, 3, 20, true, false),    // 25 - Thu.
		mkTme(2019, 1, 4, 10, true, true),     // 26 - Fri.
		mkTme(2019, 1, 4, 20, true, false),    // 27 - Fri.
		mkTme(2019, 1, 5, 10, true, true),     // 28 - Sat.
		mkTme(2019, 1, 5, 20, true, false),    // 29 - Sat.
		mkTme(2019, 1, 6, 10, true, true),     // 30 - Sun.
		mkTme(2019, 1, 6, 20, false, false),   // 31 - Sun.
		mkTme(2019, 1, 7, 10, true, true),     // 32 - Mon.
		mkTme(2019, 1, 7, 20, false, false),   // 33 - Mon.
	}

	chk.BoolSlice(
		identifyRemovals(mkTms(tme), time.Now(), time.Now()),
		mkRemovedWeekly(tme),
	)

	chk.BoolSlice(
		identifyRemovals(mkTms(tme), time.Now(), tme[0].tme),
		mkRemovedDaily(tme),
	)
}

func TestInternalTrim_Identify_Dec31_Tuesday(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	tme := []tmeEntry{
		mkTme(2024, 12, 21, 10, true, true),   //  0 - Mon.
		mkTme(2024, 12, 21, 20, true, false),  //  1 - Mon.
		mkTme(2024, 12, 22, 10, true, true),   //  2 - Sun.
		mkTme(2024, 12, 22, 20, false, false), //  3 - Sun.
		mkTme(2024, 12, 23, 10, true, true),   //  4 - Mon.
		mkTme(2024, 12, 23, 20, true, false),  //  5 - Mon.
		mkTme(2024, 12, 24, 10, true, true),   //  6 - Tue.
		mkTme(2024, 12, 24, 20, true, false),  //  7 - Tue.
		mkTme(2024, 12, 25, 10, true, true),   //  8 - Wed.
		mkTme(2024, 12, 25, 20, true, false),  //  9 - Wed.
		mkTme(2024, 12, 26, 10, true, true),   // 10 - Thu.
		mkTme(2024, 12, 26, 20, true, false),  // 11 - Thu.
		mkTme(2024, 12, 27, 10, true, true),   // 12 - Fri.
		mkTme(2024, 12, 27, 20, true, false),  // 13 - Fri.
		mkTme(2024, 12, 28, 10, true, true),   // 14 - Sat.
		mkTme(2024, 12, 28, 20, true, false),  // 15 - Sat.
		mkTme(2024, 12, 29, 10, true, true),   // 16 - Sun.
		mkTme(2024, 12, 29, 20, false, false), // 17 - Sun.
		mkTme(2024, 12, 30, 10, true, true),   // 18 - Mon.
		mkTme(2024, 12, 30, 20, true, false),  // 19 - Mon.
		mkTme(2024, 12, 31, 10, true, true),   // 20 - Tue.
		mkTme(2024, 12, 31, 20, true, false),  // 21 - Tue.
		mkTme(2025, 1, 1, 10, true, true),     // 22 - Wed.
		mkTme(2025, 1, 1, 20, true, false),    // 23 - Wed.
		mkTme(2025, 1, 2, 10, true, true),     // 24 - Thu.
		mkTme(2025, 1, 2, 20, true, false),    // 25 - Thu.
		mkTme(2025, 1, 3, 10, true, true),     // 26 - Fri.
		mkTme(2025, 1, 3, 20, true, false),    // 27 - Fri.
		mkTme(2025, 1, 4, 10, true, true),     // 28 - Sat.
		mkTme(2025, 1, 4, 20, true, false),    // 29 - Sat.
		mkTme(2025, 1, 5, 10, true, true),     // 30 - Sun.
		mkTme(2025, 1, 5, 20, false, false),   // 31 - Sun.
		mkTme(2025, 1, 6, 10, true, true),     // 32 - Mon.
		mkTme(2025, 1, 6, 20, false, false),   // 33 - Mon.
	}

	chk.BoolSlice(
		identifyRemovals(mkTms(tme), time.Now(), time.Now()),
		mkRemovedWeekly(tme),
	)

	chk.BoolSlice(
		identifyRemovals(mkTms(tme), time.Now(), tme[0].tme),
		mkRemovedDaily(tme),
	)
}
