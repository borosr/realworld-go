package domain

import (
	"context"
	"errors"
	"testing"

	persistTypes "github.com/borosr/realworld/persist/types"
	"github.com/borosr/realworld/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProfileService_GetByUsername(t *testing.T) {
	const (
		email             = "test@email.com"
		encryptedPassword = "$2a$10$6gS0vGWDr/lp2aOfUPn0me9byfCvnESOKkRg6URJOwzbBMW0zyIJ6"
	)
	ctx := context.Background()
	mockUserRepo := MockRepository[*types.User]{}
	mockFollowRepo := MockRepository[*types.Follow]{}
	expectedUser := types.User{
		Email:    email,
		Password: encryptedPassword,
		Profile: types.Profile{
			Username: email,
		},
	}
	mockUserRepo.On("GetFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.User]) bool {
		return f(&expectedUser)
	})).Return([]*types.User{
		&expectedUser,
	}, nil)
	ps := ProfileService{
		UserRepository:   &mockUserRepo,
		FollowRepository: &mockFollowRepo,
	}
	profile, err := ps.GetByUsername(ctx, email)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, email, profile.Username)
}

func TestProfileService_Follow(t *testing.T) {
	const (
		fromEmail         = "from_test@email.com"
		toEmail           = "test@email.com"
		encryptedPassword = "$2a$10$6gS0vGWDr/lp2aOfUPn0me9byfCvnESOKkRg6URJOwzbBMW0zyIJ6"
	)
	ctx := context.Background()
	mockUserRepo := MockRepository[*types.User]{}
	mockFollowRepo := MockRepository[*types.Follow]{}
	expectedUser := types.User{
		Email:    toEmail,
		Password: encryptedPassword,
		Profile: types.Profile{
			Username: toEmail,
		},
	}
	mockUserRepo.On("GetFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.User]) bool {
		return f(&expectedUser)
	})).Return([]*types.User{
		&expectedUser,
	}, nil)
	expectedFollowed := types.Follow{
		From: fromEmail,
		To:   toEmail,
	}
	mockFollowRepo.On("Get", ctx, expectedFollowed.Key()).
		Return(nil, errors.New("not found"))
	mockFollowRepo.On("Save", ctx, &expectedFollowed).
		Return(&expectedFollowed, nil)
	ps := ProfileService{
		UserRepository:   &mockUserRepo,
		FollowRepository: &mockFollowRepo,
	}
	followed, err := ps.Follow(ctx, fromEmail, toEmail)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, toEmail, followed.Username)
	assert.True(t, followed.Following)
}

func TestProfileService_FollowAlreadyFollowed(t *testing.T) {
	const (
		fromEmail         = "from_test@email.com"
		toEmail           = "test@email.com"
		encryptedPassword = "$2a$10$6gS0vGWDr/lp2aOfUPn0me9byfCvnESOKkRg6URJOwzbBMW0zyIJ6"
	)
	ctx := context.Background()
	mockUserRepo := MockRepository[*types.User]{}
	mockFollowRepo := MockRepository[*types.Follow]{}
	expectedUser := types.User{
		Email:    toEmail,
		Password: encryptedPassword,
		Profile: types.Profile{
			Username: toEmail,
		},
	}
	mockUserRepo.On("GetFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.User]) bool {
		return f(&expectedUser)
	})).Return([]*types.User{
		&expectedUser,
	}, nil)
	expectedFollowed := types.Follow{
		From: fromEmail,
		To:   toEmail,
	}
	mockFollowRepo.On("Get", ctx, expectedFollowed.Key()).
		Return(&expectedFollowed, nil)
	ps := ProfileService{
		UserRepository:   &mockUserRepo,
		FollowRepository: &mockFollowRepo,
	}
	_, err := ps.Follow(ctx, fromEmail, toEmail)
	assert.NotNil(t, err)
}

func TestProfileService_Unfollow(t *testing.T) {
	const (
		fromEmail         = "from_test@email.com"
		toEmail           = "test@email.com"
		encryptedPassword = "$2a$10$6gS0vGWDr/lp2aOfUPn0me9byfCvnESOKkRg6URJOwzbBMW0zyIJ6"
	)
	ctx := context.Background()
	mockUserRepo := MockRepository[*types.User]{}
	mockFollowRepo := MockRepository[*types.Follow]{}
	expectedUser := types.User{
		Email:    toEmail,
		Password: encryptedPassword,
		Profile: types.Profile{
			Username: toEmail,
		},
	}
	mockUserRepo.On("GetFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.User]) bool {
		return f(&expectedUser)
	})).Return([]*types.User{
		&expectedUser,
	}, nil)
	expectedFollowed := types.Follow{
		From: fromEmail,
		To:   toEmail,
	}
	mockFollowRepo.On("Get", ctx, expectedFollowed.Key()).
		Return(&expectedFollowed, nil)
	mockFollowRepo.On("Delete", ctx, expectedFollowed.Key()).
		Return(nil)
	ps := ProfileService{
		UserRepository:   &mockUserRepo,
		FollowRepository: &mockFollowRepo,
	}
	followed, err := ps.Unfollow(ctx, fromEmail, toEmail)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, toEmail, followed.Username)
	assert.False(t, followed.Following)
}

func TestProfileService_UnfollowNotFollowedYet(t *testing.T) {
	const (
		fromEmail         = "from_test@email.com"
		toEmail           = "test@email.com"
		encryptedPassword = "$2a$10$6gS0vGWDr/lp2aOfUPn0me9byfCvnESOKkRg6URJOwzbBMW0zyIJ6"
	)
	ctx := context.Background()
	mockUserRepo := MockRepository[*types.User]{}
	mockFollowRepo := MockRepository[*types.Follow]{}
	expectedUser := types.User{
		Email:    toEmail,
		Password: encryptedPassword,
		Profile: types.Profile{
			Username: toEmail,
		},
	}
	mockUserRepo.On("GetFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.User]) bool {
		return f(&expectedUser)
	})).Return([]*types.User{
		&expectedUser,
	}, nil)
	expectedFollowed := types.Follow{
		From: fromEmail,
		To:   toEmail,
	}
	mockFollowRepo.On("Get", ctx, expectedFollowed.Key()).
		Return(nil, errors.New("not found"))
	ps := ProfileService{
		UserRepository:   &mockUserRepo,
		FollowRepository: &mockFollowRepo,
	}
	_, err := ps.Unfollow(ctx, fromEmail, toEmail)
	assert.NotNil(t, err)
}
