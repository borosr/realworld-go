package types

import (
	"context"
	"errors"

	"github.com/borosr/realworld/lib/api"
)

type UserWrapper[Data UserLogin | User | UserSignUp] struct {
	User Data `json:"user"`
}

func (u UserWrapper[Data]) Validate(ctx context.Context) error {
	if v, ok := (interface{})(u.User).(api.Validator); ok {
		return v.Validate(ctx)
	}
	return nil
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u UserLogin) Validate(_ context.Context) error {
	if u.Email == "" {
		return errors.New("missing email field")
	}
	if u.Password == "" {
		return errors.New("missing password field")
	}
	return nil
}

type UserSignUp struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u UserSignUp) Validate(_ context.Context) error {
	if u.Username == "" {
		return errors.New("missing username field")
	}
	if u.Email == "" {
		return errors.New("missing email field")
	}
	if u.Password == "" {
		return errors.New("missing password field")
	}
	return nil
}

type User struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"password"`
	Profile
}

func (u User) Validate(_ context.Context) error {
	if u.Username == "" {
		return errors.New("missing username field")
	}
	if u.Email == "" {
		return errors.New("missing email field")
	}
	if u.Password == "" {
		return errors.New("missing password field")
	}
	if u.Image == "" {
		return errors.New("missing image field")
	}
	if u.Bio == "" {
		return errors.New("missing bio field")
	}
	return nil
}

func (u *User) Name() string {
	// NOTE using pointer semantic to prevent dereference panic
	return "user"
}

func (u *User) Key() string {
	return u.Email
}

func (u *User) SetKey(id string) {
	u.Email = id
}

type ProfileWrapper struct {
	Profile Profile `json:"profile"`
}

type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}
