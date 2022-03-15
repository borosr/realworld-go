package api

import (
	"context"
	"errors"
	goTypes "go/types"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/borosr/realworld/lib/api"
	"github.com/borosr/realworld/lib/middleware"
	"github.com/borosr/realworld/persist"
	persistTypes "github.com/borosr/realworld/persist/types"
	"github.com/borosr/realworld/types"
)

type articlesController struct {
	articleRepository  persist.Repository[*types.Article]
	commentRepository  persist.Repository[*types.Comment]
	favoriteRepository persist.Repository[*types.Favorite]
	userRepository     persist.Repository[*types.User]
}

func (ac articlesController) Init() {
	api.Register[
		goTypes.Nil,
		types.ArticleListResponseWrapper,
		api.ControllerFunc[goTypes.Nil, types.ArticleListResponseWrapper],
	]("/api/articles", http.MethodGet, ac.getAll)
	api.Register[
		goTypes.Nil,
		types.ArticleListResponseWrapper,
		api.ControllerFunc[goTypes.Nil, types.ArticleListResponseWrapper],
	]("/api/articles/feed", http.MethodGet, ac.feed).
		PreProcess(middleware.TokenAuthentication)
	api.Register[
		goTypes.Nil,
		types.ArticleWrapper[types.Article],
		api.ControllerSimpleFunc[goTypes.Nil, types.ArticleWrapper[types.Article]],
	]("/api/articles/{slug}", http.MethodGet, ac.get)
	api.Register[
		types.ArticleWrapper[types.ArticleRequest],
		types.ArticleWrapper[types.Article],
		api.ControllerSimpleFunc[types.ArticleWrapper[types.ArticleRequest], types.ArticleWrapper[types.Article]],
	]("/api/articles", http.MethodPost, ac.create).
		PreProcess(middleware.TokenAuthentication).
		Validated()
	api.Register[
		types.ArticleWrapper[types.ArticleRequest],
		types.ArticleWrapper[types.Article],
		api.ControllerSimpleFunc[types.ArticleWrapper[types.ArticleRequest], types.ArticleWrapper[types.Article]],
	]("/api/articles/{slug}", http.MethodPut, ac.update).
		PreProcess(middleware.TokenAuthentication).
		Validated()
	api.Register[
		goTypes.Nil,
		goTypes.Nil,
		api.ControllerSimpleFunc[goTypes.Nil, goTypes.Nil],
	]("/api/articles/{slug}", http.MethodDelete, ac.delete).
		PreProcess(middleware.TokenAuthentication)
	api.Register[
		types.CommentWrapper[types.CommentRequest],
		types.CommentWrapper[types.CommonComment],
		api.ControllerSimpleFunc[types.CommentWrapper[types.CommentRequest], types.CommentWrapper[types.CommonComment]],
	]("/api/articles/{slug}/comments", http.MethodPost, ac.createComment).
		PreProcess(middleware.TokenAuthentication).
		Validated()
	api.Register[
		goTypes.Nil,
		types.CommentListResponseWrapper,
		api.ControllerSimpleFunc[goTypes.Nil, types.CommentListResponseWrapper],
	]("/api/articles/{slug}/comments", http.MethodGet, ac.getComments)
	api.Register[
		goTypes.Nil,
		goTypes.Nil,
		api.ControllerSimpleFunc[goTypes.Nil, goTypes.Nil],
	]("/api/articles/{slug}/comments/{id}", http.MethodDelete, ac.deleteComment).
		PreProcess(middleware.TokenAuthentication)
	api.Register[
		goTypes.Nil,
		types.ArticleWrapper[types.Article],
		api.ControllerSimpleFunc[goTypes.Nil, types.ArticleWrapper[types.Article]],
	]("/api/articles/{slug}/favorite", http.MethodPost, ac.addFavoriteArticle).
		PreProcess(middleware.TokenAuthentication)
	api.Register[
		goTypes.Nil,
		types.ArticleWrapper[types.Article],
		api.ControllerSimpleFunc[goTypes.Nil, types.ArticleWrapper[types.Article]],
	]("/api/articles/{slug}/favorite", http.MethodDelete, ac.deleteFavoriteArticle).
		PreProcess(middleware.TokenAuthentication)
}

