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
	"path/filepath"

	"github.com/spf13/cobra"
	"tfsw/internal/utils"
)

var (
	deleteCmd = &cobra.Command{
		Aliases:           []string{"rm"},
		Long:              "Delete an installed version from the cache",
		PreRun:            validateVersion,
		Run:               deleteRun,
		Short:             "Delete a version",
		Use:               "delete {VERSION... | --clean}",
		ValidArgsFunction: deleteValidArgs,
	}
)

func init() {
	// Adds delete as a child command of tfsw
	rootCmd.AddCommand(deleteCmd)

	// Add any extra command line flags for remove here
	deleteCmd.Flags().BoolP("clean", "c", false, "Remove all but the active version")
}

// deleteRun is passed directly to the Cobra Run argument and executes
// the primary logic for the `delete` command
func deleteRun(cmd *cobra.Command, args []string) {
	clean, _ := cmd.Flags().GetBool("clean")
	if clean {
		args = deleteCleanArgs(config.InstalledVersions, config.CurrentVersion)
	}

	for _, arg := range args {
		err := deleteVersion(arg, config.CurrentVersion)
		switch err {
		case ErrVersionSame:
			fmt.Printf("Terraform %s is active, please switch to another version before removing\n", arg)
		case ErrVersionNotExist:
			fmt.Printf("Terraform %s has already been removed\n", arg)
		case nil:
			fmt.Printf("Terraform %s has been removed\n", arg)
		default:
			fmt.Fprintf(os.Stderr, "Encountered an unhandled error: %v\n", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}

// deleteCleanArgs returns a slice of the installed versions, minus the current
// version, to be deleted
func deleteCleanArgs(inst []string, cur string) []string {
	idx := utils.Index(inst, cur)
	return utils.Reslice(inst, idx)
}

// deleteValidArgs dynamically generates completion arguments for the remove
// command. It gathers a list of the currently installed versions, then
// modifies that list based on the current command line input:
//
//	$ tfsw delete [tab][tab]
//	1.0.0 0.15.5 0.14.7
//
//	$ tfsw delete 0.14.7 [tab][tab]
//	1.0.0 0.15.5
func deleteValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var validArgs []string

	for i := range args {
		index := utils.Index(versionsInstalled, args[i])
		validArgs = utils.Reslice(versionsInstalled, index)
	}

	return validArgs, cobra.ShellCompDirectiveNoFileComp
}

// deleteVersion takes a Terraform version unumber, and the current active
// version, and deletes the version if they don't match.
func deleteVersion(ver, cur string) error {
	if ver == cur {
		return ErrVersionSame
	}

	dst := filepath.Join(config.ConfigDirectory, ver)

	if _, err := os.Stat(dst); err != nil {
		return ErrVersionNotExist
	}

	if err := os.RemoveAll(dst); err != nil {
		return err
	}

	return nil
}
