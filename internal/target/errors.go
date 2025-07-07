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

package target

import (
	"errors"
)

// Target errors.
var (
	ErrNew                 = errors.New("new error")
	ErrInvalid             = errors.New("invalid target")
	ErrCreateTargetFailed  = errors.New("could not create new target")
	ErrCreateAlreadyExists = errors.New("already exists")
	ErrInvalidLatest       = errors.New("invalid latest symlink")
	ErrHasLatest           = errors.New("has latest failed")
	ErrSplitNotFound       = errors.New("split not found")
	ErrInvalidSplit        = errors.New("invalid directory split")
)
