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

type profilesController struct {
	profileService domain.ProfileDescriptor
	userService    domain.UserDescriptor
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
	profile, err := pc.profileService.GetByUsername(ctx, username)
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	return types.ProfileWrapper{
		Profile: profile,
	}, nil
}

func (pc profilesController) follow(ctx context.Context, _ goTypes.Nil) (types.ProfileWrapper, error) {
	username, err := api.PathVariable[string](ctx, "username")
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	email, err := api.GetValue[string](ctx, "email")
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	currentUser, err := pc.userService.GetByEmail(ctx, email)
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	profile, err := pc.profileService.Follow(ctx, currentUser.Username, username)
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	return types.ProfileWrapper{
		Profile: profile,
	}, nil
}

func (pc profilesController) unfollow(ctx context.Context, _ goTypes.Nil) (types.ProfileWrapper, error) {
	username, err := api.PathVariable[string](ctx, "username")
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	email, err := api.GetValue[string](ctx, "email")
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	currentUser, err := pc.userService.GetByEmail(ctx, email)
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	profile, err := pc.profileService.Unfollow(ctx, currentUser.Username, username)
	if err != nil {
		return types.ProfileWrapper{}, err
	}
	return types.ProfileWrapper{
		Profile: profile,
	}, nil
}
