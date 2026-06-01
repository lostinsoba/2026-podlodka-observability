package dto

import (
	"store/internal/model"
)

type MessagesData struct {
	Messages []MessageData `json:"messages"`
}

type MessageData struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func ToMessageData(message model.Message) MessageData {
	return MessageData{
		ID:        message.ID,
		TenantID:  message.TenantID,
		Message:   message.Message,
		Timestamp: message.Timestamp,
	}
}

type MessagesByCursorQueryResponseData struct {
	Messages []MessageData `json:"messages"`
	Continue string        `json:"continue"`
}

func ToMessagesByCursorQueryResponseData(response model.MessagesByCursorQueryResponse) MessagesByCursorQueryResponseData {
	messages := make([]MessageData, 0, len(response.Messages))
	for _, message := range response.Messages {
		messages = append(messages, ToMessageData(message))
	}
	return MessagesByCursorQueryResponseData{
		Messages: messages,
		Continue: response.Continue,
	}
}

type MessagesByPageResponseData struct {
	Messages   []MessageData `json:"messages"`
	TotalPages int           `json:"total_pages"`
	TotalCount int           `json:"total_count"`
}

func ToMessagesByOffsetQueryResponseData(response model.MessageByOffsetQueryResponse) MessagesByPageResponseData {
	messages := make([]MessageData, 0, len(response.Messages))
	for _, message := range response.Messages {
		messages = append(messages, ToMessageData(message))
	}
	return MessagesByPageResponseData{
		Messages:   messages,
		TotalPages: response.TotalPages,
		TotalCount: response.TotalCount,
	}
}
