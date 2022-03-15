package api

import (
	"strings"
	"testing"
)

func TestBuildParams(t *testing.T) {
	t.Logf("%+v", buildParams([]string{"api", "users", "{username}"}))
	t.Logf("%+v", buildParams(strings.Split("/api/users/{username}", "/")))
	t.Logf("%+v", buildParams([]string{"api", "users", "{username}", "other_thing", "{thing:[a-zA-Z0-9]*}"}))
	t.Logf("%+v", buildParams(strings.Split("/api/users/{username}/other_thing/{thing:[a-zA-Z0-9]*}", "/")))
}
