package domain

import (
	"context"
	"errors"
	"testing"

	persistTypes "github.com/borosr/realworld/persist/types"
	"github.com/borosr/realworld/types"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestArticleService_GetAll(t *testing.T) {
	const (
		email = "test@email.com"
	)
	expectedSlug := "test-slug" + xid.New().String()

	ctx := context.Background()

	expectedFavorite := types.Favorite{
		Slug:     expectedSlug,
		Username: email,
	}

	expectedArticle := types.Article{
		Slug: expectedSlug,
	}

	mockArticleRepo := MockRepository[*types.Article]{}
	mockCommentRepo := MockRepository[*types.Comment]{}
	mockFavoriteRepo := MockRepository[*types.Favorite]{}
	mockFavoriteRepo.On("CountFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Favorite]) bool {
		return f(&expectedFavorite)
	})).
		Return(1, nil)
	mockArticleRepo.On("GetFiltered", ctx).
		Return([]*types.Article{&expectedArticle}, nil)
	service := MockUserService{}
	as := ArticleService{
		ArticleRepository:  &mockArticleRepo,
		CommentRepository:  &mockCommentRepo,
		FavoriteRepository: &mockFavoriteRepo,
		UserService:        &service,
	}
	articles, total, err := as.GetAll(ctx, "", "", "", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(articles))
	if !t.Failed() {
		assert.Equal(t, expectedSlug, articles[0].Slug)
	}
}

func TestArticleService_GetAllTagFilter(t *testing.T) {
	const (
		email       = "test@email.com"
		expectedTag = "first_tag"
	)
	expectedSlug := "test-slug" + xid.New().String()

	ctx := context.Background()

	expectedFavorite := types.Favorite{
		Slug:     expectedSlug,
		Username: email,
	}

	expectedArticle := types.Article{
		Slug:    expectedSlug,
		TagList: []string{expectedTag},
	}

	mockArticleRepo := MockRepository[*types.Article]{}
	mockCommentRepo := MockRepository[*types.Comment]{}
	mockFavoriteRepo := MockRepository[*types.Favorite]{}
	mockFavoriteRepo.On("CountFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Favorite]) bool {
		return f(&expectedFavorite)
	})).
		Return(1, nil)
	mockArticleRepo.On("GetFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Article]) bool {
		return f(&expectedArticle)
	})).
		Return([]*types.Article{&expectedArticle}, nil)
	service := MockUserService{}
	as := ArticleService{
		ArticleRepository:  &mockArticleRepo,
		CommentRepository:  &mockCommentRepo,
		FavoriteRepository: &mockFavoriteRepo,
		UserService:        &service,
	}
	articles, total, err := as.GetAll(ctx, expectedTag, "", "", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(articles))
	if !t.Failed() {
		assert.Equal(t, expectedSlug, articles[0].Slug)
	}
}

func TestArticleService_GetAllAuthorFilter(t *testing.T) {
	const (
		email = "test@email.com"
	)
	expectedSlug := "test-slug" + xid.New().String()

	ctx := context.Background()

	expectedFavorite := types.Favorite{
		Slug:     expectedSlug,
		Username: email,
	}

	expectedArticle := types.Article{
		Slug: expectedSlug,
		Author: types.Profile{
			Username: email,
		},
	}

	mockArticleRepo := MockRepository[*types.Article]{}
	mockCommentRepo := MockRepository[*types.Comment]{}
	mockFavoriteRepo := MockRepository[*types.Favorite]{}
	mockFavoriteRepo.On("CountFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Favorite]) bool {
		return f(&expectedFavorite)
	})).
		Return(1, nil)
	mockArticleRepo.On("GetFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Article]) bool {
		return f(&expectedArticle)
	})).
		Return([]*types.Article{&expectedArticle}, nil)
	service := MockUserService{}
	as := ArticleService{
		ArticleRepository:  &mockArticleRepo,
		CommentRepository:  &mockCommentRepo,
		FavoriteRepository: &mockFavoriteRepo,
		UserService:        &service,
	}
	articles, total, err := as.GetAll(ctx, "", email, "", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(articles))
	if !t.Failed() {
		assert.Equal(t, expectedSlug, articles[0].Slug)
	}
}

