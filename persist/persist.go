package persist

import (
	"context"

	"github.com/borosr/realworld/persist/badger"
	"github.com/borosr/realworld/persist/types"
)

type Repository[Type types.Storable] interface {
	Save(ctx context.Context, data Type) (Type, error)
	Get(ctx context.Context, key string) (Type, error)
	GetFiltered(ctx context.Context, filters ...types.Filter[Type]) ([]Type, error)
	CountFiltered(ctx context.Context, filters ...types.Filter[Type]) (uint64, error)
	Delete(ctx context.Context, key string) error
	Sequence(ctx context.Context, key string) (uint64, error)
}

func Get[Type types.Storable]() Repository[Type] {
	return badger.Get[Type]()
}
