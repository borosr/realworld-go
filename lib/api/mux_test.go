package api

import (
	"context"
	goTypes "go/types"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestBuildParams(t *testing.T) {
	t.Run("buildParams", func(t *testing.T) {
		t.Parallel()
		params := buildParams([]string{"api", "users", "{username}"})
		t.Log(params)
		assert.Equal(t, 1, len(params))
		param, ok := params[2]
		assert.True(t, ok, "params has 2 as key")
		assert.Equal(t, "username", param.name)
		assert.Nil(t, param.pattern)
	})
	t.Run("buildParams_split", func(t *testing.T) {
		t.Parallel()
		params := buildParams(strings.Split("/api/users/{username}", "/"))
		t.Log(params)
		assert.Equal(t, 1, len(params))
		param, ok := params[3]
		assert.True(t, ok, "params has 3 as key")
		assert.Equal(t, "username", param.name)
		assert.Nil(t, param.pattern)
	})
	t.Run("buildParams_multiple_params", func(t *testing.T) {
		t.Parallel()
		params := buildParams([]string{"api", "users", "{username}", "other_thing", "{thing:[a-zA-Z0-9]*}"})
		t.Log(params)
		assert.Equal(t, 2, len(params))
		paramUn, ok := params[2]
		assert.True(t, ok, "params has 2 as key")
		assert.Equal(t, "username", paramUn.name)
		assert.Nil(t, paramUn.pattern)
		paramTh, ok := params[4]
		assert.True(t, ok, "params has 4 as key")
		assert.Equal(t, "thing", paramTh.name)
		assert.NotNil(t, paramTh.pattern)
	})
	t.Run("buildParams_multiple_params_split", func(t *testing.T) {
		t.Parallel()
		params := buildParams(strings.Split("/api/users/{username}/other_thing/{thing:[a-zA-Z0-9]*}", "/"))
		t.Log(params)
		assert.Equal(t, 2, len(params))
		paramUn, ok := params[3]
		assert.True(t, ok, "params has 2 as key")
		assert.Equal(t, "username", paramUn.name)
		assert.Nil(t, paramUn.pattern)
		paramTh, ok := params[5]
		assert.True(t, ok, "params has 4 as key")
		assert.Equal(t, "thing", paramTh.name)
		assert.NotNil(t, paramTh.pattern)
	})
}

func BenchmarkRegisterMultiplePathVariablesRegexp(b *testing.B) {
	// BenchmarkRegisterMultiplePathVariablesRegexp-10    	  504189	      2025 ns/op	    1873 B/op	      23 allocs/op
	type test struct{}
	b.ReportAllocs()
	handler := func(ctx context.Context, _ test) (test, error) {
		return test{}, nil
	}
	for i := 0; i < b.N; i++ {
		Register[test, test, ControllerSimpleFunc[test, test]]("/test/path/{asd:[ads]+}/thing/{other}/"+xid.New().String(), http.MethodPost, handler)
	}
}

func BenchmarkRegisterOnePathVariableNoRegex(b *testing.B) {
	// BenchmarkRegisterOnePathVariableNoRegex-10    	 1272385	       908.2 ns/op	     815 B/op	       9 allocs/op
	type test struct{}
	b.ReportAllocs()
	handler := func(ctx context.Context, _ test) (test, error) {
		return test{}, nil
	}
	for i := 0; i < b.N; i++ {
		Register[test, test, ControllerSimpleFunc[test, test]]("/test/path/thing/{other}/"+xid.New().String(), http.MethodPost, handler)
	}
}

func BenchmarkRegisterOneNoPathVariable(b *testing.B) {
	// BenchmarkRegisterOneNoPathVariable-10    	 2088630	       527.9 ns/op	     377 B/op	       6 allocs/op
	type test struct{}
	b.ReportAllocs()
	handler := func(ctx context.Context, _ test) (test, error) {
		return test{}, nil
	}
	for i := 0; i < b.N; i++ {
		Register[test, test, ControllerSimpleFunc[test, test]]("/test/path/thing/"+xid.New().String(), http.MethodPost, handler)
	}
}

func BenchmarkHandleFunc(b *testing.B) {
	s := Server{}
	type test struct{}
	handlerFunc := func(ctx context.Context, _ test) (test, error) {
		return test{}, nil
	}
	handler := methodWrapper[test, test, ControllerSimpleFunc[test, test]]("/test/path/"+xid.New().String(), http.MethodPost, handlerFunc)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.HandleFunc("/test/path/"+xid.New().String(), http.MethodPost, handler)
	}
}

func BenchmarkBuildParams(b *testing.B) {
	split := strings.Split("/test/path", "/")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buildParams(split)
	}
}

func BenchmarkBuildParamsRegex(b *testing.B) {
	split := strings.Split("/test/path/{value:.*}", "/")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buildParams(split)
	}
}

func BenchmarkBuildParamsPathVariable(b *testing.B) {
	split := strings.Split("/test/path/{value}", "/")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buildParams(split)
	}
}

func BenchmarkServer_ServeHTTP_Post(b *testing.B) {
	// BenchmarkServer_ServeHTTP-10    	  404146	      2626 ns/op	    3151 B/op	      45 allocs/op
	s := Server{}
	type test struct{}
	handlerFunc := func(ctx context.Context, _ test) (test, error) {
		return test{}, nil
	}
	var path string
	for i := 0; i < 20; i++ {
		path = "/test/path/" + xid.New().String()
		s.HandleFunc(path, http.MethodPost, methodWrapper[test, test, ControllerSimpleFunc[test, test]](path, http.MethodPost, handlerFunc))
	}
	req, _ := http.NewRequest(http.MethodPost, path, nil)
	w := httptest.NewRecorder()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.ServeHTTP(w, req)
	}
}

func BenchmarkServer_ServeHTTP_Get(b *testing.B) {
	// BenchmarkServer_ServeHTTP-10    	  404146	      2626 ns/op	    3151 B/op	      45 allocs/op
	s := Server{}
	type test struct{}
	handlerFunc := func(ctx context.Context, _ goTypes.Nil) (test, error) {
		return test{}, nil
	}
	var path string
	for i := 0; i < 20; i++ {
		path = "/test/path/" + xid.New().String()
		s.HandleFunc(path, http.MethodGet, methodWrapper[goTypes.Nil, test, ControllerSimpleFunc[goTypes.Nil, test]](path, http.MethodGet, handlerFunc))
	}
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.ServeHTTP(w, req)
	}
}

func TestServer_parse(t *testing.T) {
	s := Server{}
	type test struct{}
	handlerFunc := func(ctx context.Context, _ test) (test, error) {
		return test{}, nil
	}
	path := "/test/path/" + xid.New().String()
	handler := methodWrapper[test, test, ControllerSimpleFunc[test, test]](path, http.MethodPost, handlerFunc)
	s.HandleFunc(path, http.MethodPost, handler)
	sp := strings.Split(path, "/")
	r := route{pattern: sp, method: http.MethodPost, params: buildParams(sp), handler: http.HandlerFunc(handler)}
	r.parse(context.Background(), path, http.MethodPost)
}
