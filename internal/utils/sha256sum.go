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
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// TODO - Verify SHA256SUMS against GPG signature
func Sha256sum(sums, file string) (bool, error) {
	// TODO? - Emulate the coreutils sha256sum functionaliy where it will attempt to
	// find and vaildate all files in the sumfile

	lines, err := fileLines(sums)
	if err != nil {
		return false, err
	}

	base := filepath.Base(file)
	var sum string
	for _, l := range lines {
		l := strings.Fields(l)
		if len(l) != 2 {
			return false, errors.New(sums + " format invalid, a line does not have two elements")
		}

		if strings.TrimPrefix(l[1], "*") == base {
			sum = l[0]
			break
		}
	}

	if sum == "" {
		return false, errors.New(base + " is not in " + sums)
	}

	sha256sum, err := hashFile(file)
	if err != nil {
		return false, err
	}

	if sum == sha256sum {
		return true, nil
	}

	return false, nil
}

// TODO - Allow sha256.New() to be passed in as an argument so it can take in any hashing algo
// as long as it implements hash.Hash
func hashFile(file string) (sum string, err error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buffer := make([]byte, 32*1024)
	hash := sha256.New()
	for {
		n, err := f.Read(buffer)
		if n > 0 {
			_, err := hash.Write(buffer[:n])
			if err != nil {
				return "", err
			}
			continue
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// fileLines reads a file and returns the lines as a slice
// of strings
func fileLines(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
