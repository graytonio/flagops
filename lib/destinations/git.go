package destinations

import (
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/graytonio/flagops/lib/config"
	"github.com/graytonio/flagops/lib/gittools"
	"github.com/sirupsen/logrus"
)

var _ Output = &GitOutput{}

type GitOutput struct {
	conf config.Path
	fs   billy.Filesystem
	repo *git.Repository
}

func newGitOutput(conf config.Path) (*GitOutput, error) {
	return &GitOutput{
		conf: conf,
	}, nil
}

// Init implements Output.
func (g *GitOutput) Init() error {
	repo, err := gittools.Clone(g.conf.Destination.Repo)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Checkout branch if set
	if g.conf.Destination.Branch != "" {
		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(g.conf.Destination.Branch),
			Force:  true,
		})
		if err != nil {
			return err
		}
	}

	g.fs = w.Filesystem
	g.repo = repo

	if g.conf.Destination.UpsertMode {
		err = cleanFSDestination(g.fs, g.conf.Path)
		if err != nil {
			return err
		}
	}

	return nil
}

// ExecuteFile implements Output.
func (g *GitOutput) ExecuteFile(path string, content []byte) error {
	destPath, err := getFileOutputDestination(g.conf.Path, g.conf.Destination.Path, path)
	if err != nil {
		return err
	}

	f, err := g.fs.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}

// Finalize implements Output.
func (g *GitOutput) Finalize() error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	stat, err := w.Status()
	if err != nil {
		return err
	}

	if stat.IsClean() {
		logrus.Debug("nothing to commit")
		return nil
	}

	err = w.AddWithOptions(&git.AddOptions{
		All: true,
	})
	if err != nil {
		return err
	}

	_, err = w.Commit("flagops: Built templates", &git.CommitOptions{
		AllowEmptyCommits: false,
		Author: &object.Signature{
			Name:  "FlagOps",
			Email: "flagops@gmail.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	err = g.repo.Push(&git.PushOptions{})
	if err != nil {
		return err
	}

	return nil
}
