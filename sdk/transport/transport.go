package transport

import "net/http"

type Wrapper func(internalRoundTripper http.RoundTripper) http.RoundTripper
