package templ

import (
	"context"

	"github.com/open-feature/go-sdk/openfeature"
)

type TestProvider struct {
	Flags map[string]any
}

// BooleanEvaluation implements openfeature.FeatureProvider.
func (t *TestProvider) BooleanEvaluation(ctx context.Context, flag string, defaultValue bool, evalCtx openfeature.FlattenedContext) openfeature.BoolResolutionDetail {
	val, ok := t.Flags[flag]
	if !ok {
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewFlagNotFoundResolutionError(flag),
			},
		}
	}

	bVal, ok := val.(bool)
	if !ok {
		return openfeature.BoolResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewTypeMismatchResolutionError("not bool"),
			},
		}
	}

	return openfeature.BoolResolutionDetail{
		Value: bVal,
	}
}

// FloatEvaluation implements openfeature.FeatureProvider.
func (t *TestProvider) FloatEvaluation(ctx context.Context, flag string, defaultValue float64, evalCtx openfeature.FlattenedContext) openfeature.FloatResolutionDetail {
	val, ok := t.Flags[flag]
	if !ok {
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewFlagNotFoundResolutionError(flag),
			},
		}
	}

	fVal, ok := val.(float64)
	if !ok {
		return openfeature.FloatResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewTypeMismatchResolutionError("not float"),
			},
		}
	}

	return openfeature.FloatResolutionDetail{
		Value: fVal,
	}
}

// Hooks implements openfeature.FeatureProvider.
func (t *TestProvider) Hooks() []openfeature.Hook {
	return []openfeature.Hook{}
}

// IntEvaluation implements openfeature.FeatureProvider.
func (t *TestProvider) IntEvaluation(ctx context.Context, flag string, defaultValue int64, evalCtx openfeature.FlattenedContext) openfeature.IntResolutionDetail {
	val, ok := t.Flags[flag]
	if !ok {
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewFlagNotFoundResolutionError(flag),
			},
		}
	}

	iVal, ok := val.(int64)
	if !ok {
		return openfeature.IntResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewTypeMismatchResolutionError("not float"),
			},
		}
	}

	return openfeature.IntResolutionDetail{
		Value: iVal,
	}
}

// Metadata implements openfeature.FeatureProvider.
func (t *TestProvider) Metadata() openfeature.Metadata {
	return openfeature.Metadata{
		Name: "Test Provider",
	}
}

// ObjectEvaluation implements openfeature.FeatureProvider.
func (t *TestProvider) ObjectEvaluation(ctx context.Context, flag string, defaultValue interface{}, evalCtx openfeature.FlattenedContext) openfeature.InterfaceResolutionDetail {
	val, ok := t.Flags[flag]
	if !ok {
		return openfeature.InterfaceResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewFlagNotFoundResolutionError(flag),
			},
		}
	}

	return openfeature.InterfaceResolutionDetail{
		Value: val,
	}
}

// StringEvaluation implements openfeature.FeatureProvider.
func (t *TestProvider) StringEvaluation(ctx context.Context, flag string, defaultValue string, evalCtx openfeature.FlattenedContext) openfeature.StringResolutionDetail {
	val, ok := t.Flags[flag]
	if !ok {
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewFlagNotFoundResolutionError(flag),
			},
		}
	}

	sVal, ok := val.(string)
	if !ok {
		return openfeature.StringResolutionDetail{
			Value: defaultValue,
			ProviderResolutionDetail: openfeature.ProviderResolutionDetail{
				ResolutionError: openfeature.NewTypeMismatchResolutionError("not float"),
			},
		}
	}

	return openfeature.StringResolutionDetail{
		Value: sVal,
	}
}

var _ openfeature.FeatureProvider = &TestProvider{}
