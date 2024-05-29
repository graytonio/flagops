package templ

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"slices"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/graytonio/flagops/lib/config"
	"github.com/graytonio/flagops/lib/destinations"
	"github.com/open-feature/go-sdk/openfeature"
	"github.com/sirupsen/logrus"
)

type TemplateEngine struct {
	SourceFilesystem billy.Filesystem
	Output           destinations.Output
	FlagProvider     *openfeature.Client

	files   []string
	path    config.Path
	funcMap map[string]any
}

// TODO Move to config
var ignoreFileExtentions = []string{
	".tar",
	".tar.gz",
	".tgz",
	".zip",
	".gz",
	".7z",
}

// Parses configured paths and available providers to create list of engine tasks to be executed
func CreateEngines(paths []config.Path, providers map[string]*openfeature.Client) ([]*TemplateEngine, error) {
	engines := []*TemplateEngine{}
	for _, path := range paths {
		engine, err := NewTemplateEngine(path, providers[path.Env])
		if err != nil {
			return nil, err
		}

		engines = append(engines, engine)
	}
	return engines, nil
}

// Creates a new template engine with a source path, a destination path, and a feature flag provider
func NewTemplateEngine(path config.Path, provider *openfeature.Client) (*TemplateEngine, error) {
	engine := &TemplateEngine{
		SourceFilesystem: osfs.New("."),
		FlagProvider:     provider,

		path: path,
	}

	err := engine.ScanFiles()
	if err != nil {
		return nil, err
	}

	output, err := destinations.NewOutput(path)
	if err != nil {
		return nil, err
	}
	engine.Output = output

	funcMap := sprig.GenericFuncMap()
	funcMap["env"] = engine.env
	funcMap["toYaml"] = engine.toYAML
	funcMap["fromYaml"] = engine.fromYAML
	funcMap["fromYamlArray"] = engine.fromYAMLArray
	engine.funcMap = funcMap

	return engine, nil
}

// Can input path and create list of files to template
func (te *TemplateEngine) ScanFiles() error {
	return util.Walk(te.SourceFilesystem, te.path.Path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		logrus.WithField("file", path).Debug("Adding file to template engine")
		te.files = append(te.files, path)
		return nil
	})
}

// Iterates each file in the engine and writes it to the destination
func (te *TemplateEngine) Execute() error {
	err := te.Output.Init()
	if err != nil {
		return err
	}

	for _, file := range te.files {
		err := te.executeFileTemplate(file)
		if err != nil {
			return err
		}
	}

	return te.Output.Finalize()
}

func (te *TemplateEngine) executeFileTemplate(path string) error {
	file, err := te.SourceFilesystem.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if slices.Contains(ignoreFileExtentions, filepath.Ext(path)) {
		return te.Output.ExecuteFile(path, data)
	}

	if te.path.Destination.Header != "" {
		logrus.Info("Adding header")
		data = append([]byte(fmt.Sprintf("%s\n", te.path.Destination.Header)), data...)
	}

	if te.path.Destination.Footer != "" {
		logrus.Info("Adding footer")
		data = append(data, []byte(fmt.Sprintf("\n%s", te.path.Destination.Footer))...)
	}

	templ, err := template.New(path).
		Delims("[{", "}]").
		Funcs(te.funcMap).
		Parse(string(data))
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte{})
	err = templ.Execute(buf, nil)
	if err != nil {
		return err
	}

	return te.Output.ExecuteFile(path, buf.Bytes())
}
