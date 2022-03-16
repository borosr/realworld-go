package domain

import (
	"context"
	"fmt"

	"github.com/borosr/realworld/persist"
	"github.com/borosr/realworld/types"
)

type ProfileDescriptor interface {
	GetByUsername(ctx context.Context, username string) (types.Profile, error)
	Follow(ctx context.Context, followBy, following string) (types.Profile, error)
	Unfollow(ctx context.Context, followBy, following string) (types.Profile, error)
}

type ProfileService struct {
	UserRepository   persist.Repository[*types.User]
	FollowRepository persist.Repository[*types.Follow]
}

func (ps ProfileService) GetByUsername(ctx context.Context, username string) (types.Profile, error) {
	users, err := ps.UserRepository.GetFiltered(ctx, func(u *types.User) bool {
		return u.Username == username
	})
	if err != nil {
		return types.Profile{}, err
	}
	if len(users) != 1 {
		return types.Profile{}, fmt.Errorf("unable to find profile with username: %s", username)
	}
	return users[0].Profile, nil
}

func (ps ProfileService) Follow(ctx context.Context, from, to string) (types.Profile, error) {
	if ps.hasFollowedBy(ctx, from, to) {
		return types.Profile{}, fmt.Errorf("profile %s already followed by: %s", to, from)
	}

	following, err := ps.GetByUsername(ctx, to)
	if err != nil {
		return types.Profile{}, err
	}

	f := types.Follow{
		From: from,
		To:   to,
	}
	if _, err := ps.FollowRepository.Save(ctx, &f); err != nil {
		return types.Profile{}, err
	}

	following.Following = true
	return following, nil
}

func (ps ProfileService) Unfollow(ctx context.Context, from, to string) (types.Profile, error) {
	if !ps.hasFollowedBy(ctx, from, to) {
		return types.Profile{}, fmt.Errorf("profile %s haven't followed by: %s", to, from)
	}

	following, err := ps.GetByUsername(ctx, to)
	if err != nil {
		return types.Profile{}, err
	}

	f := types.Follow{
		From: from,
		To:   to,
	}
	if err := ps.FollowRepository.Delete(ctx, f.Key()); err != nil {
		return types.Profile{}, err
	}

	following.Following = false
	return following, nil
}

func (ps ProfileService) hasFollowedBy(ctx context.Context, from string, to string) bool {
	f := types.Follow{
		From: from,
		To:   to,
	}
	_, err := ps.FollowRepository.Get(ctx, f.Key())
	return err == nil
}
