package mid

import (
	"net/http"
)

func CORS(whitelist []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// telling the client that the server will accept requests from any origin
			// and not to cache the response
			w.Header().Add("Vary", "Origin")
			// Add the "Vary: Access-Control-Request-Method" header.
			w.Header().Add("Vary", "Access-Control-Request-Method")
			// Set the CORS headers
			origin := r.Header.Get("Origin")

			for _, o := range whitelist {
				if o == origin {
					w.Header().Add("Access-Control-Allow-Origin", o)
					// Check if the request has the HTTP method OPTIONS and contains the
					// "Access-Control-Request-Method" header. If it does, then we treat
					// it as a preflight request.
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						// Set the necessary preflight response header
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
						// Write the headers along with a 200 OK status and return from
						// the middleware with no further action.
						w.WriteHeader(http.StatusOK)
						return
					}
					break
				}
			}

			next.ServeHTTP(w, r)
		})

	}
}
