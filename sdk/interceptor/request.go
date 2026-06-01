package interceptor

import (
	"net/http"
	"net/url"
)

func WithRequestHeader(key, value string) Interceptor {
	return func(next Doer) Doer {
		return func(request *http.Request) (*http.Response, error) {
			request.Header.Set(key, value)
			return next(request)
		}
	}
}

func WithRequestQueryParams(params url.Values) Interceptor {
	return func(next Doer) Doer {
		return func(request *http.Request) (*http.Response, error) {
			request.URL.RawQuery = params.Encode()
			return next(request)
		}
	}
}