func (ac articlesController) getAll(ctx context.Context, _ goTypes.Nil, m api.Meta) (types.ArticleListResponseWrapper, error) {
	var filters []persistTypes.Filter[*types.Article]
	if m.Params.Has("tag") {
		param := m.Params.Get("tag")
		filters = append(filters, func(t *types.Article) bool {
			for _, tag := range t.TagList {
				if tag == param {
					return true
				}
			}
			return false
		})
	}
	if m.Params.Has("author") {
		param := m.Params.Get("author")
		filters = append(filters, func(t *types.Article) bool {
			return t.Author.Username == param
		})
	}
	if m.Params.Has("favorited") {
		param := m.Params.Get("favorited")
		filtered, err := ac.favoriteRepository.GetFiltered(ctx, func(f *types.Favorite) bool {
			return f.Username == param
		})
		if err != nil {
			return types.ArticleListResponseWrapper{}, err
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
	limit, offset := ac.getLimitOffset(m)
	results, err := ac.articleRepository.GetFiltered(ctx, filters...)
	if err != nil {
		return types.ArticleListResponseWrapper{}, err
	}
	results, totalCount := ac.reduceResult(results, limit, offset)
	ac.attachFavorite(ctx, results)
	return types.ArticleListResponseWrapper{
		Articles:      results,
		ArticlesCount: totalCount,
	}, nil
}

func (ac articlesController) attachFavorite(ctx context.Context, results []*types.Article) {
	for i := range results {
		favorites, _ := ac.favoriteRepository.CountFiltered(ctx, func(f *types.Favorite) bool {
			return f.Slug == results[i].Slug
		})
		results[i].FavoritesCount = int(favorites)
	}
}

func (ac articlesController) feed(ctx context.Context, _ goTypes.Nil, m api.Meta) (types.ArticleListResponseWrapper, error) {
	limit, offset := ac.getLimitOffset(m)
	results, err := ac.articleRepository.GetFiltered(ctx)
	if err != nil {
		return types.ArticleListResponseWrapper{}, err
	}
	results, totalCount := ac.reduceResult(results, limit, offset)
	ac.attachFavorite(ctx, results)
	return types.ArticleListResponseWrapper{
		Articles:      results,
		ArticlesCount: totalCount,
	}, nil
}

func (ac articlesController) getLimitOffset(m api.Meta) (int, int) {
	limit := 20
	if m.Params.Has("limit") {
		var err error
		limit, err = strconv.Atoi(m.Params.Get("limit"))
		if err != nil {
			limit = 20
		}
	}
	offset := 0
	if m.Params.Has("offset") {
		var err error
		limit, err = strconv.Atoi(m.Params.Get("offset"))
		if err != nil {
			offset = 0
		}
	}
	return limit, offset
}

func (ac articlesController) reduceResult(results []*types.Article, limit int, offset int) ([]*types.Article, int) {
	totalCount := len(results)
	if totalCount > limit+offset {
		results = results[offset : limit+offset]
	} else {
		results = results[offset:]
	}
	return results, totalCount
}

func (ac articlesController) get(ctx context.Context, _ goTypes.Nil) (types.ArticleWrapper[types.Article], error) {
	var fallbackResult types.ArticleWrapper[types.Article]
	id, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return fallbackResult, err
	}
	article, err := ac.articleRepository.Get(ctx, id)
	if err != nil {
		return fallbackResult, err
	}
	fallbackResult.Article = *article
	return fallbackResult, nil
}

func (ac articlesController) create(ctx context.Context, req types.ArticleWrapper[types.ArticleRequest]) (types.ArticleWrapper[types.Article], error) {
	var fallbackResult types.ArticleWrapper[types.Article]
	email, err := api.GetValue[string](ctx, "email")
	if err != nil {
		return fallbackResult, err
	}
	user, err := ac.getUser(ctx, email)
	if err != nil {
		return fallbackResult, err
	}
	now := time.Now()
	sort.Strings(req.Article.TagList)
	saved, err := ac.articleRepository.Save(ctx, &types.Article{
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		TagList:     req.Article.TagList,
		CreatedAt:   now,
		UpdatedAt:   now,
		Author:      user.Profile, // TODO get following from specific table
	})
	if err != nil {
		return fallbackResult, err
	}
	fallbackResult.Article = *saved
	return fallbackResult, nil
}

func (ac articlesController) update(ctx context.Context, req types.ArticleWrapper[types.ArticleRequest]) (types.ArticleWrapper[types.Article], error) {
	var fallbackResult types.ArticleWrapper[types.Article]
	slug, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return fallbackResult, err
	}
	existing, err := ac.articleRepository.Get(ctx, slug)
	if err != nil {
		return fallbackResult, err
	}
	if req.Article.Title != "" {
		existing.Title = req.Article.Title
	}
	if req.Article.Description != "" {
		existing.Description = req.Article.Description
	}
	if req.Article.Body != "" {
		existing.Body = req.Article.Body
	}
	if req.Article.TagList != nil {
		existing.TagList = req.Article.TagList
	}
	updated, err := ac.articleRepository.Save(ctx, existing)
	if err != nil {
		return fallbackResult, err
	}
	fallbackResult.Article = *updated
	return fallbackResult, nil
}

func (ac articlesController) delete(ctx context.Context, _ goTypes.Nil) (goTypes.Nil, error) {
	slug, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return goTypes.Nil{}, err
	}
	if err := ac.articleRepository.Delete(ctx, slug); err != nil {
		return goTypes.Nil{}, err
	}
	return goTypes.Nil{}, nil
}

