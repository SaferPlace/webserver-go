package middleware

import (
	"net/http"

	"github.com/rs/cors"
	"golang.org/x/exp/slices"
)

// Cors allows to only allow known domains.
func Cors(domains []string) Middleware {
	return cors.New(cors.Options{
		AllowedMethods: []string{
			// CORS preflight
			http.MethodOptions,
			// Metrics
			http.MethodGet,
			// connect RPCs
			http.MethodPost,
		},
		// Mirror the `Origin` header value in the `Access-Control-Allow-Origin`
		// preflight response header.
		// This is equivalent to `Access-Control-Allow-Origin: *`, but allows
		// for requests with credentials.
		// Note that this effectively disables CORS and is not safe for use in
		// production environments.
		AllowOriginFunc: func(origin string) bool {
			// Disable CORS when the domain list is not specified.
			// This might be a security issue long term.
			if len(domains) == 0 {
				return true
			}
			return slices.Contains(domains, origin)
		},
		// Note that rs/cors does not return `Access-Control-Allow-Headers: *`
		// in response to preflight requests with the following configuration.
		// It simply mirrors all headers listed in the `Access-Control-Request-Headers`
		// preflight request header.
		AllowedHeaders: []string{"*"},
		// We explicitly set the exposed header names instead of using the wildcard *,
		// because in requests with credentials, it is treated as the literal header
		// name "*" without special semantics.
		ExposedHeaders: []string{
			"Grpc-Status", "Grpc-Message", "Grpc-Status-Details-Bin", "X-Grpc-Test-Echo-Initial",
			"Trailer-X-Grpc-Test-Echo-Trailing-Bin"},
	}).Handler
}
