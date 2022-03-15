package api

import (
	"context"
	goTypes "go/types"
	"net/http"
	"sort"

	"github.com/borosr/realworld/lib/api"
	"github.com/borosr/realworld/persist"
	"github.com/borosr/realworld/types"
)

type tagsController struct {
	articleRepository persist.Repository[*types.Article]
}

func (tc tagsController) Init() {
	api.Register[
		goTypes.Nil,
		types.TagsWrapper,
		api.ControllerSimpleFunc[goTypes.Nil, types.TagsWrapper],
	]("/api/tags", http.MethodGet, tc.getAll)
}

func (tc tagsController) getAll(ctx context.Context, _ goTypes.Nil) (types.TagsWrapper, error) {
	articles, err := tc.articleRepository.GetFiltered(ctx)
	if err != nil {
		return types.TagsWrapper{}, err
	}
	var tagSet = make(map[string]struct{})
	for i := range articles {
		for _, tag := range articles[i].TagList {
			tagSet[tag] = struct{}{}
		}
	}
	var tagList = make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tagList = append(tagList, tag)
	}
	sort.Strings(tagList)
	return types.TagsWrapper{
		Tags: tagList,
	}, nil
}
