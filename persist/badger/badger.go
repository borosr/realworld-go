package badger

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"

	"github.com/borosr/realworld/persist/types"
	bdb "github.com/dgraph-io/badger/v3"
	"github.com/rs/xid"
)

const defaultSequenceBandwidth = 1
const sequencePrefix = "seq"

var db *bdb.DB
var mutex sync.Mutex

type Repository[Type types.Storable] struct {
	db *bdb.DB
}

func Get[Type types.Storable]() Repository[Type] {
	getDB()
	return Repository[Type]{
		db: db,
	}
}

func getDB() {
	if db == nil {
		mutex.Lock()
		defer mutex.Unlock()
		var err error
		db, err = bdb.Open(bdb.DefaultOptions("/tmp/badger"))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (r Repository[Type]) Save(_ context.Context, data Type) (Type, error) {
	if data.Key() == "" {
		data.SetKey(xid.New().String())
	}
	if err := r.db.Update(func(txn *bdb.Txn) error {
		rawData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		key := r.buildID(data.Name(), data.Key())
		return txn.Set([]byte(key), rawData)
	}); err != nil {
		return data, err
	}
	return data, nil
}

func (r Repository[Type]) Get(_ context.Context, key string) (Type, error) {
	var t Type
	if err := r.db.View(func(txn *bdb.Txn) error {
		item, err := txn.Get([]byte(r.buildID(t.Name(), key)))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &t)
		})
	}); err != nil {
		return t, err
	}
	return t, nil
}

func (r Repository[Type]) GetFiltered(_ context.Context, filters ...types.Filter[Type]) ([]Type, error) {
	var res = make([]Type, 0)
	if err := r.db.View(func(txn *bdb.Txn) error {
		options := bdb.DefaultIteratorOptions
		var t Type
		options.Prefix = []byte(t.Name())
		it := txn.NewIterator(options)
		defer it.Close()
	outer:
		for it.Rewind(); it.Valid(); it.Next() {
			if err := it.Item().Value(func(val []byte) error {
				return json.Unmarshal(val, &t)
			}); err != nil {
				//return err
				continue
			}
			for _, filter := range filters {
				if ok := filter(t); !ok {
					continue outer
				}
			}
			res = append(res, t)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (r Repository[Type]) CountFiltered(_ context.Context, filters ...types.Filter[Type]) (uint64, error) {
	var count uint64
	if err := r.db.View(func(txn *bdb.Txn) error {
		options := bdb.DefaultIteratorOptions
		var t Type
		options.Prefix = []byte(t.Name())
		it := txn.NewIterator(options)
		defer it.Close()
	outer:
		for it.Rewind(); it.Valid(); it.Next() {
			if err := it.Item().Value(func(val []byte) error {
				return json.Unmarshal(val, &t)
			}); err != nil {
				return err
			}
			for _, filter := range filters {
				if ok := filter(t); !ok {
					continue outer
				}
			}
			count++
		}
		return nil
	}); err != nil {
		return 0, err
	}
	return count, nil
}

func (r Repository[Type]) Delete(_ context.Context, key string) error {
	return r.db.Update(func(txn *bdb.Txn) error {
		var t Type
		return txn.Delete([]byte(r.buildID(t.Name(), key)))
	})
}

func (r Repository[Type]) Sequence(_ context.Context, key string) (uint64, error) {
	var t Type
	seq, err := r.db.GetSequence([]byte(r.buildID(sequencePrefix, t.Name(), key)), defaultSequenceBandwidth)
	if err != nil {
		return 0, err
	}
	next, err := seq.Next()
	if err != nil {
		return 0, err
	}
	return next, nil
}

func (r Repository[Type]) buildID(parts ...string) string {
	return strings.Join(parts, "-")
}
