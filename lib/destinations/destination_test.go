package destinations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileOutputDestination(t *testing.T) {
	var tests = []struct {
		name               string
		destPath           string
		srcPath            string
		file               string
		expectedOutputPath string
	}{
		{
			name:               "basic",
			destPath:           "apps/production",
			srcPath:            "apps/",
			file:               "apps/metallb/Chart.yaml",
			expectedOutputPath: "apps/production/metallb/Chart.yaml",
		},
		{
			name:               "basic2",
			destPath:           "build",
			srcPath:            "test/",
			file:               "test/test.txt",
			expectedOutputPath: "build/test.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			destFile, err := getFileOutputDestination(tt.srcPath, tt.destPath, tt.file)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOutputPath, destFile)
		})
	}

}
