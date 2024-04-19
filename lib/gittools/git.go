package gittools

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

// Caches git repos for reuse in multiple paths
var repos map[string]*git.Repository = map[string]*git.Repository{}

func Clone(url string) (*git.Repository, error) {
	if repo, ok := repos[url]; ok {
		return repo, nil
	}

	repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return nil, err
	}

	repos[url] = repo
	return repo, nil
}
