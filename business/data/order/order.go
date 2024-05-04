package order

import (
	"github/islamghany/blog/foundation/validate"
	"net/http"
	"strings"
)

const (
	DESC = "DESC"
	ASC  = "ASC"
)

var directions = map[string]string{
	DESC: DESC,
	ASC:  ASC,
}

type By struct {
	Field string
	Dir   string
}

func NewBy(field, dir string) By {
	return By{
		Field: field,
		Dir:   dir,
	}
}

func Parse(r *http.Request, defaultOrder By) (By, error) {
	v := r.URL.Query().Get("order")

	if v == "" {
		return defaultOrder, nil
	}

	orderParts := strings.Split(v, ",")
	var by By
	switch len(orderParts) {
	case 1:
		by = NewBy(orderParts[0], ASC)
	case 2:
		by = NewBy(orderParts[0], orderParts[1])
	default:
		return By{}, validate.NewFieldError("order", "must be in the format field,dir")
	}
	if _, exits := directions[by.Dir]; !exits {
		return By{}, validate.NewFieldError("order", "dir must be DESC or ASC")
	}
	return by, nil

}
