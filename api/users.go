package api

import (
	"context"
	goTypes "go/types"
	"net/http"

	"github.com/borosr/realworld/domain"
	"github.com/borosr/realworld/lib/api"
	"github.com/borosr/realworld/lib/middleware"
	"github.com/borosr/realworld/types"
)

type userController struct {
	userService domain.UserDescriptor
}

func (uc userController) Init() {
	api.Register[
		types.UserWrapper[types.UserLogin],
		types.UserWrapper[types.User],
		api.ControllerSimpleFunc[types.UserWrapper[types.UserLogin], types.UserWrapper[types.User]],
	]("/api/users/login", http.MethodPost, uc.login).
		Validated()
	api.Register[
		types.UserWrapper[types.UserSignUp],
		types.UserWrapper[types.User],
		api.ControllerSimpleFunc[types.UserWrapper[types.UserSignUp], types.UserWrapper[types.User]],
	]("/api/users", http.MethodPost, uc.registration).
		Validated()
	api.Register[
		goTypes.Nil,
		types.UserWrapper[types.User],
		api.ControllerSimpleFunc[goTypes.Nil, types.UserWrapper[types.User]],
	]("/api/user", http.MethodGet, uc.currentUser).
		PreProcess(middleware.TokenAuthentication).
		Validated()
	api.Register[
		types.UserWrapper[types.User],
		types.UserWrapper[types.User],
		api.ControllerSimpleFunc[types.UserWrapper[types.User], types.UserWrapper[types.User]],
	]("/api/user", http.MethodPut, uc.updateUser).
		PreProcess(middleware.TokenAuthentication)
}

func (uc userController) login(ctx context.Context, u types.UserWrapper[types.UserLogin]) (types.UserWrapper[types.User], error) {
	var fallback types.UserWrapper[types.User]
	user, err := uc.userService.Login(ctx, u.User)
	if err != nil {
		return fallback, err
	}
	fallback.User = user
	return fallback, nil
}

func (uc userController) registration(ctx context.Context, u types.UserWrapper[types.UserSignUp]) (types.UserWrapper[types.User], error) {
	var fallback types.UserWrapper[types.User]
	user, err := uc.userService.SignUp(ctx, u.User)
	if err != nil {
		return fallback, err
	}
	fallback.User = user
	return fallback, nil
}

func (uc userController) currentUser(ctx context.Context, _ goTypes.Nil) (types.UserWrapper[types.User], error) {
	var fallback types.UserWrapper[types.User]
	email, err := api.GetValue[string](ctx, "email")
	if err != nil {
		return fallback, err
	}
	user, err := uc.userService.GetByEmail(ctx, email)
	if err != nil {
		return fallback, err
	}
	fallback.User = user
	return fallback, nil
}

func (uc userController) updateUser(ctx context.Context, u types.UserWrapper[types.User]) (types.UserWrapper[types.User], error) {
	var fallback types.UserWrapper[types.User]
	user, err := uc.userService.Update(ctx, u.User)
	if err != nil {
		return fallback, err
	}
	fallback.User = user
	return fallback, nil
}
