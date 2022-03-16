package api

import (
	"log"

	"github.com/borosr/realworld/domain"
	"github.com/borosr/realworld/lib/api"
	"github.com/borosr/realworld/persist"
	"github.com/borosr/realworld/types"
)

func Service() {
	log.Println("Listening on 18000...")

	initControllers()

	if err := api.ListenAndServe(":18000"); err != nil {
		log.Fatal(err)
	}
}

func initControllers() {
	userRepository := persist.Get[*types.User]()
	articleRepository := persist.Get[*types.Article]()
	followRepository := persist.Get[*types.Follow]()
	commentRepository := persist.Get[*types.Comment]()
	favoriteRepository := persist.Get[*types.Favorite]()

	userService := domain.UserService{
		UserRepository: userRepository,
	}
	profileService := domain.ProfileService{
		UserRepository:   userRepository,
		FollowRepository: followRepository,
	}
	articleService := domain.ArticleService{
		ArticleRepository:  articleRepository,
		CommentRepository:  commentRepository,
		FavoriteRepository: favoriteRepository,
		UserService:        userService,
	}

	userController{
		userService: userService,
	}.Init()
	profilesController{
		profileService: profileService,
		userService:    userService,
	}.Init()
	articlesController{
		articleService: articleService,
	}.Init()
	tagsController{
		articleRepository: articleRepository,
	}.Init()
}
