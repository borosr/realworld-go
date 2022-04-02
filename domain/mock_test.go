package domain

import (
	"context"

	persistTypes "github.com/borosr/realworld/persist/types"
	"github.com/borosr/realworld/types"
	"github.com/stretchr/testify/mock"
)

type MockRepository[Type persistTypes.Storable] struct {
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

func (m *MockRepository[Type]) GetFiltered(ctx context.Context, filters ...persistTypes.Filter[Type]) ([]Type, error) {
	var is = make([]interface{}, 0, len(filters)+1)
	is = append(is, ctx)
	for _, f := range filters {
		is = append(is, f)
	}
	args := m.Called(is...)
	return args.Get(0).([]Type), args.Error(1)
}

func (m *MockRepository[Type]) CountFiltered(ctx context.Context, filters ...persistTypes.Filter[Type]) (uint64, error) {
	var is = make([]interface{}, 0, len(filters)+1)
	is = append(is, ctx)
	for _, f := range filters {
		is = append(is, f)
	}
	args := m.Called(is...)
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

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Login(ctx context.Context, u types.UserLogin) (types.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(types.User), args.Error(1)
}

func (m *MockUserService) SignUp(ctx context.Context, u types.UserSignUp) (types.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(types.User), args.Error(1)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (types.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(types.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, u types.User) (types.User, error) {
	args := m.Called(ctx, u)
	return args.Get(0).(types.User), args.Error(1)
}
