package paging

import (
	"github/islamghany/blog/foundation/validate"
	"net/http"
	"strconv"
)

type Response[T any] struct {
	Items    []T `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func NewResponse[T any](items []T, total, page, pageSize int) Response[T] {
	return Response[T]{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

type Page struct {
	Number int
	Size   int
}

func ParseRequest(r *http.Request) (Page, error) {
	values := r.URL.Query()
	number := 1
	size := 10
	if v := values.Get("page"); v != "" {
		var err error
		number, err = strconv.Atoi(v)
		if err != nil {
			return Page{}, validate.NewFieldError("page", "must be an integer")
		}
	}
	if v := values.Get("page_size"); v != "" {
		var err error
		size, err = strconv.Atoi(v)
		if err != nil {
			return Page{}, validate.NewFieldError("page_size", "must be an integer")
		}
	}
	return Page{Number: number, Size: size}, nil
}
