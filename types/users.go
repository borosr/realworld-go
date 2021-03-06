package types

import (
	"context"

	"github.com/borosr/realworld/lib/api"
	"github.com/borosr/realworld/lib/broken"
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
		return broken.Validation("missing email field")
	}
	if u.Password == "" {
		return broken.Validation("missing password field")
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
		return broken.Validation("missing username field")
	}
	if u.Email == "" {
		return broken.Validation("missing email field")
	}
	if u.Password == "" {
		return broken.Validation("missing password field")
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
		return broken.Validation("missing username field")
	}
	if u.Email == "" {
		return broken.Validation("missing email field")
	}
	if u.Password == "" {
		return broken.Validation("missing password field")
	}
	if u.Image == "" {
		return broken.Validation("missing image field")
	}
	if u.Bio == "" {
		return broken.Validation("missing bio field")
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

type Follow struct {
	From string `json:"followed_by"`
	To   string `json:"following"`
}

func (f *Follow) Name() string {
	return "follow"
}

func (f *Follow) Key() string {
	return f.From + "-" + f.To
}

func (f *Follow) SetKey(id string) {
	// DO NOTHING
}
