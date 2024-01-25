package templ

import (
	"github.com/graytonio/flagops/lib/config"
	"github.com/open-feature/go-sdk/openfeature"
)

func CreateEngines(paths []config.Path, providers map[string]*openfeature.Client) ([]*TemplateEngine, error) {
	engines := []*TemplateEngine{}
	for _, path := range paths {
		engine, err := NewTemplateEngine(path.Path, path.Destination, providers[path.Env])
		if err != nil {
			return nil, err
		}
		
		engines = append(engines, engine)
	}
	return engines, nil
}