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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/schollz/progressbar/v3"
)

func FetchUrl(uri, dstFile string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	_, err = http.Head(uri)
	if err != nil {
		return fmt.Errorf("failed to HEAD %s got: %v\n", uri, err)
	}

	res, err := http.Get(uri)
	if err != nil {
		return fmt.Errorf("failed to GET %s got: %v\n", uri, err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		//dstFile := filepath.Join(dir, path.Base(url))
		file, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("unable to open file handle: %v\n", err)
		}
		defer file.Close()

		bar := progressbar.DefaultBytes(
			-1,
			"downloading",
		)
		_, err = io.Copy(io.MultiWriter(file, bar), res.Body)
		if err != nil {
			return fmt.Errorf("failed io.Copy: %v\n", err)
		}

		if err = bar.Clear(); err != nil {
			return fmt.Errorf("failed to clear the progress bar: %v\n", err)
		}

		return nil
	case http.StatusForbidden:
		return fmt.Errorf("the file requested doesn't exist on %s", u.Hostname())
	default:
		return fmt.Errorf("did not get HTTP 200, got %d instead", res.StatusCode)
	}

	return nil
}
