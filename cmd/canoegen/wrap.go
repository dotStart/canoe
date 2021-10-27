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
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/dotstart/canoe/build"
	"github.com/dotstart/canoe/internal"
	"github.com/dotstart/canoe/internal/metadata"
	"github.com/google/subcommands"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const defaultRuntimeVersion = 16

var targetPattern = regexp.MustCompile("^[a-z0-9_-]+$")

type wrapCommand struct {
	inputFile  string
	outputFile string
	mainClass  string

	target        string
	wrapperFile   string
	useGuiWrapper bool

	runtimeMinimumVersion uint
	runtimeMaximumVersion uint
	runtimeInitialMemory  string
	runtimeMemoryLimit    string
	runtimeArguments      string

	verbose bool
}

func (*wrapCommand) Name() string {
	return "wrap"
}

func (*wrapCommand) Synopsis() string {
	return "creates a canoe wrapped Java archive"
}

func (*wrapCommand) Usage() string {
	wrappers, _ := build.GetFilesystem().ReadDir("wrappers")

	targets := ""
	for _, target := range wrappers {
		if target.IsDir() {
			targets += "  - " + target.Name() + "\n"
		}
	}

	return `canoegen wrap -in <file> [-out <file>] [-target <name>] [args]

Generates a canoe self contained executable which automatically locates compatible runtime 
installations to launch a given embedded Java application.

In order to generate a simple self contained executable the following command may be used:

  $ canoegen wrap -in foo.jar

Where "foo.jar" contains the application code which shall be wrapped. By default, Linux, Mac OS and
Windows versions of the executable will be placed in the current working directory.

Alternatively, a target directory may be specified via the "-out" parameter:

  $ canoegen wrap -in foo.jar -out ./target

The wrap subcommand may also be used to explicitly generate an image for a given target 
platform via the "-target" parameter:

  $ canoegen wrap -in foo.jar -target windows-amd64 -out foo.exe

When a single target is given, the "-out" parameter specifies the full file path for the resulting
executable.

The following targets are available via the "-target" option:

` + targets + `
If desired, the built-in targets may be replaced by custom executables via the "-wrapper" option:

  $ canoegen wrap -in foo.jar -wrapper mywrapper.exe -out foo.exe

Custom wrappers are expected to rely on the implementations provided by the canoew-cli or canoew-gui
packages respectively and will _NOT_ be validated by this tool. Please ensure that passed wrapper
executables are actually compatible with this revision of the tool as wrapped executables may 
otherwise fail to launch or produce other undesired side effects.

The following configuration options are provided by this command:

`
}

func (cmd *wrapCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.inputFile, "in", "", "selects an input archive (required)")
	f.StringVar(&cmd.outputFile, "out", ".", "selects an output file or directory")
	f.StringVar(&cmd.mainClass, "main-class", "", "selects a specific main class to launch (defaults to the Main-Class attribute within the archive manifest)")

	f.StringVar(&cmd.target, "target", "", "selects a target platform (defaults to all)")
	f.StringVar(&cmd.wrapperFile, "wrapper", "", "selects an alternative wrapper executable (defaults to embedded executables)")
	f.BoolVar(&cmd.useGuiWrapper, "gui", false, "selects a GUI focused wrapper executable on supported platforms (only applies to Windows targets; ignored otherwise)")

	f.UintVar(&cmd.runtimeMinimumVersion, "runtime-version", defaultRuntimeVersion, fmt.Sprintf("defines the minimum required runtime version (defaults to %d)", defaultRuntimeVersion))
	f.UintVar(&cmd.runtimeMaximumVersion, "runtime-max-version", 0, "defines the maximum permitted runtime version (unset by default)")
	f.StringVar(&cmd.runtimeInitialMemory, "runtime-initial-memory", "", "defines the initial runtime memory (unset by default)")
	f.StringVar(&cmd.runtimeMemoryLimit, "runtime-memory-limit", "", "defines the runtime memory limit (unset by default)")
	f.StringVar(&cmd.runtimeArguments, "runtime-args", "", "supplies additional arguments to be passed to the runtime upon application startup")

	f.BoolVar(&cmd.verbose, "verbose", false, "prints additional information when generating executables")
}