func (ac articlesController) createComment(ctx context.Context, req types.CommentWrapper[types.CommentRequest]) (types.CommentWrapper[types.CommonComment], error) {
	var fallbackResult types.CommentWrapper[types.CommonComment]
	slug, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return fallbackResult, err
	}
	var commentID int
	seq, err := ac.commentRepository.Sequence(ctx, slug)
	if err != nil {
		return fallbackResult, err
	}
	commentID = int(seq)
	email, err := api.GetValue[string](ctx, "email")
	if err != nil {
		return fallbackResult, err
	}
	user, err := ac.getUser(ctx, email)
	if err != nil {
		return fallbackResult, err
	}
	now := time.Now()
	saved, err := ac.commentRepository.Save(ctx, &types.Comment{
		CommonComment: types.CommonComment{
			ID:        commentID,
			CreatedAt: now,
			UpdatedAt: now,
			Body:      req.Comment.Body,
			Author:    user.Profile,
		},
		Slug: slug,
	})
	if err != nil {
		return fallbackResult, err
	}
	fallbackResult.Comment = saved.CommonComment
	return fallbackResult, nil
}

func (ac articlesController) getComments(ctx context.Context, _ goTypes.Nil) (types.CommentListResponseWrapper, error) {
	slug, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return types.CommentListResponseWrapper{}, err
	}
	results, err := ac.commentRepository.GetFiltered(ctx, func(c *types.Comment) bool {
		return c.Slug == slug
	})
	if err != nil {
		return types.CommentListResponseWrapper{}, err
	}
	var comments []types.CommonComment
	for _, res := range results {
		comments = append(comments, res.CommonComment)
	}
	return types.CommentListResponseWrapper{
		Comments: comments,
	}, nil
}

func (ac articlesController) deleteComment(ctx context.Context, _ goTypes.Nil) (goTypes.Nil, error) {
	slug, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return goTypes.Nil{}, err
	}
	id, err := api.PathVariable[int](ctx, "id")
	if err != nil {
		return goTypes.Nil{}, err
	}
	c := types.Comment{
		CommonComment: types.CommonComment{
			ID: id,
		},
		Slug: slug,
	}
	return goTypes.Nil{}, ac.commentRepository.Delete(ctx, c.Key())
}

func (ac articlesController) addFavoriteArticle(ctx context.Context, _ goTypes.Nil) (types.ArticleWrapper[types.Article], error) {
	var fallbackResult types.ArticleWrapper[types.Article]
	slug, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return fallbackResult, err
	}
	article, err := ac.articleRepository.Get(ctx, slug)
	if err != nil {
		return fallbackResult, err
	}
	email, err := api.GetValue[string](ctx, "email")
	if err != nil {
		return fallbackResult, err
	}
	user, err := ac.getUser(ctx, email)
	if err != nil {
		return fallbackResult, err
	}
	favorite := types.Favorite{
		Slug:     slug,
		Username: user.Username,
	}
	if _, err := ac.favoriteRepository.Get(ctx, favorite.Key()); err == nil {
		return fallbackResult, errors.New("already added to favorite")
	}
	if _, err := ac.favoriteRepository.Save(ctx, &favorite); err != nil {
		return fallbackResult, err
	}
	article.Favorited = true
	favoriteCount, err := ac.favoriteRepository.CountFiltered(ctx, func(f *types.Favorite) bool {
		return f.Slug == slug
	})
	if err != nil {
		return fallbackResult, err
	}
	article.FavoritesCount = int(favoriteCount)
	fallbackResult.Article = *article
	return fallbackResult, nil
}

func (ac articlesController) deleteFavoriteArticle(ctx context.Context, _ goTypes.Nil) (types.ArticleWrapper[types.Article], error) {
	var fallbackResult types.ArticleWrapper[types.Article]
	slug, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return fallbackResult, err
	}
	article, err := ac.articleRepository.Get(ctx, slug)
	if err != nil {
		return fallbackResult, err
	}
	email, err := api.GetValue[string](ctx, "email")
	if err != nil {
		return fallbackResult, err
	}
	user, err := ac.getUser(ctx, email)
	if err != nil {
		return fallbackResult, err
	}
	favorite := types.Favorite{
		Slug:     slug,
		Username: user.Username,
	}
	article.Favorited = false
	if err := ac.favoriteRepository.Delete(ctx, favorite.Key()); err != nil {
		return fallbackResult, err
	}
	fallbackResult.Article = *article
	return fallbackResult, nil
}

func (ac articlesController) getUser(ctx context.Context, email string) (types.User, error) {
	user, err := ac.userRepository.Get(ctx, email)
	if err != nil {
		return types.User{}, err
	}
	return *user, nil
}
