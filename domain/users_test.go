package domain

import (
	"context"
	"testing"

	"github.com/borosr/realworld/lib/auth"
	"github.com/borosr/realworld/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Login(t *testing.T) {
	const (
		email             = "test@email.com"
		encryptedPassword = "$2a$10$6gS0vGWDr/lp2aOfUPn0me9byfCvnESOKkRg6URJOwzbBMW0zyIJ6"
	)
	ctx := context.Background()
	mock := MockRepository[*types.User]{}
	expectedUser := types.User{
		Email:    email,
		Password: encryptedPassword,
		Profile: types.Profile{
			Username: email,
		},
	}
	mock.On("Get", ctx, email).Return(&expectedUser, nil)
	mock.On("Save", ctx, &expectedUser).Return(&expectedUser, nil)
	us := UserService{
		UserRepository: &mock,
	}
	login, err := us.Login(ctx, types.UserLogin{
		Email:    email,
		Password: "password",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, login.Token)
	verify, err := auth.Verify(login.Token)
	assert.Nil(t, err)
	assert.Equal(t, email, verify["email"])
}

func TestUserService_SignUp(t *testing.T) {
	const (
		username          = "some_username"
		email             = "test@email.com"
		password          = "password"
		encryptedPassword = "$2a$10$6gS0vGWDr/lp2aOfUPn0me9byfCvnESOKkRg6URJOwzbBMW0zyIJ6"
	)
	ctx := context.Background()
	mockRepo := MockRepository[*types.User]{}
	expectedUser := types.User{
		Email:    email,
		Password: encryptedPassword,
		Profile: types.Profile{
			Username: username,
		},
	}
	mockRepo.On("Save", ctx, mock.MatchedBy(func(u *types.User) bool {
		return u.Email == email &&
			u.Username == username &&
			u.Password != ""
	})).Return(&expectedUser, nil)
	us := UserService{
		UserRepository: &mockRepo,
	}
	u, err := us.SignUp(ctx, types.UserSignUp{
		Username: username,
		Email:    email,
		Password: password,
	})
	assert.Nil(t, err)
	assert.Equal(t, email, u.Email)
	assert.Equal(t, username, u.Username)
	assert.NotEmpty(t, u.Password)
}

func TestUserService_Update(t *testing.T) {
	const (
		email             = "test@email.com"
		encryptedPassword = "$2a$10$6gS0vGWDr/lp2aOfUPn0me9byfCvnESOKkRg6URJOwzbBMW0zyIJ6"
		expectedBio       = "random bio description"
		expectedToken     = "token1234"
		expectedImageURL  = "localhost/test_img.jpg"
	)
	ctx := context.Background()
	mock := MockRepository[*types.User]{}
	expectedUser := types.User{
		Email:    email,
		Password: encryptedPassword,
		Profile: types.Profile{
			Username: email,
		},
	}
	mock.On("Get", ctx, email).Return(&expectedUser, nil)
	us := UserService{
		UserRepository: &mock,
	}
	changedUser := *(&expectedUser)
	changedUser.Bio = expectedBio
	changedUser.Token = expectedToken
	changedUser.Image = expectedImageURL
	mock.On("Save", ctx, &changedUser).Return(&changedUser, nil)

	updated, err := us.Update(ctx, changedUser)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedBio, updated.Bio)
}
