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
	"github.com/dotstart/canoe/internal/metadata"
	"github.com/google/subcommands"
	"os"
)

type infoCommand struct {
	inputFile string
}

func (*infoCommand) Name() string {
	return "info"
}

func (*infoCommand) Synopsis() string {
	return "decodes wrapper metadata from a canoe executable"
}

func (*infoCommand) Usage() string {
	return `canoegen info -in <file> [args]

Displays the configuration information stored within a given canoe executable which has previously
been generated using canoe wrap:

  $ canoegen info -in foo.exe

The following configuration options are provided by this command:

`
}

func (cmd *infoCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.inputFile, "in", "", "selects an input executable")
}

func (cmd *infoCommand) Execute(context.Context, *flag.FlagSet, ...interface{}) subcommands.ExitStatus {
	if len(cmd.inputFile) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "invalid parameters: input file is required")
		return subcommands.ExitUsageError
	}

	meta, err := internal.ReadExecutableFooter(cmd.inputFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read executable: %s\n", err)
		return subcommands.ExitFailure
	}

	fmt.Printf("==== configuration of %s ====\n\n", cmd.inputFile)

	fmt.Println("==> canoe metadata")
	fmt.Println()
	fmt.Printf("          version: %s\n", meta.CanoeVersion)
	fmt.Printf(" custom generator: %v\n", meta.CustomWrapper)
	fmt.Println()

	fmt.Println("==> runtime configuration")
	fmt.Println()
	fmt.Printf("      minimum version: %d\n", meta.Runtime.MinimumVersion)
	fmt.Printf("      maximum version: %d\n", meta.Runtime.MaximumVersion)
	fmt.Println()

	fmt.Printf("       initial memory: %s\n", metadata.AppendByteSuffix(meta.Runtime.InitialMemory))
	fmt.Printf("         memory limit: %s\n", metadata.AppendByteSuffix(meta.Runtime.MemoryLimit))
	fmt.Printf(" additional arguments: \"%s\"\n", meta.Runtime.AdditionalArguments)
	fmt.Println()

	fmt.Println("==> application configuration")
	fmt.Println()
	fmt.Printf(" main class: %s\n", meta.Application.MainClass)
	fmt.Println()

	fmt.Println("-- end of readout --")

	return subcommands.ExitSuccess
}
