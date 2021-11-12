package service

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"fmt"
	"github.com/dotstart/githoard/internal/config"
	"net/url"
	"strings"
)

import (
	"github.com/google/go-github/v40/github"
)

var gitHubHosts = []string{"github.com"}

type gitHubProvider struct {
	client *github.Client
}

func NewGitHubProvider() Provider {
	return &gitHubProvider{
		client: github.NewClient(nil),
	}
}

func (gh *gitHubProvider) GetType() gitea.GitServiceType {
	return gitea.GitServiceGithub
}

func (gh *gitHubProvider) GetHosts() []string {
	return gitHubHosts
}

func (gh *gitHubProvider) createRepoOptions(repo *github.Repository, login *config.Login) *Options {
	return &Options{
		OwnerName:      repo.GetOwner().GetLogin(),
		RepositoryName: repo.GetName(),
		CloneURL:       repo.GetCloneURL(),
		Token:          login.GitHubToken,
		Description:    repo.GetDescription(),
	}
}

func (gh *gitHubProvider) RetrieveRepoOptions(ctx context.Context, repo *url.URL, login *config.Login) (*Options, error) {
	p := repo.Path
	if p[0] == '/' {
		p = p[1:]
	}

	splitPath := strings.Split(p, "/")
	if len(splitPath) != 2 {
		return nil, fmt.Errorf("invalid repository location: %s", repo.Path)
	}

	originalOwner := splitPath[0]
	originalRepo := splitPath[1]

	meta, _, err := gh.client.Repositories.Get(ctx, originalOwner, originalRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to contact GitHub: %s", err)
	}

	return gh.createRepoOptions(meta, login), nil
}

func (gh *gitHubProvider) RetrieveProfileOptions(ctx context.Context, owner *url.URL, login *config.Login) ([]Options, error) {
	p := owner.Path
	if p[0] == '/' {
		p = p[1:]
	}

	repos := make([]*github.Repository, 0)

	i := 1
	for {
		opts := &github.RepositoryListOptions{
			ListOptions: github.ListOptions{
				Page:    i,
				PerPage: 100,
			},
		}

		page, _, err := gh.client.Repositories.List(ctx, p, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve repositories: %s", err)
		}

		if len(page) == 0 {
			break
		}

		repos = append(repos, page...)
		i++
	}

	opts := make([]Options, len(repos))
	for i, repo := range repos {
		opts[i] = *gh.createRepoOptions(repo, login)
	}
	return opts, nil
}
