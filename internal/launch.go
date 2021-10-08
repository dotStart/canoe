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
package internal

import (
	"errors"
	"fmt"
	"github.com/dotstart/canoe/internal/metadata"
	"github.com/dotstart/canoe/internal/runtime"
	"github.com/gen2brain/dlgs"
	"os"
	"os/exec"
	"path"
	"strings"
)

func Launch(runtimeExecutable string) int {
	executable, err := os.Executable()
	if err != nil {
		_, _ = dlgs.Error("Application Error", "Failed to open application executable.")
		return -1
	}

	return LaunchApplication(executable, runtimeExecutable)
}

func LaunchApplication(executable string, runtimeExecutable string) int {
	cfg, err := ReadExecutableFooter(executable)
	if err != nil {
		_, _ = dlgs.Error("Application Error", "Failed to load application configuration.")
		return -2
	}

	home, err := runtime.Find(cfg.Runtime.MinimumVersion, cfg.Runtime.MaximumVersion)
	if errors.Is(err, runtime.ErrNotFound) {
		home = ""
		err = runtime.FindInPath(runtimeExecutable, cfg.Runtime.MinimumVersion, cfg.Runtime.MaximumVersion)
	}
	if err != nil {
		_, _ = dlgs.Error("Runtime Error", fmt.Sprintf("Failed to locate valid Java Runtime: %s", err))
		return -3
	}

	executablePath := path.Join(home, runtimeExecutable)
	if _, err := os.Stat(executablePath); err != nil {
		_, _ = dlgs.Error("Runtime Error", "Invalid Java Runtime installation: Cannot find executable")
		return -4
	}

	arguments := make([]string, 0)

	if cfg.Runtime.InitialMemory != 0 {
		arguments = append(arguments, "-Xms"+metadata.AppendByteSuffix(cfg.Runtime.InitialMemory))
	}
	if cfg.Runtime.MemoryLimit != 0 {
		arguments = append(arguments, "-Xmx"+metadata.AppendByteSuffix(cfg.Runtime.MemoryLimit))
	}

	if len(cfg.Runtime.AdditionalArguments) != 0 {
		arguments = append(arguments, strings.Split(cfg.Runtime.AdditionalArguments, " ")...)
	}

	arguments = append(arguments, "-cp", executable)
	arguments = append(arguments, cfg.Application.MainClass)

	cmd := exec.Command(executablePath, arguments...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}

		return -5
	}

	return 0
}
