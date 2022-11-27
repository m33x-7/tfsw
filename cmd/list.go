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
)

var (
	listCmd = &cobra.Command{
		Aliases: []string{"ls"},
		Long:    "Lists all currently installed versions in the cache, and marks the active version",
		Run:     listRun,
		Short:   "List installed versions",
		Use:     "list",
	}
)

func init() {
	// Add list as a child command of tfsw
	rootCmd.AddCommand(listCmd)

	// Add any extra command line flags for list here
	listCmd.Flags().BoolP("remote", "r", false, "list ")
}

// listRun is passed directly to the Cobra Run argument and executes
// the primary logic for the `list` command
func listRun(cmd *cobra.Command, args []string) {
	err := listVersions(config.InstalledVersions, config.CurrentVersion)
	switch err {
	case ErrNoneInstalled:
		fmt.Printf("No versions of Terraform have been installed with %s\n", basename)
		os.Exit(0)
	case nil:
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Encountered an unhandled error: %v\n", err)
		os.Exit(1)
	}
}

// listVersions takes the currently installed versions, and the current
// active version, and prints out a pretty list
func listVersions(inst []string, cur string) error {
	total := len(inst)
	if total == 0 {
		return ErrNoneInstalled
	}

	enabled := cur != ""

	fmt.Printf("| %s | installed: %d | enabled: %t\n", filepath.Base(os.Args[0]), total, enabled)
	for i, v := range inst {
		if maxIdx := total - 1; i == maxIdx {
			if v == cur {
				fmt.Printf("└ %s *\n", v)
				break
			}
			fmt.Printf("└ %s\n", v)
			break
		}

		if v == cur {
			fmt.Printf("├ %s *\n", v)
			continue
		}

		fmt.Printf("├ %s\n", v)
	}

	return nil
}
