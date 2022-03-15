package api

import (
	"log"

	"github.com/borosr/realworld/lib/api"
	"github.com/borosr/realworld/persist"
	"github.com/borosr/realworld/types"
)

func Service() {
	log.Println("Listening on 18000...")

	userRepository := persist.Get[*types.User]()
	articleRepository := persist.Get[*types.Article]()

	userController{
		userRepository: userRepository,
	}.Init()
	profilesController{
		userRepository: userRepository,
	}.Init()
	articlesController{
		articleRepository:  articleRepository,
		commentRepository:  persist.Get[*types.Comment](),
		favoriteRepository: persist.Get[*types.Favorite](),
		userRepository:     userRepository,
	}.Init()
	tagsController{
		articleRepository: articleRepository,
	}.Init()

	if err := api.ListenAndServe(":18000"); err != nil {
		log.Fatal(err)
	}
}
