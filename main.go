package main

import (
	"github.com/graytonio/flagops/lib/config"
	"github.com/graytonio/flagops/lib/provider"
	"github.com/graytonio/flagops/lib/templ"
	"github.com/sirupsen/logrus"
)

func main() {
    logrus.SetLevel(logrus.DebugLevel)
    conf, err := config.LoadConfig("")
    if err != nil {
        panic(err)
    }

    providers, err := provider.ConfigureProviders(conf.Envs)
    if err != nil {
        panic(err)
    }

    engines, err := templ.CreateEngines(conf.Paths, providers)
    if err != nil {
        panic(err)
    }

    for _, engine := range engines {
        err = engine.Execute()
        if err != nil {
            panic(err)
        }
    }
}
