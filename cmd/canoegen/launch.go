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
	"flag"
	"fmt"
	"github.com/dotstart/canoe/internal"
	"github.com/dotstart/canoe/internal/runtime"
	"github.com/google/subcommands"
	"os"
)

type launchCommand struct {
	inputFile string
}

func (*launchCommand) Name() string {
	return "launch"
}

func (*launchCommand) Synopsis() string {
	return "launches the wrapped application within a canoe executable"
}

func (*launchCommand) Usage() string {
	return `canoegen launch -in <file> [args]

Launches the archive contained within a given canoe wrapped executable:

  $ canoegen launch -in foo.exe

This command is primarily provided for development purposes.

The following configuration options are provided by this command:

`
}

func (cmd *launchCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.inputFile, "in", "", "selects an input executable")
}

func (cmd *launchCommand) Execute(context.Context, *flag.FlagSet, ...interface{}) subcommands.ExitStatus {
	if len(cmd.inputFile) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "invalid parameters: input file is required")
		return subcommands.ExitUsageError
	}

	return subcommands.ExitStatus(internal.LaunchApplication(cmd.inputFile, runtime.CliExecutableName))
}
