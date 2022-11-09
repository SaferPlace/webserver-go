// Copyright 2022 SaferPlace

package webserver

import "net/http"

// Service registers the service with the server, it
// returns the path on which the server needs to be
// registered and the handler for that path.
type Service func() (string, http.Handler)
