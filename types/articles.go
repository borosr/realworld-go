package types

import (
	"strconv"
	"time"
)

type ArticleListResponseWrapper struct {
	Articles      []*Article `json:"articles"`
	ArticlesCount int        `json:"articlesCount"`
}

type ArticleWrapper[SpecificArticle Article | ArticleRequest] struct {
	Article SpecificArticle `json:"article"`
}

type Article struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	TagList        []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         Profile   `json:"author"`
}

func (a *Article) Name() string {
	return "article"
}

func (a *Article) Key() string {
	return a.Slug
}

func (a *Article) SetKey(id string) {
	a.Slug = id
}

type Favorite struct {
	Slug     string `json:"slug"`
	Username string `json:"username"`
}

func (f *Favorite) Name() string {
	return "favorite"
}

func (f *Favorite) Key() string {
	return f.Slug + "-" + f.Username
}

func (f *Favorite) SetKey(id string) {
	// DO NOTHING
}

type ArticleRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Body        string   `json:"body"`
	TagList     []string `json:"tagList"`
}

type CommentWrapper[SpecificComment CommonComment | CommentRequest] struct {
	Comment SpecificComment `json:"comment"`
}

type CommentListResponseWrapper struct {
	Comments []CommonComment `json:"comments"`
}

type CommentRequest struct {
	Body string `json:"body"`
}

type CommonComment struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
	Author    Profile   `json:"author"`
}

type Comment struct {
	CommonComment
	Slug string `json:"slug"`
}

func (c *Comment) Name() string {
	return "comment"
}

func (c *Comment) Key() string {
	return c.Slug + "-" + strconv.Itoa(c.ID)
}

func (c *Comment) SetKey(_ string) {
	// DO NOTHING
}
