package mid

import (
	"context"
	"fmt"
	"github/islamghany/blog/business/web/v1/response"
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/validate"
	"github/islamghany/blog/foundation/web"
	"net/http"
)

func Errors(log *logger.Logger) web.Middleware {
	return func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			if err := handler(ctx, w, r); err != nil {
				log.Error(ctx, "message", "msg", err)
				var er response.ErrorDocument
				var status int
				switch {
				case response.IsError(err):
					respErr := response.GetError(err)
					// Check if the error is a validation error.
					if validate.IsFieldErrors(respErr.Err) {
						fieldsErros := validate.GetFieldsErrors(respErr.Err)
						er = response.ErrorDocument{
							Error:  "data validation error",
							Fields: fieldsErros.Fields(),
						}
						status = http.StatusBadRequest
						break
					}
					er = response.ErrorDocument{
						Error: respErr.Error(),
					}
					status = respErr.Status
					break
				default:
					er = response.ErrorDocument{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}
				if err = web.Response(ctx, w, er, status); err != nil {
					return fmt.Errorf("sending error: %w", err)
				}
			}
			return nil

		}
		return h
	}
}
