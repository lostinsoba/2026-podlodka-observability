package model

import "fmt"

type Message struct {
	ID         string
	ExternalID string
	TenantID   string
	Message    string
	Timestamp  int64
}

func (m Message) GetUniqueKey() string {
	return fmt.Sprintf("%s-%s", m.TenantID, m.ExternalID)
}

type UpdateMessagesRequest struct {
	Messages map[string]Message
}

type UpdateMessagesResponse struct {
	MessagesUpdated int64
}

type MessagesByCursorQueryRequest struct {
	TenantID string
	Continue string
	Limit    int
}

func (mqr MessagesByCursorQueryRequest) IsFirstOne() bool {
	return mqr.Continue == ""
}

type MessagesByCursorQueryResponse struct {
	Messages []Message
	Continue string
}

type MessagesStats struct {
	MessageCount int64
}

type MessageByOffsetQueryRequest struct {
	TenantID string
	Page     int
	PageSize int
}

type MessageByOffsetQueryResponse struct {
	Messages   []Message
	TotalPages int
	TotalCount int
}
