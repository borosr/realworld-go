package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/rs/xid"
)

func BenchmarkMethodWrapper(b *testing.B) {
	// BenchmarkMethodWrapper-10    	11476627	       101.0 ns/op	      96 B/op	       2 allocs/op
	type test struct{}
	b.ReportAllocs()
	handler := func(ctx context.Context, _ test) (test, error) {
		return test{}, nil
	}
	for i := 0; i < b.N; i++ {
		methodWrapper[test, test, ControllerSimpleFunc[test, test]]("/test/path/"+xid.New().String(), http.MethodPost, handler)
	}
}
