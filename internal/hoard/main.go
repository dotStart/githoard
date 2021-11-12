package hoard

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"errors"
	"fmt"
	"github.com/dotstart/githoard/internal/config"
	"github.com/dotstart/githoard/internal/service"
	"net/url"
	"os"
	"strings"
)

type Hoard struct {
	registry *service.Registry
	login    *config.Login

	upstream *gitea.Client
	usr      *gitea.User
}

func New(registry *service.Registry) (*Hoard, error) {
	login, err := config.ReadLogin()
	if err != nil {
		return nil, fmt.Errorf("failed to load login information: %s", err)
	}
	if login == nil {
		return nil, errors.New("login required")
	}

	client, err := gitea.NewClient(login.InstanceUri, gitea.SetToken(login.Token))
	if err != nil {
		return nil, fmt.Errorf("failed to contact upstream: %s", err)
	}

	usr, _, err := client.GetMyUserInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user info: %w", err)
	}

	return &Hoard{
		registry: registry,
		login:    login,
		upstream: client,
		usr:      usr,
	}, nil
}

func (h *Hoard) doMirrorRepo(migration gitea.MigrateRepoOption, opts MigrationOptions) error {
	if strings.ToLower(h.usr.UserName) != strings.ToLower(migration.RepoOwner) {
		_, resp, err := h.upstream.GetOrg(migration.RepoOwner)
		if err != nil {
			if resp == nil || resp.StatusCode != 404 {
				return fmt.Errorf("failed to retrieve organization info: %w", err)
			}

			org, _, err := h.upstream.CreateOrg(gitea.CreateOrgOption{
				Name: migration.RepoOwner,
			})
			if err != nil {
				return fmt.Errorf("failed to create organization: %w", err)
			}

			fmt.Printf("created organization %s#%d\n", org.UserName, org.ID)
		}
	}

	if opts.Force {
		existingRepo, resp, err := h.upstream.GetRepo(migration.RepoOwner, migration.RepoName)
		if err != nil {
			if resp == nil || resp.StatusCode != 404 {
				return fmt.Errorf("failed to verify existing repository: %s", err)
			}
		} else {
			_, err := h.upstream.DeleteRepo(migration.RepoOwner, migration.RepoName)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "failed to delete existing repository: %s\n", err)
			}

			fmt.Printf("deleted existing repository %s#%d (-force)\n", existingRepo.FullName, existingRepo.ID)
		}
	}

	repo, _, err := h.upstream.MigrateRepo(migration)
	if err != nil {
		return fmt.Errorf("failed to create migration: %s", err)
	}

	fmt.Printf("created repository %s#%d - %s\n", repo.FullName, repo.ID, repo.HTMLURL)
	return nil
}

func (h *Hoard) MirrorRepo(ctx context.Context, repoUrl *url.URL, opts MirrorRepoOptions) error {
	migration, err := h.registry.RetrieveRepoOptions(ctx, repoUrl, h.login)
	if err != nil {
		return fmt.Errorf("failed to retrieve migration options: %w", err)
	}

	if migration == nil {
		// TODO
		return fmt.Errorf("unknown service provider - manually select repository name and owner")
	}

	if opts.OwnerName != "" {
		migration.RepoOwner = opts.OwnerName
	}
	if opts.RepositoryName != "" {
		migration.RepoName = opts.RepositoryName
	}

	return h.doMirrorRepo(*migration, opts.MigrationOptions)
}

func (h *Hoard) MirrorProfile(ctx context.Context, profileUrl *url.URL, opts MirrorProfileOptions) error {
	migrations, err := h.registry.RetrieveProfileOptions(ctx, profileUrl, h.login)
	if err != nil {
		return fmt.Errorf("failed to retrieve migration options: %w", err)
	}

	if migrations == nil {
		return fmt.Errorf("unknown service provider")
	}

	for _, migration := range migrations {
		if opts.OwnerName != "" {
			migration.RepoOwner = opts.OwnerName
		}

		if err := h.doMirrorRepo(migration, opts.MigrationOptions); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to mirror %s/%s: %s\n", migration.RepoOwner, migration.RepoName, err)
		}
	}

	return nil
}
