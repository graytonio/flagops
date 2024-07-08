package gittools

import (
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/graytonio/flagops/lib/config"
)

// Caches git repos for reuse in multiple paths
var repos map[string]*git.Repository = map[string]*git.Repository{}

func Clone(repoConf config.GitRepo) (*git.Repository, error) {
	if repo, ok := repos[repoConf.Repo]; ok {
		return repo, nil
	}

	opts := &git.CloneOptions{
		URL:    repoConf.Repo,
		Mirror: true,
	}

	authMethod, err := buildAuthMethod(repoConf.Auth)
	if err != nil {
		return nil, err
	}

	if authMethod != nil {
		opts.Auth = authMethod
	}

	repo, err := git.Clone(memory.NewStorage(), memfs.New(), opts)
	if err != nil {
		return nil, err
	}

	repos[repoConf.Repo] = repo
	return repo, nil
}

func buildAuthMethod(auth config.Auth) (transport.AuthMethod, error) {
	if auth.Username != "" && auth.Password != "" {
		return buildBasicAuth(auth)
	}

	if auth.PrivateKey != "" {
		return buildPrivateKeyAuth(auth)
	}

	return nil, nil
}

// Builds auth using private keys. Checks if the supplied value is a file if so reads the key from the file. Otherwise attempts to read the supplied env variable for private key data.
func buildPrivateKeyAuth(auth config.Auth) (transport.AuthMethod, error) {
	privKey, ok := os.LookupEnv(auth.PrivateKey)
	if ok {
		return ssh.NewPublicKeys("git", []byte(privKey), "")
	}

	return ssh.NewPublicKeysFromFile("git", auth.PrivateKey, "")
}

// Builds a basic auth method using the passed username and reading the password from the supplied env variable
func buildBasicAuth(auth config.Auth) (transport.AuthMethod, error) {
	pass, ok := os.LookupEnv(auth.Password)
	if !ok {
		return nil, fmt.Errorf("password could not be read from env %s", auth.Password)
	}

	basicAuth := &http.BasicAuth{
		Username: auth.Username,
		Password: pass,
	}

	return basicAuth, nil
}
