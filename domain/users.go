package domain

import (
	"context"
	"time"

	"github.com/borosr/realworld/lib/auth"
	"github.com/borosr/realworld/persist"
	"github.com/borosr/realworld/types"
	"golang.org/x/crypto/bcrypt"
)

type UserDescriptor interface {
	Login(ctx context.Context, u types.UserLogin) (types.User, error)
	SignUp(ctx context.Context, u types.UserSignUp) (types.User, error)
	GetByEmail(ctx context.Context, email string) (types.User, error)
	Update(ctx context.Context, u types.User) (types.User, error)
}

type UserService struct {
	UserRepository persist.Repository[*types.User]
}

func (us UserService) Login(ctx context.Context, u types.UserLogin) (types.User, error) {
	user, err := us.UserRepository.Get(ctx, u.Email)
	if err != nil {
		return types.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return types.User{}, err
	}
	if user.Token == "" {
		token, err := auth.Sign(map[string]interface{}{
			"email": user.Email,
			"iat":   time.Now().UTC().Unix(),
		})
		if err != nil {
			return types.User{}, err
		}
		user.Token = token
		if _, err := us.UserRepository.Save(ctx, user); err != nil {
			return types.User{}, err
		}
	}
	return *user, nil
}

func (us UserService) SignUp(ctx context.Context, u types.UserSignUp) (types.User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return types.User{}, err
	}
	saved, err := us.UserRepository.Save(ctx, &types.User{
		Email: u.Email,
		Profile: types.Profile{
			Username: u.Username,
		},
		Password: string(encryptedPassword),
	})
	if err != nil {
		return types.User{}, err
	}
	return *saved, nil
}

func (us UserService) GetByEmail(ctx context.Context, email string) (types.User, error) {
	user, err := us.UserRepository.Get(ctx, email)
	if err != nil {
		return types.User{}, err
	}
	return *user, nil
}

func (us UserService) Update(ctx context.Context, u types.User) (types.User, error) {
	user, err := us.UserRepository.Get(ctx, u.Email)
	if err != nil {
		return types.User{}, err
	}
	if u.Email != "" {
		user.Email = u.Email
	}
	if u.Token != "" {
		user.Token = u.Token
	}
	if u.Username != "" {
		user.Username = u.Username
	}
	if u.Bio != "" {
		user.Bio = u.Bio
	}
	if u.Image != "" {
		user.Image = u.Image
	}
	if u.Password != "" {
		user.Password = u.Password
	}
	saved, err := us.UserRepository.Save(ctx, user)
	if err != nil {
		return types.User{}, err
	}
	return *saved, nil
}
