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

func TestInternalSettings_ValRetentionDaily_InvalidBlank(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.Err(
		cfg.validateKeepDaily(""),
		""+
			ErrInvalidKeepDaily.Error()+
			": "+
			ErrMissing.Error()+
			"",
	)
}

func TestInternalSettings_ValRetentionDaily_InvalidSyntax(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.Err(
		cfg.validateKeepDaily("5days"),
		""+
			ErrInvalidKeepDaily.Error()+
			": "+
			ErrSyntax.Error()+
			"",
	)
}

func TestInternalSettings_ValRetentionDaily_InvalidNumber(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.Err(
		cfg.validateKeepDaily("5x days"),
		""+
			ErrInvalidKeepDaily.Error()+
			": "+
			ErrSyntax.Error()+
			"",
	)
}

func TestInternalSettings_ValRetentionDaily_InvalidUnits(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.Err(
		cfg.validateKeepDaily("5 day"),
		""+
			ErrInvalidKeepDaily.Error()+
			": "+
			ErrInvalidUnit.Error()+
			": "+ValidUnits+
			"",
	)
}

func TestInternalSettings_ValRetentionDaily_Low(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.Err(
		cfg.validateKeepDaily("1 days"),
		""+
			ErrInvalidKeepDaily.Error()+
			": "+
			ErrRetentionDailyMin.Error()+
			"",
	)
}

func TestInternalSettings_ValRetentionDaily_ValidHours(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.NoErr(cfg.validateKeepDaily("200 hours"))
}

func TestInternalSettings_ValRetentionDaily_ValidDays(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.NoErr(cfg.validateKeepDaily("30 days"))
}

func TestInternalSettings_ValRetentionDaily_Duplicate(t *testing.T) {
	chk := sztestlog.CaptureNothing(t, szlog.LevelAll)
	defer chk.Release()

	var cfg Config

	chk.NoErr(cfg.validateKeepDaily("24 days"))
	chk.Err(
		cfg.validateKeepDaily("200 hours"),
		""+
			ErrInvalidKeepDaily.Error()+
			": "+
			ErrDuplicate.Error()+
			": 'keepDaily'"+
			"",
	)
}