func TestArticleService_GetAllFavoriteFilter(t *testing.T) {
	const (
		email = "test@email.com"
	)
	expectedSlug := "test-slug" + xid.New().String()

	ctx := context.Background()

	expectedFavorite := types.Favorite{
		Slug:     expectedSlug,
		Username: email,
	}

	expectedArticle := types.Article{
		Slug: expectedSlug,
		Author: types.Profile{
			Username: email,
		},
	}

	mockArticleRepo := MockRepository[*types.Article]{}
	mockCommentRepo := MockRepository[*types.Comment]{}
	mockFavoriteRepo := MockRepository[*types.Favorite]{}
	mockFavoriteRepo.On("CountFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Favorite]) bool {
		return f(&expectedFavorite)
	})).
		Return(1, nil)
	mockFavoriteRepo.On("GetFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Favorite]) bool {
		return f(&expectedFavorite)
	})).
		Return([]*types.Favorite{&expectedFavorite}, nil)
	mockArticleRepo.On("GetFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Article]) bool {
		return f(&expectedArticle)
	})).
		Return([]*types.Article{&expectedArticle}, nil)
	service := MockUserService{}
	as := ArticleService{
		ArticleRepository:  &mockArticleRepo,
		CommentRepository:  &mockCommentRepo,
		FavoriteRepository: &mockFavoriteRepo,
		UserService:        &service,
	}
	articles, total, err := as.GetAll(ctx, "", "", email, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(articles))
	if !t.Failed() {
		assert.Equal(t, expectedSlug, articles[0].Slug)
	}
}

func TestArticleService_GetAllAllFilter(t *testing.T) {
	const (
		email       = "test@email.com"
		expectedTag = "first_tag"
	)
	expectedSlug := "test-slug" + xid.New().String()

	ctx := context.Background()

	expectedFavorite := types.Favorite{
		Slug:     expectedSlug,
		Username: email,
	}

	expectedArticle := types.Article{
		Slug:    expectedSlug,
		TagList: []string{expectedTag},
		Author: types.Profile{
			Username: email,
		},
	}

	mockArticleRepo := MockRepository[*types.Article]{}
	mockCommentRepo := MockRepository[*types.Comment]{}
	mockFavoriteRepo := MockRepository[*types.Favorite]{}
	mockFavoriteRepo.On("CountFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Favorite]) bool {
		return f(&expectedFavorite)
	})).
		Return(1, nil)
	mockFavoriteRepo.On("GetFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Favorite]) bool {
		return f(&expectedFavorite)
	})).
		Return([]*types.Favorite{&expectedFavorite}, nil)
	filterMatcher := mock.MatchedBy(func(f persistTypes.Filter[*types.Article]) bool {
		return f(&expectedArticle)
	})
	mockArticleRepo.On("GetFiltered", ctx, filterMatcher, filterMatcher, filterMatcher).
		Return([]*types.Article{&expectedArticle}, nil)
	service := MockUserService{}
	as := ArticleService{
		ArticleRepository:  &mockArticleRepo,
		CommentRepository:  &mockCommentRepo,
		FavoriteRepository: &mockFavoriteRepo,
		UserService:        &service,
	}
	articles, total, err := as.GetAll(ctx, expectedTag, email, email, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(articles))
	if !t.Failed() {
		assert.Equal(t, expectedSlug, articles[0].Slug)
	}
}

func TestArticleService_Feed(t *testing.T) {
	const (
		email = "test@email.com"
	)
	expectedSlug := "test-slug" + xid.New().String()

	ctx := context.Background()

	expectedFavorite := types.Favorite{
		Slug:     expectedSlug,
		Username: email,
	}

	expectedArticle := types.Article{
		Slug: expectedSlug,
	}

	mockArticleRepo := MockRepository[*types.Article]{}
	mockCommentRepo := MockRepository[*types.Comment]{}
	mockFavoriteRepo := MockRepository[*types.Favorite]{}
	mockFavoriteRepo.On("CountFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Favorite]) bool {
		return f(&expectedFavorite)
	})).
		Return(1, nil)
	mockArticleRepo.On("GetFiltered", ctx).
		Return([]*types.Article{&expectedArticle}, nil)
	service := MockUserService{}
	as := ArticleService{
		ArticleRepository:  &mockArticleRepo,
		CommentRepository:  &mockCommentRepo,
		FavoriteRepository: &mockFavoriteRepo,
		UserService:        &service,
	}
	articles, total, err := as.Feed(ctx, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, total)
	assert.Equal(t, 1, len(articles))
	if !t.Failed() {
		assert.Equal(t, expectedSlug, articles[0].Slug)
	}
}

