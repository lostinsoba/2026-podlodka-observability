package interceptor

import "net/http"

type Doer func(*http.Request) (*http.Response, error)

type Interceptor func(next Doer) Doer
