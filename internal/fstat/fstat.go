/*
   Golang rsync backup utility wrapper: szbck.
   Copyright (C) 2026 Leslie Dancsecs

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

package fstat

import (
	"fmt"
	"strings"

	"github.com/dancsecs/szbck/internal/out"
	"golang.org/x/sys/unix"
)

// StatFS captures file system sizes and reports deltas.
type StatFS struct {
	path        string
	totalBytes  uint64
	freeBytes   uint64
	totalINodes uint64
	freeINodes  uint64
}

// New creates a file system object populated with the total and free
// bytes available.
func New(path string) (*StatFS, error) {
	var initial unix.Statfs_t

	statFS := new(StatFS)
	statFS.path = path

	unix.Sync()
	err := unix.Statfs(path, &initial)

	//nolint:gosec // Ok.
	if err == nil {
		statFS.totalBytes = uint64(initial.Bsize) * initial.Blocks
		statFS.freeBytes = uint64(initial.Bsize) * initial.Bfree
		statFS.totalINodes = initial.Files
		statFS.freeINodes = initial.Ffree

		return statFS, nil
	}

	return nil, fmt.Errorf("statfs failed: %w", err)
}

func balancePct(pct string) string {
	if pct == "" {
		return ""
	}

	return fmt.Sprintf("(%7s)", pct)
}

func report(
	name, byteSize, bytePercent, iNode, iNodePct string,
) string {
	return strings.TrimRight(
		fmt.Sprintf(
			"%13s %20s %9s %20s %9s",
			name,
			byteSize,
			balancePct(bytePercent),
			iNode,
			balancePct(iNodePct),
		),
		" ",
	)
}

// TotalStatus returns a string representing the file system total bytes and
// iNodes.
func (a *StatFS) TotalStatus() string {
	return fmt.Sprint(
		report(
			"",
			"Bytes",
			"",
			"INodes",
			"",
		),
		"\n",
		report(
			"Capacity:",
			out.Uint(a.totalBytes),
			"",
			out.Uint(a.totalINodes),
			"",
		),
	)
}

// FreeStatus returns a string representing the file system usage.
func (a *StatFS) FreeStatus(title string) string {
	return report(
		"Free "+title+":",
		out.Uint(a.freeBytes),
		out.Pct(float64(a.freeBytes)/float64(a.totalBytes)),
		out.Uint(a.freeINodes),
		out.Pct(float64(a.freeINodes)/float64(a.totalINodes)),
	)
}

// Delta returns a string representing the file system usage changes.
func (a *StatFS) Delta() string {
	deltaStatFS, _ := New(a.path)

	//nolint:gosec // OK.
	deltaBytes := int64(a.freeBytes - deltaStatFS.freeBytes)
	//nolint:gosec // OK.
	deltaINodes := int64(a.freeINodes - deltaStatFS.freeINodes)

	return fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		a.TotalStatus(),
		a.FreeStatus("Before"),
		deltaStatFS.FreeStatus("After"),
		report(
			"Delta:",
			out.Int(deltaBytes),
			out.Pct(float64(deltaBytes)/float64(a.totalBytes)),
			out.Int(deltaINodes),
			out.Pct(float64(deltaINodes)/float64(a.totalINodes)),
		),
	)
}
