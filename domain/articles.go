package domain

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/borosr/realworld/lib/api"
	"github.com/borosr/realworld/persist"
	persistTypes "github.com/borosr/realworld/persist/types"
	"github.com/borosr/realworld/types"
)

type ArticleDescriptor interface {
	GetAll(ctx context.Context, tag, author, favorite string, limit, offset int) ([]*types.Article, int, error)
	Feed(ctx context.Context, limit, offset int) ([]*types.Article, int, error)
	Get(ctx context.Context, slug string) (types.Article, error)
	Create(ctx context.Context, a types.ArticleRequest, ownerEmail string) (types.Article, error)
	Update(ctx context.Context, slug string, a types.ArticleRequest) (types.Article, error)
	Delete(ctx context.Context, slug string) error
	CreateComment(ctx context.Context, slug string, c types.CommentRequest) (types.CommonComment, error)
	GetComments(ctx context.Context, slug string) ([]types.CommonComment, error)
	DeleteComment(ctx context.Context, slug string, id int) error
	AddFavoriteArticle(ctx context.Context, slug, username string) (types.Article, error)
	DeleteFavoriteArticle(ctx context.Context, slug, username string) (types.Article, error)
}

type ArticleService struct {
	ArticleRepository  persist.Repository[*types.Article]
	CommentRepository  persist.Repository[*types.Comment]
	FavoriteRepository persist.Repository[*types.Favorite]
	UserService        UserDescriptor
}

func (as ArticleService) GetAll(ctx context.Context, tag, author, favorite string, limit, offset int) ([]*types.Article, int, error) {
	var filters []persistTypes.Filter[*types.Article]
	if tag != "" {
		filters = append(filters, func(t *types.Article) bool {
			for _, t := range t.TagList {
				if t == tag {
					return true
				}
			}
			return false
		})
	}
	if author != "" {
		filters = append(filters, func(t *types.Article) bool {
			return t.Author.Username == author
		})
	}
	if favorite != "" {
		filtered, err := as.FavoriteRepository.GetFiltered(ctx, func(f *types.Favorite) bool {
			return f.Username == favorite
		})
		if err != nil {
			return nil, 0, err
		}
		filters = append(filters, func(t *types.Article) bool {
			for i := range filtered {
				if t.Slug == filtered[i].Slug {
					return true
				}
			}
			return false
		})
	}
	results, err := as.ArticleRepository.GetFiltered(ctx, filters...)
	if err != nil {
		return nil, 0, err
	}
	results, totalCount := as.reduceResult(results, limit, offset)
	as.attachFavorite(ctx, results)
	return results, totalCount, nil
}

func (as ArticleService) Feed(ctx context.Context, limit, offset int) ([]*types.Article, int, error) {
	results, err := as.ArticleRepository.GetFiltered(ctx)
	if err != nil {
		return nil, 0, err
	}
	results, totalCount := as.reduceResult(results, limit, offset)
	as.attachFavorite(ctx, results)
	return results, totalCount, nil
}

func (as ArticleService) Get(ctx context.Context, slug string) (types.Article, error) {
	article, err := as.ArticleRepository.Get(ctx, slug)
	if err != nil {
		return types.Article{}, err
	}
	return *article, nil
}

func (as ArticleService) Create(ctx context.Context, a types.ArticleRequest, ownerEmail string) (types.Article, error) {
	user, err := as.UserService.GetByEmail(ctx, ownerEmail)
	if err != nil {
		return types.Article{}, err
	}
	now := time.Now()
	sort.Strings(a.TagList)
	saved, err := as.ArticleRepository.Save(ctx, &types.Article{
		Title:       a.Title,
		Description: a.Description,
		Body:        a.Body,
		TagList:     a.TagList,
		CreatedAt:   now,
		UpdatedAt:   now,
		Author:      user.Profile,
	})
	if err != nil {
		return types.Article{}, err
	}
	return *saved, nil
}

func (as ArticleService) Update(ctx context.Context, slug string, a types.ArticleRequest) (types.Article, error) {
	existing, err := as.ArticleRepository.Get(ctx, slug)
	if err != nil {
		return types.Article{}, err
	}
	if a.Title != "" {
		existing.Title = a.Title
	}
	if a.Description != "" {
		existing.Description = a.Description
	}
	if a.Body != "" {
		existing.Body = a.Body
	}
	if a.TagList != nil {
		existing.TagList = a.TagList
	}
	updated, err := as.ArticleRepository.Save(ctx, existing)
	if err != nil {
		return types.Article{}, err
	}
	return *updated, nil
}

