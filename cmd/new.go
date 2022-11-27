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

package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"tfsw/internal/utils"
)

var (
	newCmd = &cobra.Command{
		Args:   cobra.MinimumNArgs(1),
		Long:   "Installs the specified Terraform versions to the local cache",
		PreRun: validateVersion,
		Run:    newRun,
		Short:  "Install new versions",
		Use:    "new VERSION...",
	}
)

func init() {
	rootCmd.AddCommand(newCmd)
}

// newRun is passed directly to the Cobra Run argument and executes
// the primary logic for the `new` command
func newRun(cmd *cobra.Command, args []string) {
	for _, version := range args {
		err := newVersion(version)
		switch err {
		case ErrVersionExists:
			fmt.Printf("Terraform %s already exists\n", version)
		case nil:
			fmt.Printf("Terraform %s has been added\n", version)
		default:
			fmt.Fprintf(os.Stderr, "Encountered an unhandled error: %v\n", err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}

// newVersion takes a Terraform version number, downloads it,
// checks shasums, then unzips it into the correct location
func newVersion(ver string) error {
	if _, err := os.Stat(filepath.Join(config.ConfigDirectory, ver, terraform)); err == nil {
		return ErrVersionExists
	}

	os.MkdirAll(config.TempDirectory, 0755)
	defer os.RemoveAll(config.TempDirectory)

	zip := strings.Join([]string{"terraform", ver, runtime.GOOS, runtime.GOARCH}, "_") + ".zip"
	sums := strings.Join([]string{"terraform", ver, "SHA256SUMS"}, "_")
	urlRoot := path.Join(config.RepositoryDomain, "terraform", ver)

	err := utils.FetchUrl("https://"+path.Join(urlRoot, zip), filepath.Join(config.TempDirectory, zip))
	if err != nil {
		return err
	}

	err = utils.FetchUrl("https://"+path.Join(urlRoot, sums), filepath.Join(config.TempDirectory, sums))
	if err != nil {
		return err
	}

	ok, err := utils.Sha256sum(filepath.Join(config.TempDirectory, sums), filepath.Join(config.TempDirectory, zip))
	if !ok && err != nil {
		return err
	}

	versionDir := filepath.Join(config.ConfigDirectory, ver)
	os.MkdirAll(versionDir, 0755)

	_, err = utils.Unzip(filepath.Join(config.TempDirectory, zip), versionDir)
	if err != nil {
		return err
	}

	return nil
}
