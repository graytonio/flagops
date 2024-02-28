package templ

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/graytonio/flagops/lib/config"
	"github.com/graytonio/flagops/lib/gittools"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/sirupsen/logrus"
)

type TemplateEngine struct {
	SourceFilesystem      billy.Filesystem
	DestinationFilesystem billy.Filesystem
	FlagProvider          *openfeature.Client
	RootPath              string

	files []string
	dest  config.Destination
	funcMap map[string]any

	gitRepo *git.Repository
}

func NewTemplateEngine(rootPath string, dest config.Destination, provider *openfeature.Client) (*TemplateEngine, error) {
	engine := &TemplateEngine{
		SourceFilesystem: osfs.New("."),
		RootPath:         rootPath,
		FlagProvider:     provider,

		dest: dest,
	}

	err := engine.setupFilesystem()
	if err != nil {
		return nil, err
	}

	err = engine.ScanFiles()
	if err != nil {
		return nil, err
	}

	funcMap := sprig.GenericFuncMap()
	funcMap["env"] = engine.env
	funcMap["toYaml"] = engine.toYAML
	funcMap["fromYaml"] = engine.fromYAML
	engine.funcMap = funcMap

	return engine, nil
}

// Init configured destination fs
func (te *TemplateEngine) setupFilesystem() error {
	switch te.dest.Type {
	case "file":
		te.DestinationFilesystem = osfs.New(".")
	case "git":
		repo, err := gittools.Clone(te.dest.Repo)
		if err != nil {
			return err
		}

		w, err := repo.Worktree()
		if err != nil {
			return err
		}

		te.DestinationFilesystem = w.Filesystem
		te.gitRepo = repo
	default:
		return fmt.Errorf("unsupported destination: %s", te.dest.Type)
	}
	return nil
}

// Clean up any steps for filesystems like git commit
func (te *TemplateEngine) finalizeFilesystem() error {
	switch te.dest.Type {
	case "file":
		return nil
	case "git":
		return te.finalizeGitRepo()
	default:
		return fmt.Errorf("unsupported destination: %s", te.dest.Type)
	}
}

func (te *TemplateEngine) finalizeGitRepo() error {
	w, err := te.gitRepo.Worktree()
	if err != nil {
		return err
	}

	err = w.AddGlob("*")
	if err != nil {
		return err
	}

	_, err = w.Commit("flagops: Built templates", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "FlagOps",
			Email: "flagops@gmail.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	err = te.gitRepo.Push(&git.PushOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Can input path and create list of files to template
func (te *TemplateEngine) ScanFiles() error {
	return util.Walk(te.SourceFilesystem, te.RootPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		logrus.WithField("file", path).Debug("Adding file to template engine")
		te.files = append(te.files, path)
		return nil
	})
}

// Calculate desired output file location
func (te *TemplateEngine) getFileOutputDestination(originalPath string) (string, error) {
	relPath, err := filepath.Rel(te.RootPath, originalPath)
	if err != nil {
		return "", err
	}

	return path.Join(te.dest.Path, relPath), nil
}

// Iterates each file in the engine and writes it to the destination
func (te *TemplateEngine) Execute() error {
	for _, file := range te.files {
		err := te.executeFileTemplate(file)
		if err != nil {
			return err
		}
	}
	return te.finalizeFilesystem()
}

func (te *TemplateEngine) executeFileTemplate(path string) error {
	file, err := te.SourceFilesystem.Open(path)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	templ, err := template.New(path).
		Delims("[{", "}]").
		Funcs(te.funcMap).Parse(string(data))
	if err != nil {
		return err
	}

	destPath, err := te.getFileOutputDestination(path)
	if err != nil {
		return err
	}

	logrus.WithField("input_file", path).WithField("output_file", destPath).Debug("Writing file output")
	destFile, err := te.DestinationFilesystem.Create(destPath)
	if err != nil {
		return err
	}

	err = templ.Execute(destFile, nil)
	if err != nil {
		return err
	}

	return nil
}