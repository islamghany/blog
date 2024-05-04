package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

const (
	MAX_BODY_SIZE = int64(1024 * 1024 * 2) // 2MB
)

// Decode is a helper function to decode JSON data from an HTTP request and runs validation on it.
type validator interface {
	Validate() error
}

func Decode(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_BODY_SIZE)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("decoding request: %w", err)
	}
	if d, ok := dst.(validator); ok {
		if err := d.Validate(); err != nil {
			return fmt.Errorf("validating request: %w", err)
		}
	}
	return nil
}

func ParamID(r *http.Request, key string) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())
	val := params.ByName(key)
	if val == "" {
		return "", fmt.Errorf("url parameter %q missing", key)
	}
	return val, nil
}

func ParamUUID(r *http.Request, key string) (uuid.UUID, error) {
	params := httprouter.ParamsFromContext(r.Context())
	val := params.ByName(key)
	if val == "" {
		return uuid.UUID{}, fmt.Errorf("url parameter %q missing", key)
	}
	uid, err := uuid.Parse(val)
	if err != nil {
		return uuid.UUID{}, errors.New("url parameter id is not a valid uuid")
	}
	return uid, nil
}

func ReadString(q url.Values, key string, defaultValue string) string {
	if q.Get(key) == "" {
		return defaultValue
	}
	return q.Get(key)
}

func ReadInt(q url.Values, key string, defaultValue int) int {
	s := q.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return i
}

func ReadBool(q url.Values, key string, defaultValue bool) bool {
	s := q.Get(key)
	if s == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return defaultValue
	}
	return b
}

func ReadCSV(q url.Values, key string, defaultValue []string) []string {
	s := q.Get(key)
	if s == "" {
		return defaultValue
	}
	return strings.Split(s, ",")

}
