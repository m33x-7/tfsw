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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

const (
	expr string = `^([0-9]+\.){2}[0-9]+(-(alpha|beta|oci|rc)[0-9]*)?$`
)

var (
	basename                          = filepath.Base(os.Args[0])
	config                            = &configuration{}
	ErrNoneInstalled   error          = errors.New("no versions installed")
	ErrVersionNotExist error          = errors.New("file or directory doesn't exist")
	ErrVersionExists   error          = errors.New("version already exists")
	ErrVersionSame     error          = errors.New("new version is the same as old version")
	regex              *regexp.Regexp = regexp.MustCompile(expr)
	rootCmd                           = &cobra.Command{
		Long:  "Terraform Switch allows adding, removing, and switching, between multiple versions of Terraform",
		Short: "tfsw manages Terraform versions",
		Use:   "tfsw",
	}
	versionCurrent    string
	versionsInstalled []string
)

// TODO - If terraform is in path, but it's not TF, it allow add to work, but if terraform
//  is then removed you can't get it to work without manually

func Execute() error {
	// TODO: Detect if `terraform` is already on $PATH and exit
	// TODO: CLI doesn't work without ~/.config/tfsw being created
	// in advance
	err := config.load()
	if err != nil {
		return err
	}

	return rootCmd.Execute()
}

type configuration struct {
	BinaryDirectory        string
	CacheDirectory         string
	ConfigDirectory        string
	CurrentVersion         string
	HomeDirectory          string
	InstalledVersions      []string
	RepositoryDomain       string
	TempDirectory          string
	TerraformSymlinkTarget string
}

func (c *configuration) load() error {
	if err := c.binDir(); err != nil {
		return err
	}

	if err := c.confDir(); err != nil {
		return err
	}

	if err := c.cacheDir(); err != nil {
		return err
	}

	if err := c.tmpDir(); err != nil {
		return err
	}

	if err := c.currentVersion(); err != nil {
		return err
	}

	c.TerraformSymlinkTarget = filepath.Join(c.BinaryDirectory, terraform)
	c.RepositoryDomain = "releases.hashicorp.com"

	return nil
}

// binDir sets the directory the Terraform binaries will be symlinked
// to. Defaults to ${HOME}/bin
func (c *configuration) binDir() error {
	if err := c.homeDir(); err != nil {
		return err
	}

	c.BinaryDirectory = filepath.Join(c.HomeDirectory, "bin")

	return nil
}

// cacheDir sets the directory to sore temporary files e.g. downloads,
// available remote versions. Defaults to the following:
//
//	UNIX: ${HOME}/.cache/tfsw
//	Windows: TBC
func (c *configuration) cacheDir() error {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	c.CacheDirectory = filepath.Join(userCacheDir, "tfsw")

	return nil
}

// confDir sets the directory to store permanent files e.g. Terraform
// versions that have been downloaded, and configuration. Defaults to
// the following:
//
//	UNIX: ${HOME}/.config/tfsw
//	Windows: TBC
func (c *configuration) confDir() error {
	userConfDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	c.ConfigDirectory = filepath.Join(userConfDir, "tfsw")

	return nil
}

// homeDir returns the configured home directory. Defaults to the
// users ${HOME}
func (c *configuration) homeDir() error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	c.HomeDirectory = userHomeDir

	return nil
}

// tmpDir a temporary directory to store transient files. Defaults to
// ${cacheDir}/tmp
func (c *configuration) tmpDir() error {
	if c.CacheDirectory == "" {
		return fmt.Errorf("asdsda")
	}

	c.TempDirectory = filepath.Join(c.CacheDirectory, "tmp")
	return nil
}

// currentVersion finds the active version of Terraform and adds
// it to the configuration struct as CurrentVersion
func (c *configuration) currentVersion() error {
	if err := c.installedVersions(); err != nil {
		return err
	}

	if len(c.InstalledVersions) == 0 {
		c.CurrentVersion = ""
		return nil
	}

	link, err := filepath.EvalSymlinks(filepath.Join(c.BinaryDirectory, terraform))
	if err != nil {
		return err
	}

	var version string
	for _, i := range c.InstalledVersions {
		if strings.Contains(link, i) {
			version = i
			break
		}
	}

	c.CurrentVersion = version
	return nil
}

// installedVersions finds the currently installed versions of Terraform
// and adds them to the configuration struct in InstalledVersions
func (c *configuration) installedVersions() error {
	fh, err := os.Open(c.ConfigDirectory)
	if err != nil {
		// NOTE: If this is the first time tfsw is being used the
		// config dir won't exist yet
		if errors.Is(err, os.ErrNotExist) {
			c.InstalledVersions = nil
			return nil
		}
		return err
	}

	defer fh.Close()

	dirs, err := fh.ReadDir(0)
	if err != nil {
		return err
	}

	var iv []string
	for _, d := range dirs {
		if d.IsDir() {
			if regex.MatchString(d.Name()) {
				iv = append(iv, d.Name())
			}
		}
	}

	sort.Strings(iv)
	c.InstalledVersions = iv
	return nil
}

// validateVersion is used by commands to put some guard rails around
// the version of Terraform we're downloading. It reads through any
// arguments and validates them against a regex
func validateVersion(cmd *cobra.Command, args []string) {
	for i := range args {
		if !regex.MatchString(args[i]) {
			fmt.Fprintf(os.Stderr, "%s is not a valid version\n", args[i])
			os.Exit(1)
		}
	}
}
