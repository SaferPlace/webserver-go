package middleware

import "net/http"

// Middleware transforms the request or response in some way.
type Middleware func(http.Handler) http.Handler