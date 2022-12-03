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

	"github.com/jedib0t/go-pretty/v6/table"
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
	releaseNotesURL = "https://github.com/hashicorp/terraform/releases/tag/v%s"
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
// active version, and prints out a pretty table
func listVersions(inst []string, cur string) error {
	if t := len(inst); t == 0 {
		return ErrNoneInstalled
	}

	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Version", "Active", "Release Notes"})
	for _, v := range inst {
		rn := fmt.Sprintf(releaseNotesURL, v)

		if v == cur {
			tw.AppendRow(table.Row{v, "true", rn})
			continue
		}

		tw.AppendRow(table.Row{v, "", rn})
	}

	fmt.Println(tw.Render())
	return nil
}
