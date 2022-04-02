package domain

import (
	"context"

	"github.com/borosr/realworld/persist/types"
	"github.com/stretchr/testify/mock"
)

type MockRepository[Type types.Storable] struct {
	mock.Mock
}

func (m *MockRepository[Type]) Save(ctx context.Context, data Type) (Type, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(Type), args.Error(1)
}

func (m *MockRepository[Type]) Get(ctx context.Context, key string) (Type, error) {
	args := m.Called(ctx, key)
	res := args.Get(0)
	err := args.Error(1)
	if res == nil {
		var t Type
		return t, err
	}
	return res.(Type), err
}

func (m *MockRepository[Type]) GetFiltered(ctx context.Context, filters ...types.Filter[Type]) ([]Type, error) {
	var is = make([]interface{}, 0, len(filters)+1)
	is = append(is, ctx)
	for _, f := range filters {
		is = append(is, f)
	}
	args := m.Called(is...)
	return args.Get(0).([]Type), args.Error(1)
}

func (m *MockRepository[Type]) CountFiltered(ctx context.Context, filters ...types.Filter[Type]) (uint64, error) {
	args := m.Called(ctx, filters)
	return uint64(args.Int(0)), args.Error(1)
}

func (m *MockRepository[Type]) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockRepository[Type]) Sequence(ctx context.Context, key string) (uint64, error) {
	args := m.Called(ctx, key)
	return uint64(args.Int(0)), args.Error(1)
}
