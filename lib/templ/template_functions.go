package templ

import (
	"context"
	"fmt"
	"strings"

	"github.com/open-feature/go-sdk/openfeature"
	"gopkg.in/yaml.v3"
)


func (te *TemplateEngine) env(feature string) (any, error) {
	if strings.Contains(feature, ".") {
		parts := strings.Split(feature, ".")
		featureName := parts[0]
		data, err := te.FlagProvider.ObjectValue(context.Background(), featureName, map[string]any{}, openfeature.EvaluationContext{})
		if err != nil {
			return nil, err
		}

		fmt.Printf("%+v", data)

		var value = data
		for i := 1; i < len(parts); i++ {
			v, ok := value.(map[string]any)
			if !ok {
				break
			}

			value = v[parts[i]]
			continue
		}

		return value, nil
	}

	if strings.HasSuffix(feature, "_enabled") {
		data, err := te.FlagProvider.BooleanValue(context.Background(), strings.TrimSuffix(feature, "_enabled"), false, openfeature.EvaluationContext{})
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


// This has been copied from helm and may be removed as soon as it is retrofited in sprig
// toYAML takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func (te *TemplateEngine) toYAML(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return "", err
	}
	return strings.TrimSuffix(string(data), "\n"), nil
}

// This has been copied from helm and may be removed as soon as it is retrofited in sprig
// fromYAML converts a YAML document into a map[string]interface{}.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string into
// m["Error"] in the returned map.
func (te *TemplateEngine) fromYAML(str string) (map[string]interface{}, error) {
	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		return nil, err
	}
	return m, nil
}

// This has been copied from helm and may be removed as soon as it is retrofited in sprig
// fromYAMLArray converts a YAML array into a []interface{}.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string as
// the first and only item in the returned array.
func (te *TemplateEngine) fromYAMLArray(str string) ([]interface{}, error) {
	a := []interface{}{}

	if err := yaml.Unmarshal([]byte(str), &a); err != nil {
		return nil, err
	}
	return a, nil
}