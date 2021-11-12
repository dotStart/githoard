package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/dotstart/githoard/internal/hoard"
	"github.com/dotstart/githoard/internal/service"
	"github.com/google/subcommands"
	"net/url"
	"os"
)

type repoCommand struct {
	registry *service.Registry

	ownerName      string
	repositoryName string

	force bool
}

func RepoCommand(registry *service.Registry) subcommands.Command {
	return &repoCommand{
		registry: registry,
	}
}

func (*repoCommand) Name() string {
	return "repo"
}

func (*repoCommand) Synopsis() string {
	return "creates a mirror of a given repository"
}

func (*repoCommand) Usage() string {
	return `repo [options] <uri>

Creates a mirror for a given target repository. For example:

  $ githoard repo https://github.com/dotStart/Beacon

Alternatively, you may also specify the owner of the newly created repository if it should differ
from the original:

  $ githoard repo -owner=foo https://github.com/dotStart/Beacon.git

When mirroring repositories from supported providers, their information will be copied to the 
target.

All available command line options:

`
}

func (cmd *repoCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.ownerName, "owner", "", "specifies the target repository owner")
	f.StringVar(&cmd.repositoryName, "repo", "", "specifies the target repository name")

	f.BoolVar(&cmd.force, "force", false, "causes existing repositories with the same name to be purged first")
}

func (cmd *repoCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 1 {
		_, _ = fmt.Fprintln(os.Stderr, "invalid command line: expected repository url")
		return subcommands.ExitUsageError
	}

	repoUrl, err := url.Parse(f.Arg(0))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "invalid repository URL: %s\n", err)
		return subcommands.ExitUsageError
	}

	h, err := hoard.New(cmd.registry)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		return subcommands.ExitFailure
	}

	err = h.MirrorRepo(ctx, repoUrl, hoard.MirrorRepoOptions{
		MigrationOptions: hoard.MigrationOptions{
			Force: cmd.force,
		},
		OwnerName:      cmd.ownerName,
		RepositoryName: cmd.repositoryName,
	})
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
