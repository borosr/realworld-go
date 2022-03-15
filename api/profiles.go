package api

import (
	"context"
	"errors"
	goTypes "go/types"
	"net/http"

	"github.com/borosr/realworld/lib/api"
	"github.com/borosr/realworld/lib/middleware"
	"github.com/borosr/realworld/persist"
	"github.com/borosr/realworld/types"
)

type profilesController struct {
	userRepository persist.Repository[*types.User]
}

func (pc profilesController) Init() {
	api.Register[
		goTypes.Nil,
		types.ProfileWrapper,
		api.ControllerSimpleFunc[goTypes.Nil, types.ProfileWrapper],
	]("/api/profiles/{username}", http.MethodGet, pc.get)
	api.Register[
		goTypes.Nil,
		types.ProfileWrapper,
		api.ControllerSimpleFunc[goTypes.Nil, types.ProfileWrapper],
	]("/api/profiles/{username}/follow", http.MethodPost, pc.follow).
		PreProcess(middleware.TokenAuthentication)
	api.Register[
		goTypes.Nil,
		types.ProfileWrapper,
		api.ControllerSimpleFunc[goTypes.Nil, types.ProfileWrapper],
	]("/api/profiles/{username}/follow", http.MethodDelete, pc.unfollow).
		PreProcess(middleware.TokenAuthentication)
}

func (pc profilesController) get(ctx context.Context, _ goTypes.Nil) (types.ProfileWrapper, error) {
	username, err := api.PathVariable[string](ctx, "username")
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	users, err := pc.userRepository.GetFiltered(ctx, func(u *types.User) bool {
		return u.Username == username
	})
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	if len(users) != 1 {
		return types.ProfileWrapper{}, errors.New("")
	}
	return types.ProfileWrapper{
		Profile: users[0].Profile,
	}, nil
}

func (pc profilesController) follow(ctx context.Context, _ goTypes.Nil) (types.ProfileWrapper, error) {
	return pc.setFollow(ctx, true)
}

func (pc profilesController) unfollow(ctx context.Context, _ goTypes.Nil) (types.ProfileWrapper, error) {
	return pc.setFollow(ctx, false)
}

func (pc profilesController) setFollow(ctx context.Context, f bool) (types.ProfileWrapper, error) {
	username, err := api.PathVariable[string](ctx, "username")
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	users, err := pc.userRepository.GetFiltered(ctx, func(u *types.User) bool {
		return u.Username == username
	})
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	if len(users) != 1 {
		return types.ProfileWrapper{}, errors.New("")
	}
	user := users[0]
	user.Following = f
	saved, err := pc.userRepository.Save(ctx, user)
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	return types.ProfileWrapper{
		Profile: saved.Profile,
	}, nil
}
