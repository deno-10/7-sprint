package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	city := "moscow"
	totalCount := len(cafeList[city])

	wantFor100 := 100
	if totalCount < 100 {
		wantFor100 = totalCount
	}

	requests := []struct {
		count int
		want  int
	}{
		{count: 0, want: 0},
		{count: 1, want: 1},
		{count: 2, want: 2},
		{count: 100, want: wantFor100},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city="+city+"&count="+strconv.Itoa(v.count), nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		body := strings.TrimSpace(response.Body.String())

		got := 0
		if body != "" {
			got = len(strings.Split(body, ","))
		}

		assert.Equal(t, v.want, got, "count=%d", v.count)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{search: "фасоль", wantCount: 0},
		{search: "кофе", wantCount: 2},
		{search: "вилка", wantCount: 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city=moscow&search="+v.search, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		body := strings.TrimSpace(response.Body.String())

		var cafes []string
		if body != "" {
			cafes = strings.Split(body, ",")
		}

		assert.Len(t, cafes, v.wantCount, "search=%s", v.search)

		for _, cafe := range cafes {
			assert.True(t,
				strings.Contains(strings.ToLower(cafe), strings.ToLower(v.search)),
				"кафе %q не содержит подстроку %q", cafe, v.search)
		}
	}
}