func (cmd *wrapCommand) Execute(context.Context, *flag.FlagSet, ...interface{}) subcommands.ExitStatus {
	if len(cmd.inputFile) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "invalid parameters: input file is required")
		return subcommands.ExitUsageError
	}

	archive, err := ioutil.ReadFile(cmd.inputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read input file: %s\n", err)
		return subcommands.ExitFailure
	}

	inputBase := filepath.Base(cmd.inputFile)
	extensionOffset := strings.LastIndex(inputBase, ".")
	inferredOutputName := inputBase[:extensionOffset]

	runtimeInitialMemory := uint64(0)
	if len(cmd.runtimeInitialMemory) != 0 {
		runtimeInitialMemory, err = metadata.ParseByteSuffix(cmd.runtimeInitialMemory)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "invalid initial memory size: %s\n", err)
			return subcommands.ExitUsageError
		}
	}

	runtimeMemoryLimit := uint64(0)
	if len(cmd.runtimeMemoryLimit) != 0 {
		runtimeMemoryLimit, err = metadata.ParseByteSuffix(cmd.runtimeMemoryLimit)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "invalid memory limit: %s\n", err)
			return subcommands.ExitUsageError
		}
	}

	meta := &metadata.ApplicationContainer{
		CanoeVersion:  internal.Version(),
		CustomWrapper: len(cmd.wrapperFile) != 0,
		Runtime: &metadata.RuntimeConfiguration{
			MinimumVersion:      uint64(cmd.runtimeMinimumVersion),
			MaximumVersion:      uint64(cmd.runtimeMaximumVersion),
			InitialMemory:       runtimeInitialMemory,
			MemoryLimit:         runtimeMemoryLimit,
			AdditionalArguments: cmd.runtimeArguments,
		},
		Application: &metadata.ApplicationConfiguration{
			MainClass: cmd.mainClass,
		},
	}

	if len(cmd.target) == 0 && len(cmd.wrapperFile) == 0 {
		targets, err := build.GetFilesystem().ReadDir("wrappers")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to load target list: %s\n", err)
		}

		for _, target := range targets {
			if !target.IsDir() {
				continue
			}

			if cmd.verbose {
				fmt.Printf("generating target %s\n", target.Name())
			}

			outputFileName := inferredOutputName + "-" + target.Name()
			output := filepath.Join(cmd.outputFile, outputFileName)

			if err := cmd.generateFromTarget(meta, target.Name(), archive, output); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "failed to generate target %s: %s\n", target.Name(), err)
			}
		}

		return subcommands.ExitSuccess
	}

	output := cmd.outputFile
	if output == "." {
		output = inferredOutputName
	}

	if len(cmd.wrapperFile) == 0 {
		if !targetPattern.MatchString(cmd.target) {
			_, _ = fmt.Fprintf(os.Stderr, "invalid target: %s\n", cmd.target)
			return subcommands.ExitUsageError
		}

		err = cmd.generateFromTarget(meta, cmd.target, archive, output)
	} else {
		err = cmd.generateFromExecutable(meta, cmd.wrapperFile, archive, output)
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to generate executable: %s\n", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (cmd *wrapCommand) generateFromExecutable(meta *metadata.ApplicationContainer, input string, archive []byte, output string) error {
	inFile, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("failed to open wrapper %s: %w", input, err)
	}

	return cmd.generate(meta, inFile, archive, output)
}

func (cmd *wrapCommand) generateFromTarget(meta *metadata.ApplicationContainer, target string, archive []byte, output string) error {
	filename := "canoew"
	if strings.Contains(target, "windows") {
		if cmd.useGuiWrapper {
			filename = "canoew-gui"
		}

		filename += ".exe"
		output += ".exe"
	}

	inFile, err := build.GetFilesystem().ReadFile("wrappers/" + target + "/" + filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("invalid target: %s", target)
		}

		return fmt.Errorf("failed to open wrapper for target %s: %w", target, err)
	}

	return cmd.generate(meta, inFile, archive, output)
}

func (cmd *wrapCommand) generate(meta *metadata.ApplicationContainer, wrapper []byte, archive []byte, output string) error {
	parent := filepath.Dir(output)
	if _, err := os.Stat(parent); err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(parent, 0755); err != nil {
			return fmt.Errorf("failed to create directory structure at %s: %w", parent, err)
		}
	}

	outFile, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("failed to open output file %s: %w", output, err)
	}

	if _, err := outFile.Write(wrapper); err != nil {
		return fmt.Errorf("failed to create output file %s: %w", output, err)
	}

	if _, err := outFile.Write(archive); err != nil {
		return fmt.Errorf("failed to write to output file: %s: %w", output, err)
	}

	if _, err := internal.WriteExecutableFooter(outFile, meta); err != nil {
		return fmt.Errorf("failed to finalize output file %s: %w", output, err)
	}

	return nil
}
