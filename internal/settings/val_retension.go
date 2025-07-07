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

package settings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	keepHourly = "keepHourly"
	keepDaily  = "keepDaily"
)

// Valid time units message.
const (
	UnitHours  = "hours"
	UnitDays   = "days"
	ValidUnits = "must be '" +
		UnitHours + "'" +
		"' or '" + UnitDays + "'" +
		""
)

// Retention minimum durations.
const (
	minHourly = time.Hour * 24
	minDaily  = time.Hour * 48
)

// Option errors.
var (
	ErrInvalidKeepHourly  = errors.New("invalid hourly retention")
	ErrInvalidKeepDaily   = errors.New("invalid daily retention")
	ErrInvalidUnit        = errors.New("invalid time unit")
	ErrRetentionHourlyMin = errors.New("must be >= 24 hours")
	ErrRetentionDailyMin  = errors.New(
		"must be > retention hours or >= 48 hours",
	)
)

//nolint:cyclop,funlen // Ok.
func validateTimeUnit(
	name string,
	currentValue *time.Duration,
	value string,
) error {
	const (
		base10      = 10
		bits64      = 64
		hoursPerDay = 24
	)

	var (
		amountStr string
		amount    int64
		units     string
		found     bool
		err       error
	)

	if *currentValue != 0 {
		err = fmt.Errorf("%w: '%s'", ErrDuplicate, name)
	}

	if err == nil && value == "" {
		err = ErrMissing
	}

	if err == nil {
		amountStr, units, found = strings.Cut(value, " ")
		if !found {
			err = ErrSyntax
		} else {
			units = strings.TrimSpace(units)
		}
	}

	if err == nil {
		amount, err = strconv.ParseInt(amountStr, base10, bits64)
		if err != nil {
			err = ErrSyntax
		}
	}

	if err == nil {
		switch units {
		case UnitHours:
			*currentValue = time.Duration(amount) * time.Hour

			return nil
		case UnitDays:
			*currentValue = time.Duration(amount) * time.Hour * hoursPerDay

			return nil
		default:
			err = fmt.Errorf(
				"%w: %s",
				ErrInvalidUnit, ValidUnits,
			)
		}
	}

	return err
}

func (cfg *Config) validateKeepHourly(value string) error {
	err := validateTimeUnit(
		keepHourly, &cfg.KeepHourly, value,
	)

	if err == nil && cfg.KeepHourly < minHourly {
		err = ErrRetentionHourlyMin
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", ErrInvalidKeepHourly, err)
}

func (cfg *Config) validateKeepDaily(value string) error {
	err := validateTimeUnit(
		keepDaily, &cfg.KeepDaily, value,
	)

	if err == nil && cfg.KeepDaily < minDaily {
		err = ErrRetentionDailyMin
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", ErrInvalidKeepDaily, err)
}
