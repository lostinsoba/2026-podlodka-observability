package query

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"store/internal/model"
)

const (
	defaultLimit = 10
)

type MessagesByCursorQuery struct {
	Continue string
	Limit    int
}

func (q *MessagesByCursorQuery) ToModel(tenantID string) model.MessagesByCursorQueryRequest {
	return model.MessagesByCursorQueryRequest{
		TenantID: tenantID,
		Continue: q.Continue,
		Limit:    q.Limit,
	}
}

func ParseMessagesByCursorQuery(request *http.Request) (*MessagesByCursorQuery, error) {
	urlQuery, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		return nil, err
	}
	continueStr := urlQuery.Get("continue")
	var limit int
	limitStr := urlQuery.Get("limit")
	if limitStr == "" {
		limit = defaultLimit
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse limit: %s", err)
		}
	}
	return &MessagesByCursorQuery{
		Continue: continueStr,
		Limit:    limit,
	}, nil
}

const (
	defaultPage     = 1
	defaultPageSize = 10
)

type MessagesByOffsetQuery struct {
	Page     int
	PageSize int
}

func (q *MessagesByOffsetQuery) ToModel(tenantID string) model.MessageByOffsetQueryRequest {
	return model.MessageByOffsetQueryRequest{
		TenantID: tenantID,
		Page:     q.Page,
		PageSize: q.PageSize,
	}
}

func ParseMessagesByOffsetQuery(request *http.Request) (*MessagesByOffsetQuery, error) {
	urlQuery, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		return nil, err
	}
	var (
		page     int
		pageSize int
	)
	pageStr := urlQuery.Get("page")
	if pageStr == "" {
		page = defaultPage
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse page: %s", err)
		}
	}
	pageSizeStr := urlQuery.Get("page_size")
	if pageSizeStr == "" {
		pageSize = defaultPageSize
	} else {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse page_size: %s", err)
		}
	}
	return &MessagesByOffsetQuery{
		Page:     page,
		PageSize: pageSize,
	}, nil
}
