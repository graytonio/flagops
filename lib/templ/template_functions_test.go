package templ

import (
	"testing"

	"github.com/open-feature/go-sdk/openfeature"
	"github.com/stretchr/testify/assert"
)

func TestParseObjectFlagSimpleIndex(t *testing.T) {
	provider := &TestProvider{
		Flags: map[string]any{
			"test0": map[string]any{
				"key": "foo",
			},
		},
	}

	openfeature.SetProvider(provider)
	engine := TemplateEngine{
		FlagProvider: openfeature.NewClient("test"),
	}

	res, err := engine.env("test0.key")
	assert.NoError(t, err)
	assert.Equal(t, "foo", res)
}

func TestParseObjectFlagNestedIndex(t *testing.T) {
	provider := &TestProvider{
		Flags: map[string]any{
			"test0": map[string]any{
				"nested": map[string]any{
					"obj": map[string]any{
						"key": "foo",
					},
				},
			},
		},
	}

	openfeature.SetProvider(provider)
	engine := TemplateEngine{
		FlagProvider: openfeature.NewClient("test"),
	}

	res, err := engine.env("test0.nested.obj.key")
	assert.NoError(t, err)
	assert.Equal(t, "foo", res)
}

func TestParseBoolFlag(t *testing.T) {
	provider := &TestProvider{
		Flags: map[string]any{
			"test0": true,
		},
	}

	openfeature.SetProvider(provider)
	engine := TemplateEngine{
		FlagProvider: openfeature.NewClient("test"),
	}

	res, err := engine.env("test0_enabled")
	assert.NoError(t, err)
	assert.Equal(t, true, res)
}

func TestParseString(t *testing.T) {
	provider := &TestProvider{
		Flags: map[string]any{
			"test0": "foo",
		},
	}

	openfeature.SetProvider(provider)
	engine := TemplateEngine{
		FlagProvider: openfeature.NewClient("test"),
	}

	res, err := engine.env("test0")
	assert.NoError(t, err)
	assert.Equal(t, "foo", res)

}

func TestParseStringInt(t *testing.T) {
	provider := &TestProvider{
		Flags: map[string]any{
			"test1": int64(12),
		},
	}

	openfeature.SetProvider(provider)
	engine := TemplateEngine{
		FlagProvider: openfeature.NewClient("test"),
	}

	res, err := engine.env("test1")
	assert.NoError(t, err)
	assert.Equal(t, "12", res)
}

func TestParseStringFloat(t *testing.T) {
	provider := &TestProvider{
		Flags: map[string]any{
			"test2": float64(12.2),
		},
	}

	openfeature.SetProvider(provider)
	engine := TemplateEngine{
		FlagProvider: openfeature.NewClient("test"),
	}

	res, err := engine.env("test2")
	assert.NoError(t, err)
	assert.Equal(t, "12.2", res)
}
