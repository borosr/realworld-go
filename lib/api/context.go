package api

import (
	"context"
	"errors"
	"strconv"
)

const ctxPathVariablePrefix = "path_var_"

type ContextTypeConstraint interface {
	~string | ~int
}

func PathVariable[VariableType ContextTypeConstraint](ctx context.Context, key string) (VariableType, error) {
	return GetValue[VariableType](ctx, ctxPathVariablePrefix+key)
}

func GetValue[VariableType ContextTypeConstraint](ctx context.Context, key string) (VariableType, error) {
	var fallbackResult VariableType
	value := ctx.Value(key)
	if value == nil {
		return fallbackResult, errors.New("key not found")
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
		return VariableType(res), nil
	}

	return fallbackResult, errors.New("unable to parse context value")
}
