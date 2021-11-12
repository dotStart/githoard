package cmd

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"flag"
	"fmt"
	"github.com/dotstart/githoard/internal/config"
	"github.com/google/subcommands"
	"os"
)

type loginCommand struct {
	githubToken string
}

func LoginCommand() subcommands.Command {
	return &loginCommand{}
}

func (*loginCommand) Name() string {
	return "login"
}

func (*loginCommand) Synopsis() string {
	return "authenticates with gitea"
}

func (*loginCommand) Usage() string {
	return `login <uri> <token>

Configures and authenticates against a given gitea instance. For example:

  $ githoard login https://gitea.example.org abcdef01234

Available command line options:

`
}

func (cmd *loginCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.githubToken, "github-token", "", "specifies a GitHub token")
}

func (cmd *loginCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 2 {
		_, _ = fmt.Fprintln(os.Stderr, "invalid command line: Must specify gitea URI and token")
		return subcommands.ExitUsageError
	}

	rootUri := f.Arg(0)
	token := f.Arg(1)

	client, err := gitea.NewClient(rootUri, gitea.SetToken(token))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to establish connection: %s\n", err)
		return subcommands.ExitFailure
	}

	usr, _, err := client.GetMyUserInfo()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to retrieve user information: %s\n", err)
		return subcommands.ExitFailure
	}

	cfg := &config.Login{
		InstanceUri: rootUri,
		Token:       token,
		GitHubToken: cmd.githubToken,
	}

	if err := cfg.Write(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to write configuration file: %s\n", err)
		return subcommands.ExitFailure
	}

	fmt.Printf("authenticated with %s as %s#%d\n", rootUri, usr.UserName, usr.ID)
	return subcommands.ExitSuccess
}
