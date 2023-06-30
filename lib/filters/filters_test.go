package filters_test

import (
	"testing"

	"github.com/graytonio/flagops/lib/filters"
	"github.com/sirupsen/logrus"
)

func TestBlankLineFilter_Parse(t *testing.T) {
	filter := filters.BlankLineFilter{}
	log := logrus.NewEntry(logrus.StandardLogger())

	testCases := []struct {
		name           string
		input          string
		expectedOutput string
		expectedError  error
	}{
		{
			name:           "Case 1",
			input:          "Hello\n\nWorld",
			expectedOutput: "Hello\nWorld",
			expectedError:  nil,
		},
		{
			name:           "Case 2",
			input:          "   \n\n   \n",
			expectedOutput: "",
			expectedError:  nil,
		},
		{
			name:           "Case 3",
			input:          "NoBlankLines",
			expectedOutput: "NoBlankLines",
			expectedError:  nil,
		},
		{
			name:           "Case 4",
			input:          "",
			expectedOutput: "",
			expectedError:  nil,
		},
		{
			name:           "Case 5",
			input:          "Single Line\n",
			expectedOutput: "Single Line\n",
			expectedError:  nil,
		},
		{
			name:           "Case 6",
			input:          "   \n\n   \nText",
			expectedOutput: "Text",
			expectedError:  nil,
		},
		{
			name:           "Case 7",
			input:          "   \n\n   \nText\n",
			expectedOutput: "Text\n",
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := filter.Parse(tc.input, log)

			if output != tc.expectedOutput {
				t.Errorf("Expected output: %q, but got: %q", tc.expectedOutput, output)
			}

			if err != tc.expectedError {
				t.Errorf("Expected error: %v, but got: %v", tc.expectedError, err)
			}
		})
	}
}