package api

import (
	"context"
	goTypes "go/types"
	"net/http"
	"time"

	"github.com/borosr/realworld/lib/api"
	"github.com/borosr/realworld/lib/auth"
	"github.com/borosr/realworld/lib/middleware"
	"github.com/borosr/realworld/persist"
	"github.com/borosr/realworld/types"
	"golang.org/x/crypto/bcrypt"
)

type userController struct {
	userRepository persist.Repository[*types.User]
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
	user, err := uc.userRepository.Get(ctx, u.User.Email)
	if err != nil {
		return fallback, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.User.Password)); err != nil {
		return fallback, err
	}
	if user.Token == "" {
		token, err := auth.Sign(map[string]interface{}{
			"email": user.Email,
			"iat":   time.Now().UTC().Unix(),
		})
		if err != nil {
			return fallback, err
		}
		user.Token = token
		if _, err := uc.userRepository.Save(ctx, user); err != nil {
			return fallback, err
		}
	}

	fallback.User = *user
	return fallback, nil
}

func (uc userController) registration(ctx context.Context, u types.UserWrapper[types.UserSignUp]) (types.UserWrapper[types.User], error) {
	var fallback types.UserWrapper[types.User]
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return fallback, err
	}
	saved, err := uc.userRepository.Save(ctx, &types.User{
		Email: u.User.Email,
		Profile: types.Profile{
			Username: u.User.Username,
		},
		Password: string(encryptedPassword),
	})
	if err != nil {
		return fallback, err
	}
	fallback.User = *saved
	return fallback, nil
}

func (uc userController) currentUser(ctx context.Context, u goTypes.Nil) (types.UserWrapper[types.User], error) {
	var fallback types.UserWrapper[types.User]
	username, err := api.GetValue[string](ctx, "email")
	if err != nil {
		return fallback, err
	}
	user, err := uc.getAuthenticatedUser(ctx, username)
	if err != nil {
		return fallback, err
	}
	fallback.User = user
	return fallback, nil
}

func (uc userController) updateUser(ctx context.Context, u types.UserWrapper[types.User]) (types.UserWrapper[types.User], error) {
	var fallback types.UserWrapper[types.User]
	user, err := uc.userRepository.Get(ctx, u.User.Email)
	if err != nil {
		return fallback, err
	}
	if u.User.Email != "" {
		user.Email = u.User.Email
	}
	if u.User.Token != "" {
		user.Token = u.User.Token
	}
	if u.User.Username != "" {
		user.Username = u.User.Username
	}
	if u.User.Bio != "" {
		user.Bio = u.User.Bio
	}
	if u.User.Image != "" {
		user.Image = u.User.Image
	}
	if u.User.Password != "" {
		user.Password = u.User.Password
	}
	saved, err := uc.userRepository.Save(ctx, user)
	if err != nil {
		return fallback, err
	}
	fallback.User = *saved
	return fallback, nil
}

func (uc userController) getAuthenticatedUser(ctx context.Context, username string) (types.User, error) {
	user, err := uc.userRepository.Get(ctx, username)
	if err != nil {
		return types.User{}, err
	}
	return *user, nil
}
