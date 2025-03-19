package server

import (
	"net/http"
	"strconv"

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
