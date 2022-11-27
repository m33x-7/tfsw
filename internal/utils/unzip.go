/*
Terraform Switch - A commandline utility to manage multiple versions
of HashiCorps infrastructure as code tool, Terraform

Copyright (C) 2022  Tom Cole <tom@m33x-7.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License along
with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package utils

import (
	"archive/zip"
	"github.com/schollz/progressbar/v3"
	"io"
	"os"
	"path/filepath"
)

func Unzip(src, dst string) (int, error) {
	archive, err := zip.OpenReader(src)
	if err != nil {
		return 0, err
	}
	defer archive.Close()

	var unarchived int = 0
	for _, f := range archive.File {
		path := filepath.Join(dst, f.Name)

		// TODO: If this fails it should clean up rather than
		// leave a partically unarchived zip?
		dstfile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return unarchived, err
		}
		defer dstfile.Close()

		srcfile, err := f.Open()
		if err != nil {
			return unarchived, err
		}
		defer srcfile.Close()

		bar := progressbar.DefaultBytes(
			-1,
			"extracting ",
		)

		if _, err := io.Copy(io.MultiWriter(dstfile, bar), srcfile); err != nil {
			return unarchived, err
		}

		_ = bar.Clear()

		unarchived++
	}

	return unarchived, nil
}
