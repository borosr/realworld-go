package api

import (
	"context"
	"strconv"

	"github.com/borosr/realworld/lib/broken"
)

const ctxPathVariablePrefix = "path_var_"

type ContextTypeConstraint interface {
	~string | ~int
}

func PathVariable[VariableType ContextTypeConstraint](ctx context.Context, key string) (VariableType, error) {
	return GetValue[VariableType](ctx, ctxPathVariablePrefix+key)
}
func SetPathVariable[VariableType ContextTypeConstraint](ctx context.Context, key string, value VariableType) context.Context {
	return context.WithValue(ctx, ctxPathVariablePrefix+key, value)
}

func GetValue[VariableType ContextTypeConstraint](ctx context.Context, key string) (VariableType, error) {
	var fallbackResult VariableType
	value := ctx.Value(key)
	if value == nil {
		return fallbackResult, broken.Internal("key not found")
	}

	switch (interface{})(fallbackResult).(type) {
	case string:
		return value.(VariableType), nil
	case int:
		valStr := value.(string)
		res, err := strconv.Atoi(valStr)
		if err != nil {
			return fallbackResult, err
		}
		return (interface{})(res).(VariableType), nil
	}

	return fallbackResult, broken.Internal("unable to parse context value")
}
