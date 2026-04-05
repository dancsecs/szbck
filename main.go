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

package main

import (
	"os"

	"github.com/dancsecs/szbck/internal"
)

/*
Simply invokes the internal version of main which returns an int as a classic
type main function.  This is returned to the operating system via the os.Exit
function which cannot be tested.  Therefore this wrapper is the only function
in this utility that is not tested.
*/
func main() {
	os.Exit(internal.Main(os.Args))
}
