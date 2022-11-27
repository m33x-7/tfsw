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

	"github.com/spf13/cobra"
)

var (
	buildVer    string
	buildCommit string
	versionCmd  = &cobra.Command{
		Use:   "version",
		Long:  "Print verion, commit, etc about tfsw",
		Run:   versionRun,
		Short: "Print version information",
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// versionRun is passed directly to the Cobra Run argument and executes
// the primary logic for the `version` command
func versionRun(cmd *cobra.Command, args []string) {
	fmt.Printf("Terraform Switch %s\n", buildVer)
	return
}
