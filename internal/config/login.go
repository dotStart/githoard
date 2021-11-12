package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const LoginFileName = "login.json"

type Login struct {
	InstanceUri string `json:"instance_uri"`
	Token       string `json:"token"`
	GitHubToken string `json:"github-token"`
}

func ReadLogin() (*Login, error) {
	home, err := DirectoryPath()
	if err != nil {
		return nil, err
	}

	location := filepath.Join(home, LoginFileName)
	enc, err := os.ReadFile(location)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to read: %w", err)
	}

	var l Login
	if err := json.Unmarshal(enc, &l); err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}
	return &l, nil
}

func (l *Login) Write() error {
	home, err := DirectoryPath()
	if err != nil {
		return err
	}

	location := filepath.Join(home, LoginFileName)

	enc, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("failed to encode: %w", err)
	}

	if err := os.WriteFile(location, enc, 0600); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}