func (as ArticleService) Delete(ctx context.Context, slug string) error {
	if err := as.ArticleRepository.Delete(ctx, slug); err != nil {
		return err
	}
	return nil
}

func (as ArticleService) CreateComment(ctx context.Context, slug string, c types.CommentRequest) (types.CommonComment, error) {
	var commentID int
	seq, err := as.CommentRepository.Sequence(ctx, slug)
	if err != nil {
		return types.CommonComment{}, err
	}
	commentID = int(seq)
	email, err := api.GetValue[string](ctx, "email")
	if err != nil {
		return types.CommonComment{}, err
	}
	user, err := as.UserService.GetByEmail(ctx, email)
	if err != nil {
		return types.CommonComment{}, err
	}
	now := time.Now()
	saved, err := as.CommentRepository.Save(ctx, &types.Comment{
		CommonComment: types.CommonComment{
			ID:        commentID,
			CreatedAt: now,
			UpdatedAt: now,
			Body:      c.Body,
			Author:    user.Profile,
		},
		Slug: slug,
	})
	if err != nil {
		return types.CommonComment{}, err
	}
	return saved.CommonComment, nil
}

func (as ArticleService) GetComments(ctx context.Context, slug string) ([]types.CommonComment, error) {
	results, err := as.CommentRepository.GetFiltered(ctx, func(c *types.Comment) bool {
		return c.Slug == slug
	})
	if err != nil {
		return nil, err
	}
	var comments []types.CommonComment
	for _, res := range results {
		comments = append(comments, res.CommonComment)
	}
	return comments, nil
}

func (as ArticleService) DeleteComment(ctx context.Context, slug string, id int) error {
	c := types.Comment{
		CommonComment: types.CommonComment{
			ID: id,
		},
		Slug: slug,
	}
	return as.CommentRepository.Delete(ctx, c.Key())
}

func (as ArticleService) AddFavoriteArticle(ctx context.Context, slug, email string) (types.Article, error) {
	article, err := as.ArticleRepository.Get(ctx, slug)
	if err != nil {
		return types.Article{}, err
	}
	user, err := as.UserService.GetByEmail(ctx, email)
	if err != nil {
		return types.Article{}, err
	}
	favorite := types.Favorite{
		Slug:     slug,
		Username: user.Username,
	}
	if _, err := as.FavoriteRepository.Get(ctx, favorite.Key()); err == nil {
		return types.Article{}, errors.New("already added to favorite")
	}
	if _, err := as.FavoriteRepository.Save(ctx, &favorite); err != nil {
		return types.Article{}, err
	}
	article.Favorited = true
	favoriteCount, err := as.FavoriteRepository.CountFiltered(ctx, func(f *types.Favorite) bool {
		return f.Slug == slug
	})
	if err != nil {
		return types.Article{}, err
	}
	article.FavoritesCount = int(favoriteCount)
	return *article, nil
}

func (as ArticleService) DeleteFavoriteArticle(ctx context.Context, slug, email string) (types.Article, error) {
	user, err := as.UserService.GetByEmail(ctx, email)
	if err != nil {
		return types.Article{}, err
	}
	article, err := as.ArticleRepository.Get(ctx, slug)
	if err != nil {
		return types.Article{}, err
	}
	favorite := types.Favorite{
		Slug:     slug,
		Username: user.Username,
	}
	article.Favorited = false
	if err := as.FavoriteRepository.Delete(ctx, favorite.Key()); err != nil {
		return types.Article{}, err
	}
	return *article, nil
}

func (_ ArticleService) reduceResult(results []*types.Article, limit int, offset int) ([]*types.Article, int) {
	totalCount := len(results)
	if totalCount > limit+offset {
		results = results[offset : limit+offset]
	} else {
		results = results[offset:]
	}
	return results, totalCount
}

func (as ArticleService) attachFavorite(ctx context.Context, results []*types.Article) {
	for i := range results {
		favorites, _ := as.FavoriteRepository.CountFiltered(ctx, func(f *types.Favorite) bool {
			return f.Slug == results[i].Slug
		})
		results[i].FavoritesCount = int(favorites)
	}
}
