package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/query"
)

const (
	defaultLimit uint64 = 10
	defaultPage  uint64 = 1
)

func ParsePagination(r *http.Request) query.PaginationInput {
	pagination := query.PaginationInput{}
	pagination.Limit = parseLimit(r)
	pagination.Page = parsePage(r)
	return pagination
}

func ParseContentFilters(r *http.Request) content.ContentListFilters {
	return content.ContentListFilters{
		Watched:   parseBool("watched", r),
		Category:  parseString("category", r),
		Genres:    parseGenres(r),
		Name:      parseString("name", r),
		Summary:   parseString("summary", r),
		WishLevel: parseInt("wishLevel", r),
	}
}

func parseLimit(r *http.Request) uint64 {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		return defaultLimit
	}

	intLimit, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		return defaultLimit
	}

	return intLimit
}

func parsePage(r *http.Request) uint64 {
	limit := r.URL.Query().Get("page")
	if limit == "" {
		return defaultPage
	}

	intPage, err := strconv.ParseUint(limit, 10, 64)
	if err != nil {
		return defaultPage
	}

	return intPage
}

func parseBool(param string, r *http.Request) *bool {
	value := r.URL.Query().Get(param)
	if value == "" {
		return nil
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return nil
	}

	return &b
}

func parseInt(param string, r *http.Request) *int {
	value := r.URL.Query().Get(param)
	if value == "" {
		return nil
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}

	return &i
}

func parseString(param string, r *http.Request) *string {
	value := r.URL.Query().Get(param)
	if value == "" {
		return nil
	}
	return &value
}

func parseGenres(r *http.Request) *[]string {
	genres := r.URL.Query().Get("genres")
	if genres == "" {
		return nil
	}
	genresSlice := strings.Split(genres, ",")
	return &genresSlice
}
