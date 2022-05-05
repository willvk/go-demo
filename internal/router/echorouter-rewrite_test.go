package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Verify stax regex used with rewrite
// Special attn for query params which are tested via separate secondary assertion
// see: https://github.com/labstack/echo/issues/1798 and PR: https://github.com/labstack/echo/pull/1802/files
func TestStaxURLRewrite(t *testing.T) {
	e := New()

	var rec *httptest.ResponseRecorder
	var req *http.Request

	testCases := []struct {
		requestPath string
		expectPath  string
		query       string
	}{
		{"/unmatched", "/unmatched", ""},
		{"/orders/1234/items", "/orders/1234/items", ""},
		{"/90210/test", "/test", ""},
		{"/10222010/foo/c/bar/baz", "/foo/c/bar/baz", ""},
		{"/c/ignore/test", "/c/ignore/test", ""},
		{"/1/ignore/test", "/ignore/test", ""},
		{"/112/sometext/sdad?", "/sometext/sdad", ""},
		{"/abc123/sometext/sdad?", "/abc123/sometext/sdad", ""},
		{"/20200627/groups?organizationIDs=de653482-72f5-4051-bfca-fcbe58370261&status=ACTIVE&userIDs=09832b50-bc04-4769-b4d2-e1da7888dbc2",
			"/groups", "organizationIDs=de653482-72f5-4051-bfca-fcbe58370261&status=ACTIVE&userIDs=09832b50-bc04-4769-b4d2-e1da7888dbc2"},
	}

	for _, tc := range testCases {
		requestPath, expectPath, query := tc.requestPath, tc.expectPath, tc.query
		t.Run(requestPath, func(t *testing.T) {
			req = httptest.NewRequest(http.MethodGet, requestPath, nil)
			rec = httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, expectPath, req.URL.EscapedPath())
			assert.Equal(t, query, req.URL.RawQuery)
		})
	}
}
