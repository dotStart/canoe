/*
 * Copyright 2021 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package runtime

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const cmdVersionPrefix = "version "

// attempts to locate a given Java executable with the desired version number
func FindInPath(executableName string, minimumVersion uint64, maximumVersion uint64) error {
	cmd := exec.Command(executableName, "-version")

	pipe, err := cmd.StderrPipe()
	if err != nil {
		return ErrNotFound
	}

	scanner := bufio.NewScanner(pipe)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%w: failed to launch Java process", ErrInvalidInstallation)
	}

	if !scanner.Scan() {
		return fmt.Errorf("%w: Java process did not provide version information", ErrInvalidInstallation)
	}

	versionLine := scanner.Text()
	versionOffset := strings.Index(versionLine, cmdVersionPrefix)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%w: Java process terminated abnormally", ErrInvalidInstallation)
	}

	if versionOffset == -1 {
		return fmt.Errorf("%w: Java process did not provide valid version information: missing version number", ErrInvalidInstallation)
	}

	versionNumber := versionLine[(versionOffset + len(cmdVersionPrefix)):]
	if versionNumber[0] == '"' {
		versionNumber = versionNumber[1:]
	}

	majorSeparator := strings.IndexRune(versionNumber, '.')
	if majorSeparator == -1 {
		return fmt.Errorf("%w: Java process did not provide valid version information", ErrInvalidInstallation)
	}

	majorString := versionNumber[:majorSeparator]
	majorNumber, err := strconv.ParseUint(majorString, 10, 32)
	if err != nil {
		return fmt.Errorf("%w: Java process did not provide valid version information (%s)", ErrInvalidInstallation, err)
	}

	if majorNumber < minimumVersion {
		return fmt.Errorf("%w: %d required (%d found)", ErrUnsupported, minimumVersion, majorNumber)
	}
	if maximumVersion != 0 && majorNumber > maximumVersion {
		return fmt.Errorf("%w: %d and newer are unsupported (%d found)", ErrUnsupported, maximumVersion, majorNumber)
	}

	return nil
}
