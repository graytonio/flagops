package templ

import (
	"context"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/graytonio/flagops/lib/config"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/sirupsen/logrus"
)

type TemplateEngine struct {
	Filesystem billy.Filesystem
	FlagProvider *openfeature.Client
	RootPath string
	
	files []string
	dest config.Destination
}

func NewTemplateEngine(rootPath string, dest config.Destination, provider *openfeature.Client) (*TemplateEngine, error) {
	engine := &TemplateEngine{
		Filesystem: osfs.New("."),
		RootPath: rootPath,
		FlagProvider: provider,

		dest: dest,
	}

	err := engine.ScanFiles()
	if err != nil {
		return nil, err
	}

	return engine, nil
}

func (te *TemplateEngine) ScanFiles() error {
	return util.Walk(te.Filesystem, te.RootPath, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		logrus.WithField("file", path).Debug("Adding file to template engine")
		te.files = append(te.files, path)
		return nil
	})
}

func (te *TemplateEngine) getFileOutputDestination(originalPath string) (string, error) {
	relPath, err := filepath.Rel(te.RootPath, originalPath)
	if err != nil {
		return "", err
	}
	
	return path.Join(te.dest.Path, relPath), nil
}

// Iterates each file in the engine and writes it to the destination
func (te *TemplateEngine) Execute() (error) {
	for _, file := range te.files {
		err := te.executeFileTemplate(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (te *TemplateEngine) executeFileTemplate(path string) (error) {
	file, err := te.Filesystem.Open(path)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	templ, err := template.New(path).
	Delims("[{", "}]").
	Funcs(template.FuncMap{
		"env": func (feature string) any {
			data, err := te.lookupFlagValue(feature)
			if err != nil {
				return nil
			}
			return data
		},
	}).Parse(string(data))
	if err != nil {
		return err
	}

	destPath, err := te.getFileOutputDestination(path)
	if err != nil {
		return err
	}

	logrus.WithField("input_file", path).WithField("output_file", destPath).Debug("Writing file output")
	destFile, err := te.Filesystem.Create(destPath)
	if err != nil {
		return err
	}

	err = templ.Execute(destFile, nil)
	if err != nil {
		return err
	}

	return nil
}

func (te *TemplateEngine) lookupFlagValue(feature string) (any, error) {
	if strings.Contains(feature, ".") {
		data, err := te.FlagProvider.ObjectValue(context.Background(), feature, nil, openfeature.EvaluationContext{})
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	data, err := te.FlagProvider.StringValue(context.Background(), feature, "", openfeature.EvaluationContext{})
	if err != nil {
		return nil, err
	}

	return data, nil
}