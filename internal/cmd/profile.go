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

type profileCommand struct {
	registry *service.Registry

	ownerName string

	force bool
}

func ProfileCommand(registry *service.Registry) subcommands.Command {
	return &profileCommand{
		registry: registry,
	}
}

func (*profileCommand) Name() string {
	return "profile"
}

func (*profileCommand) Synopsis() string {
	return "creates an archive of all repositories owned by a given profile"
}

func (*profileCommand) Usage() string {
	return `profile [options] <uri>

Creates a mirror for a given target profile. For example:

  $ githoard profile https://github.com/dotStart

Alternatively, you may also specify the owner of the newly created repositories if they should differ
from the original:

  $ githoard profile -owner=foo https://github.com/dotStart

All available command line options:

`
}

func (cmd *profileCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.ownerName, "owner", "", "specifies the target repository owner")

	f.BoolVar(&cmd.force, "force", false, "causes existing repositories with the same name to be purged first")
}

func (cmd *profileCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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

	err = h.MirrorProfile(ctx, repoUrl, hoard.MirrorProfileOptions{
		MigrationOptions: hoard.MigrationOptions{
			Force: cmd.force,
		},
		OwnerName: cmd.ownerName,
	})
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
