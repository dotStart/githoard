package service

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"github.com/dotstart/githoard/internal/config"
	"net/url"
)

type Provider interface {
	GetType() gitea.GitServiceType
	GetHosts() []string

	RetrieveRepoOptions(ctx context.Context, repo *url.URL, login *config.Login) (*Options, error)
	RetrieveProfileOptions(ctx context.Context, owner *url.URL, login *config.Login) ([]Options, error)
}

type Options struct {
	OwnerName      string
	RepositoryName string
	CloneURL       string
	Description    string
	Token          string
}
