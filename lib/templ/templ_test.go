package templ

import (
	"testing"

	"github.com/graytonio/flagops/lib/config"
	"github.com/stretchr/testify/assert"
)

func TestGetFileOutputDestination(t *testing.T) {
	var tests = []struct {
		name               string
		destPath           string
		rootPath           string
		match              string
		expectedOutputPath string
	}{
		{
			name:               "basic",
			destPath:           "apps/production",
			rootPath:           "apps/",
			match:              "apps/metallb/Chart.yaml",
			expectedOutputPath: "apps/production/metallb/Chart.yaml",
		},
		{
			name: "basic2",
			destPath: "build",
			rootPath: "test/",
			match: "test/test.txt",
			expectedOutputPath: "build/test.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := TemplateEngine{
				RootPath: tt.rootPath,
				dest: config.Destination{
					Path: tt.destPath,
				},
			}

			destination, err := engine.getFileOutputDestination(tt.match)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOutputPath, destination)
		})
	}

}