func TestArticleService_Create(t *testing.T) {
	const (
		email         = "test@email.com"
		expectedTitle = "random_title"
	)
	expectedSlug := "test-slug" + xid.New().String()

	ctx := context.Background()

	expectedArticle := types.Article{
		Slug:    expectedSlug,
		Title:   expectedTitle,
		TagList: []string{"a", "b", "c"},
	}

	mockArticleRepo := MockRepository[*types.Article]{}
	mockCommentRepo := MockRepository[*types.Comment]{}
	mockFavoriteRepo := MockRepository[*types.Favorite]{}
	mockArticleRepo.On("Save", ctx, mock.MatchedBy(func(a *types.Article) bool {
		return a.Title == expectedArticle.Title &&
			len(a.TagList) == 3 &&
			a.TagList[0] == "a" &&
			a.TagList[1] == "b" &&
			a.TagList[2] == "c"
	})).
		Return(&expectedArticle, nil)
	service := MockUserService{}
	service.On("GetByEmail", ctx, email).
		Return(types.User{
			Profile: types.Profile{
				Username: email,
			},
		}, nil)
	as := ArticleService{
		ArticleRepository:  &mockArticleRepo,
		CommentRepository:  &mockCommentRepo,
		FavoriteRepository: &mockFavoriteRepo,
		UserService:        &service,
	}
	article, err := as.Create(ctx, types.ArticleRequest{
		Title:   expectedTitle,
		TagList: []string{"b", "c", "a"},
	}, email)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedSlug, article.Slug)
	assert.NotEmpty(t, article.Slug)
	assert.Equal(t, []string{"a", "b", "c"}, article.TagList)
}

func TestArticleService_Update(t *testing.T) {
	const (
		email               = "test@email.com"
		expectedTitle       = "random_title"
		expectedDescription = "random_description"
		expectedBody        = "content_1"
	)
	expectedSlug := "test-slug" + xid.New().String()

	ctx := context.Background()

	expectedArticle := types.Article{
		Slug:        expectedSlug,
		Title:       expectedTitle,
		Description: expectedDescription,
		Body:        expectedBody,
		TagList:     []string{"a", "b", "c"},
	}

	mockArticleRepo := MockRepository[*types.Article]{}
	mockCommentRepo := MockRepository[*types.Comment]{}
	mockFavoriteRepo := MockRepository[*types.Favorite]{}
	mockArticleRepo.On("Get", ctx, expectedSlug).
		Return(&expectedArticle, nil)
	mockArticleRepo.On("Save", ctx, mock.MatchedBy(func(a *types.Article) bool {
		return a.Title == expectedArticle.Title &&
			len(a.TagList) == 3 &&
			a.TagList[0] == "a" &&
			a.TagList[1] == "b" &&
			a.TagList[2] == "c"
	})).
		Return(&expectedArticle, nil)
	service := MockUserService{}
	service.On("GetByEmail", ctx, email).
		Return(types.User{
			Profile: types.Profile{
				Username: email,
			},
		}, nil)
	as := ArticleService{
		ArticleRepository:  &mockArticleRepo,
		CommentRepository:  &mockCommentRepo,
		FavoriteRepository: &mockFavoriteRepo,
		UserService:        &service,
	}
	article, err := as.Update(ctx, expectedSlug, types.ArticleRequest{
		Title:       expectedTitle,
		Description: expectedDescription,
		Body:        expectedBody,
		TagList:     []string{"b", "c", "a"},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedSlug, article.Slug)
	assert.NotEmpty(t, article.Slug)
	assert.Equal(t, []string{"a", "b", "c"}, article.TagList)
}

func TestArticleService_CreateComment(t *testing.T) {
	const (
		email        = "test@email.com"
		expectedBody = "random comment"
	)
	expectedSlug := "test-slug" + xid.New().String()

	ctx := context.WithValue(context.Background(), "email", email)

	expectedComment := types.Comment{
		Slug: expectedSlug,
		CommonComment: types.CommonComment{
			ID:   1,
			Body: expectedBody,
			Author: types.Profile{
				Username: email,
			},
		},
	}

	mockArticleRepo := MockRepository[*types.Article]{}
	mockCommentRepo := MockRepository[*types.Comment]{}
	mockFavoriteRepo := MockRepository[*types.Favorite]{}
	mockCommentRepo.On("Save", ctx, mock.MatchedBy(func(a *types.Comment) bool {
		return a.Slug == expectedSlug
	})).
		Return(&expectedComment, nil)
	mockCommentRepo.On("Sequence", ctx, expectedSlug).
		Return(1, nil)
	service := MockUserService{}
	service.On("GetByEmail", ctx, email).
		Return(types.User{
			Profile: types.Profile{
				Username: email,
			},
		}, nil)
	as := ArticleService{
		ArticleRepository:  &mockArticleRepo,
		CommentRepository:  &mockCommentRepo,
		FavoriteRepository: &mockFavoriteRepo,
		UserService:        &service,
	}
	comment, err := as.CreateComment(ctx, expectedSlug, types.CommentRequest{
		Body: expectedBody,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, email, comment.Author.Username)
	assert.NotEmpty(t, comment.Body)
}

func TestArticleService_AddFavoriteArticle(t *testing.T) {
	const (
		email               = "test@email.com"
		expectedTitle       = "random_title"
		expectedDescription = "random_description"
		expectedBody        = "content_1"
	)
	expectedSlug := "test-slug" + xid.New().String()

	ctx := context.Background()

	expectedArticle := types.Article{
		Slug:        expectedSlug,
		Title:       expectedTitle,
		Description: expectedDescription,
		Body:        expectedBody,
		TagList:     []string{"a", "b", "c"},
	}

	expectedFavorite := types.Favorite{
		Slug:     expectedSlug,
		Username: email,
	}

	mockArticleRepo := MockRepository[*types.Article]{}
	mockCommentRepo := MockRepository[*types.Comment]{}
	mockFavoriteRepo := MockRepository[*types.Favorite]{}
	mockArticleRepo.On("Get", ctx, expectedSlug).
		Return(&expectedArticle, nil)
	mockFavoriteRepo.On("Get", ctx, expectedFavorite.Key()).
		Return(&expectedFavorite, errors.New("not found"))
	mockFavoriteRepo.On("Save", ctx, mock.MatchedBy(func(a *types.Favorite) bool {
		return a.Slug == expectedSlug && a.Username == email
	})).
		Return(&expectedFavorite, nil)
	mockFavoriteRepo.On("CountFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Favorite]) bool {
		return f(&expectedFavorite)
	})).
		Return(1, nil)
	service := MockUserService{}
	service.On("GetByEmail", ctx, email).
		Return(types.User{
			Profile: types.Profile{
				Username: email,
			},
		}, nil)
	as := ArticleService{
		ArticleRepository:  &mockArticleRepo,
		CommentRepository:  &mockCommentRepo,
		FavoriteRepository: &mockFavoriteRepo,
		UserService:        &service,
	}
	article, err := as.AddFavoriteArticle(ctx, expectedSlug, email)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedSlug, article.Slug)
	assert.NotEmpty(t, article.Slug)
	assert.True(t, article.Favorited)
	assert.Equal(t, 1, article.FavoritesCount)
}

