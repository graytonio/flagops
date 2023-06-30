package filters

import (
	"regexp"

	"github.com/sirupsen/logrus"
)

var FilterTable map[string]Filter

func init() {
	FilterTable = map[string]Filter{
		"blank": BlankLineFilter{},
	}
}

func FilterString(input string, filters []string, log *logrus.Entry) (output string, err error) {
	output = input

	for _, filter_id := range filters {
		output, err = FilterTable[filter_id].Parse(output, log)
		if err != nil {
			return "", err
		}
	}

	return output, nil
}

type Filter interface {
	Parse(string, *logrus.Entry) (string, error)
}

type BlankLineFilter struct {}

func (f BlankLineFilter) Parse(input string, log *logrus.Entry) (string, error) {
	var blankLineFilter = regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	log.Debug("Applying Blank Line Filter")
	return blankLineFilter.ReplaceAllString(input, ""), nil
}