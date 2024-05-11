package templ

import (
	"context"
	"fmt"
	"strings"

	"github.com/open-feature/go-sdk/openfeature"
	"gopkg.in/yaml.v3"
)

func (te *TemplateEngine) getEvaluationContext() openfeature.EvaluationContext {
	return openfeature.NewEvaluationContext(te.path.Identity, te.path.Properties)
}

func (te *TemplateEngine) env(feature string) (any, error) {
	switch {
	case strings.Contains(feature, "."): // env "my_flag.subKey"
		parts := strings.Split(feature, ".")
		return te.parseObjectFlag(parts[0], parts[1:])
	case strings.HasSuffix(feature, "_enabled"): // env "my_flag_enabled"
		return te.parseBooleanFlag(strings.TrimSuffix(feature, "_enabled"))
	default: // env "my_flag"
		return te.parseStringFlag(feature)
	}
}

// Fetch feature flag flagKey and interpret the value as a json object. Recursively index the object until the desired key is found
func (te *TemplateEngine) parseObjectFlag(flagKey string, subKeys []string) (any, error) {
	data, err := te.FlagProvider.ObjectValue(context.Background(), flagKey, map[string]any{}, te.getEvaluationContext())
	if err != nil {
		return nil, err
	}

	var value = data
	for i := 0; i < len(subKeys); i++ {
		v, ok := value.(map[string]any)
		if !ok {
			value = "nil"
			break
		}

		value, ok = v[subKeys[i]]
		if !ok {
			value = "nil"
			break
		}
		continue
	}

	return value, nil
}

func (te *TemplateEngine) parseBooleanFlag(flagKey string) (bool, error) {
	return te.FlagProvider.BooleanValue(context.Background(), flagKey, false, te.getEvaluationContext())
}

func (te *TemplateEngine) parseStringFlag(flagKey string) (string, error) {
	iData, err := te.FlagProvider.IntValue(context.Background(), flagKey, -1, te.getEvaluationContext())
	if err == nil { // If flag can be interpreted as Int
		return fmt.Sprint(iData), nil
	}

	fData, err := te.FlagProvider.FloatValue(context.Background(), flagKey, -1, te.getEvaluationContext())
	if err == nil { // If flag can be interpreted as Float
		return fmt.Sprint(fData), nil
	}

	data, err := te.FlagProvider.StringValue(context.Background(), flagKey, "", te.getEvaluationContext())
	if err != nil {
		// Edge case where evaluating a disabled flag in some providers results in a null value
		if strings.Contains(err.Error(), "TYPE_MISMATCH") {
			return "", nil
		}

		return "", err
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