func TestArticleService_DeleteFavoriteArticle(t *testing.T) {
	const (
		email               = "test@email.com"
		expectedTitle       = "random_title"
		expectedDescription = "random_description"
		expectedBody        = "content_1"
	)
	expectedSlug := "test-slug" + xid.New().String()

	ctx := context.Background()

	expectedArticle := types.Article{
		Slug:        expectedSlug,
		Title:       expectedTitle,
		Description: expectedDescription,
		Body:        expectedBody,
		TagList:     []string{"a", "b", "c"},
	}

	expectedFavorite := types.Favorite{
		Slug:     expectedSlug,
		Username: email,
	}

	mockArticleRepo := MockRepository[*types.Article]{}
	mockCommentRepo := MockRepository[*types.Comment]{}
	mockFavoriteRepo := MockRepository[*types.Favorite]{}
	mockArticleRepo.On("Get", ctx, expectedSlug).
		Return(&expectedArticle, nil)
	mockFavoriteRepo.On("Delete", ctx, expectedFavorite.Key()).
		Return(nil)
	mockFavoriteRepo.On("CountFiltered", ctx, mock.MatchedBy(func(f persistTypes.Filter[*types.Favorite]) bool {
		return f(&expectedFavorite)
	})).
		Return(1, nil)
	service := MockUserService{}
	service.On("GetByEmail", ctx, email).
		Return(types.User{
			Profile: types.Profile{
				Username: email,
			},
		}, nil)
	as := ArticleService{
		ArticleRepository:  &mockArticleRepo,
		CommentRepository:  &mockCommentRepo,
		FavoriteRepository: &mockFavoriteRepo,
		UserService:        &service,
	}
	article, err := as.DeleteFavoriteArticle(ctx, expectedSlug, email)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expectedSlug, article.Slug)
	assert.NotEmpty(t, article.Slug)
	assert.False(t, article.Favorited)
	assert.Equal(t, 1, article.FavoritesCount)
}
