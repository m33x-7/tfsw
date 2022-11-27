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

	"github.com/spf13/cobra"
	"tfsw/internal/utils"
)

var (
	selectCmd = &cobra.Command{
		Args:              cobra.ExactValidArgs(1),
		Long:              "Select the active Terraform version and install it if missing",
		PreRun:            validateVersion,
		Run:               selectRun,
		Short:             "Select the active version",
		Use:               "select VERSION...",
		ValidArgsFunction: selectValidArgs,
	}
)

func init() {
	rootCmd.AddCommand(selectCmd)
}

// selectRun is passed directly to the Cobra Run argument and executes
// the primary logic for the `select` command
func selectRun(cmd *cobra.Command, args []string) {
	err := selectVersion(config.CurrentVersion, args[0])
	switch err {
	case ErrVersionSame:
		fmt.Printf("Terraform %s already active!\n", args[0])
		os.Exit(0)
	case nil:
		fmt.Printf("Terraform %s is now active\n", args[0])
		os.Exit(0)
	default:
		// NOTE: This doesn't need a trailing \n to be set
		fmt.Fprintf(os.Stderr, "Error setting Terraform version: %v", err)
		os.Exit(1)
	}
}

// selectVaildArgs only allows completion on the first argument as only
// one argument is accepted
func selectValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) >= 1 {
		return []string{}, cobra.ShellCompDirectiveNoFileComp
	}

	validArgs := config.InstalledVersions

	return validArgs, cobra.ShellCompDirectiveNoFileComp
}

// selectVersion takes the current version, and a new version. If they're
// the same it informs the user. If it's missing, it downloads it. Once the
// version is available it updates the symlink
func selectVersion(cur, new string) error {
	if new == cur {
		return ErrVersionSame
	}

	if err := newVersion(new); err != nil && !errors.Is(err, ErrVersionExists) {
		return err
	}

	tgt := filepath.Join(config.BinaryDirectory, terraform)
	src := filepath.Join(config.ConfigDirectory, new, terraform)
	if err := utils.Symlink(src, tgt); err != nil {
		return err
	}

	return nil
}
