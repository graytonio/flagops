package destinations

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/graytonio/flagops/lib/config"
)

var _ Output = &FileOutput{}

type FileOutput struct {
	conf config.Path
	fs   billy.Filesystem
}

func newFileOutput(conf config.Path) (*FileOutput, error) {
	return &FileOutput{
		conf: conf,
	}, nil
}

// Init implements Output.
func (fo *FileOutput) Init() error {
	fo.fs = osfs.New(".")
	return nil
}

// ExecuteFile implements Output.
func (fo *FileOutput) ExecuteFile(path string, content []byte) error {
	destPath, err := getFileOutputDestination(fo.conf.Path, fo.conf.Destination.Path, path)
	if err != nil {
	  return err
	}

	f, err := fo.fs.Create(destPath)
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
func (fo *FileOutput) Finalize() error {
	return nil
}