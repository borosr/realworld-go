package types

type Storable interface {
	Name() string
	Key() string
	SetKey(id string)
}

type Filter[Type Storable] func(t Type) bool
