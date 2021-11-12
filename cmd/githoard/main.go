package main

import (
	"context"
	"flag"
	"github.com/dotstart/githoard/internal/cmd"
	"github.com/dotstart/githoard/internal/service"
	"github.com/google/subcommands"
	"os"
)

func main() {
	registry := service.NewRegistry()
	registry.Register(service.NewGitHubProvider())

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(cmd.LoginCommand(), "")
	subcommands.Register(cmd.ProfileCommand(registry), "")
	subcommands.Register(cmd.RepoCommand(registry), "")
	subcommands.Register(cmd.VersionCommand(), "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
