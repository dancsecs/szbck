/*
   Golang rsync backup utility  wrapper: szbck.
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

package settings

import (
	"testing"

	"github.com/dancsecs/szlog"
	"github.com/dancsecs/sztestlog"
)

func TestInternalSettings_ValRetentionHourly_InvalidBlank(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.Err(
		cfg.validateKeepHourly(""),
		""+
			ErrInvalidKeepHourly.Error()+
			": "+
			ErrMissing.Error()+
			"",
	)
}

func TestInternalSettings_ValRetentionHourly_InvalidSyntax(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.Err(
		cfg.validateKeepHourly("5hours"),
		""+
			ErrInvalidKeepHourly.Error()+
			": "+
			ErrSyntax.Error()+
			"",
	)
}

func TestInternalSettings_ValRetentionHourly_InvalidNumber(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.Err(
		cfg.validateKeepHourly("5x hours"),
		""+
			ErrInvalidKeepHourly.Error()+
			": "+
			ErrSyntax.Error()+
			"",
	)
}

func TestInternalSettings_ValRetentionHourly_InvalidUnits(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.Err(
		cfg.validateKeepHourly("5 hour"),
		""+
			ErrInvalidKeepHourly.Error()+
			": "+
			ErrInvalidUnit.Error()+
			": "+ValidUnits+
			"",
	)
}

func TestInternalSettings_ValRetentionHourly_Low(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.Err(
		cfg.validateKeepHourly("1 hours"),
		""+
			ErrInvalidKeepHourly.Error()+
			": "+
			ErrRetentionHourlyMin.Error()+
			"",
	)
}

func TestInternalSettings_ValRetentionHourly_ValidHours(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.NoErr(cfg.validateKeepHourly("24 hours"))
}

func TestInternalSettings_ValRetentionHourly_ValidDays(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.NoErr(cfg.validateKeepHourly("1 days"))
}

func TestInternalSettings_ValRetentionHourly_Duplicate(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.NoErr(cfg.validateKeepHourly("24 hours"))
	chk.Err(
		cfg.validateKeepHourly("2 days"),
		""+
			ErrInvalidKeepHourly.Error()+
			": "+
			ErrDuplicate.Error()+
			": 'keepHourly'"+
			"",
	)
}
