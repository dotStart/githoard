package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/dotstart/githoard/internal"
	"github.com/google/subcommands"
)

type versionCommand struct {
}

func VersionCommand() subcommands.Command {
	return &versionCommand{}
}

func (*versionCommand) Name() string {
	return "version"
}

func (*versionCommand) Synopsis() string {
	return "retrieves the application version"
}

func (*versionCommand) Usage() string {
	return `version

Exposes the current application version number:

  $ githoard version

All available command line options:

`
}

func (*versionCommand) SetFlags(*flag.FlagSet) {
}

func (*versionCommand) Execute(context.Context, *flag.FlagSet, ...interface{}) subcommands.ExitStatus {
	fmt.Printf("githoard v%s\n", internal.FullVersion())
	fmt.Println()
	if internal.HasCommitHash() {
		fmt.Printf("Commit Hash: %s\n", internal.CommitHash())
		fmt.Printf("Build Date: %s\n", internal.BuildTimestamp())
	} else {
		fmt.Println("Build metadata is missing - This is a development build")
	}

	return subcommands.ExitSuccess
}
