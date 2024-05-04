package web

import (
	"context"
	"encoding/json"
	"net/http"
)

func Response(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	// set the statusCode
	SetStatusCode(ctx, statusCode)
	w.WriteHeader(statusCode)
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil

}
