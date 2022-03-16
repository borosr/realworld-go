package api

import (
	"context"
	goTypes "go/types"
	"net/http"
	"strconv"

	"github.com/borosr/realworld/domain"
	"github.com/borosr/realworld/lib/api"
	"github.com/borosr/realworld/lib/middleware"
	"github.com/borosr/realworld/types"
)

type articlesController struct {
	articleService domain.ArticleDescriptor
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
	limit, offset := ac.getLimitOffset(m)
	results, totalCount, err := ac.articleService.GetAll(ctx,
		m.Params.Get("tag"),
		m.Params.Get("author"),
		m.Params.Get("favorited"),
		limit, offset)
	if err != nil {
		return types.ArticleListResponseWrapper{}, err
	}
	return types.ArticleListResponseWrapper{
		Articles:      results,
		ArticlesCount: totalCount,
	}, nil
}

func (ac articlesController) feed(ctx context.Context, _ goTypes.Nil, m api.Meta) (types.ArticleListResponseWrapper, error) {
	limit, offset := ac.getLimitOffset(m)
	results, totalCount, err := ac.articleService.Feed(ctx, limit, offset)
	if err != nil {
		return types.ArticleListResponseWrapper{}, err
	}
	return types.ArticleListResponseWrapper{
		Articles:      results,
		ArticlesCount: totalCount,
	}, nil
}

func (ac articlesController) get(ctx context.Context, _ goTypes.Nil) (types.ArticleWrapper[types.Article], error) {
	var fallbackResult types.ArticleWrapper[types.Article]
	id, err := api.PathVariable[string](ctx, "slug")
	article, err := ac.articleService.Get(ctx, id)
	if err != nil {
		return fallbackResult, err
	}
	fallbackResult.Article = article
	return fallbackResult, nil
}

func (ac articlesController) create(ctx context.Context, req types.ArticleWrapper[types.ArticleRequest]) (types.ArticleWrapper[types.Article], error) {
	var fallbackResult types.ArticleWrapper[types.Article]
	email, err := api.GetValue[string](ctx, "email")
	if err != nil {
		return fallbackResult, err
	}
	saved, err := ac.articleService.Create(ctx, req.Article, email)
	if err != nil {
		return fallbackResult, err
	}
	fallbackResult.Article = saved
	return fallbackResult, nil
}

func (ac articlesController) update(ctx context.Context, req types.ArticleWrapper[types.ArticleRequest]) (types.ArticleWrapper[types.Article], error) {
	var fallbackResult types.ArticleWrapper[types.Article]
	slug, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return fallbackResult, err
	}
	updated, err := ac.articleService.Update(ctx, slug, req.Article)
	if err != nil {
		return fallbackResult, err
	}
	fallbackResult.Article = updated
	return fallbackResult, nil
}

func (ac articlesController) delete(ctx context.Context, _ goTypes.Nil) (goTypes.Nil, error) {
	slug, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return goTypes.Nil{}, err
	}
	if err := ac.articleService.Delete(ctx, slug); err != nil {
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
	comment, err := ac.articleService.CreateComment(ctx, slug, req.Comment)
	if err != nil {
		return fallbackResult, err
	}
	fallbackResult.Comment = comment
	return fallbackResult, nil
}

func (ac articlesController) getComments(ctx context.Context, _ goTypes.Nil) (types.CommentListResponseWrapper, error) {
	slug, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return types.CommentListResponseWrapper{}, err
	}
	comments, err := ac.articleService.GetComments(ctx, slug)
	if err != nil {
		return types.CommentListResponseWrapper{}, err
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
	return goTypes.Nil{}, ac.articleService.DeleteComment(ctx, slug, id)
}

func (ac articlesController) addFavoriteArticle(ctx context.Context, _ goTypes.Nil) (types.ArticleWrapper[types.Article], error) {
	var fallbackResult types.ArticleWrapper[types.Article]
	slug, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return fallbackResult, err
	}
	email, err := api.GetValue[string](ctx, "email")
	if err != nil {
		return fallbackResult, err
	}
	article, err := ac.articleService.AddFavoriteArticle(ctx, slug, email)
	if err != nil {
		return fallbackResult, err
	}
	fallbackResult.Article = article
	return fallbackResult, nil
}

func (ac articlesController) deleteFavoriteArticle(ctx context.Context, _ goTypes.Nil) (types.ArticleWrapper[types.Article], error) {
	var fallbackResult types.ArticleWrapper[types.Article]
	slug, err := api.PathVariable[string](ctx, "slug")
	if err != nil {
		return fallbackResult, err
	}
	email, err := api.GetValue[string](ctx, "email")
	if err != nil {
		return fallbackResult, err
	}
	article, err := ac.articleService.DeleteFavoriteArticle(ctx, slug, email)
	if err != nil {
		return fallbackResult, err
	}
	fallbackResult.Article = article
	return fallbackResult, nil
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
