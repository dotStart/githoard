package service

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"fmt"
	"github.com/dotstart/githoard/internal/config"
	"net/url"
	"strings"
)

type Registry struct {
	providers map[string]Provider
}

func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]Provider),
	}
}

func (r *Registry) Register(p Provider) {
	for _, host := range p.GetHosts() {
		r.providers[host] = p
	}
}

func (r *Registry) createMigration(typ gitea.GitServiceType, opts Options) gitea.MigrateRepoOption {
	return gitea.MigrateRepoOption{
		RepoOwner:    opts.OwnerName,
		RepoName:     opts.RepositoryName,
		CloneAddr:    opts.CloneURL,
		Service:      typ,
		AuthToken:    opts.Token,
		Mirror:       true,
		Description:  opts.Description,
		Wiki:         true,
		Milestones:   true,
		Labels:       true,
		Issues:       true,
		PullRequests: true,
		Releases:     true,
		LFS:          true,
	}
}

func (r *Registry) RetrieveRepoOptions(ctx context.Context, repo *url.URL, login *config.Login) (*gitea.MigrateRepoOption, error) {
	providerId := strings.ToLower(repo.Hostname())
	provider, ok := r.providers[providerId]
	if !ok {
		return nil, nil
	}

	opts, err := provider.RetrieveRepoOptions(ctx, repo, login)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve options via provider %s: %w", providerId, err)
	}

	migration := r.createMigration(provider.GetType(), *opts)
	return &migration, nil
}

func (r *Registry) RetrieveProfileOptions(ctx context.Context, profile *url.URL, login *config.Login) ([]gitea.MigrateRepoOption, error) {
	providerId := strings.ToLower(profile.Hostname())
	provider, ok := r.providers[providerId]
	if !ok {
		return nil, nil
	}

	opts, err := provider.RetrieveProfileOptions(ctx, profile, login)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve options via provider %s: %w", providerId, err)
	}

	migrations := make([]gitea.MigrateRepoOption, len(opts))
	for i, opt := range opts {
		migrations[i] = r.createMigration(provider.GetType(), opt)
	}
	return migrations, nil
}
