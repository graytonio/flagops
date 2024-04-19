package destinations

import (
	"os"
	"path"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
	"github.com/graytonio/flagops/lib/config"
)

type Output interface {
	// Initialize configured file system preparing it to accept files to write
	Init() error

	// Takes output content and writes it to destination at path
	ExecuteFile(path string, content []byte) error

	// Finalizes filesystem and cleans up any necessary resources
	Finalize() error
}

func NewOutput(conf config.Path) (Output, error) {
	switch conf.Destination.Type {
	case config.Git:
		return newGitOutput(conf)
	case config.File:
		return newFileOutput(conf)
	case config.Console:
		return newConsoleOutput()
	}

	return nil, nil
}

func cleanFSDestination(fs billy.Filesystem, path string) error {
	stat, err := fs.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if !stat.IsDir() {
		return os.Remove(path)
	}

	names, err := fs.ReadDir(path)
	if err != nil {
	  return err
	}

	for _, f := range names {
		err = util.RemoveAll(fs, filepath.Join(path, f.Name()))
		if err != nil {
		  return err
		}
	}

	return nil
}

func getFileOutputDestination(srcRoot string, destRoot string, file string) (string, error) {
	relPath, err := filepath.Rel(srcRoot, file)
	if err != nil {
		return "", err
	}

	return path.Join(destRoot, relPath), nil
}
