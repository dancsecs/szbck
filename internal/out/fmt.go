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

package out

import (
	"github.com/dancsecs/szlog"
	"golang.org/x/text/message"
)

//nolint:goCheckNoGlobals // OK
var (
	printer = message.NewPrinter(message.MatchLanguage("en"))
)

// Int returns the number formatted with local separators.
func Int(n int64) string {
	return printer.Sprintf("%d", n)
}

// Uint returns the number formatted with local separators.
func Uint(n uint64) string {
	return printer.Sprintf("%d", n)
}

// Pct returns the number formatted as a percent with local separators.
func Pct(n float64) string {
	return printer.Sprintf("%.2f", n*100.0) + "%" //nolint:mnd  // Ok.
}

// Print writes the provided text to os.Stdout of szlog has warnings enabled.
func Print(msg ...any) {
	if szlog.Verbose() > 0 {
		_, _ = printer.Print(msg...)
	}
}

// Printf writes the provided text to os.Stdout of szlog has warnings enabled.
func Printf(msgFmt string, msgArgs ...any) {
	if szlog.Verbose() > 0 {
		_, _ = printer.Printf(msgFmt, msgArgs...)
	}
}
