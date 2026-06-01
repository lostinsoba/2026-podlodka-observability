package dto

import (
	"encoding/json"
	"net/http"

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

func (md *MessagesData) ToModel() []model.Message {
	list := make([]model.Message, 0, len(md.Messages))
	for _, message := range md.Messages {
		list = append(list, message.ToModel())
	}
	return list
}

func (mdi *MessageData) ToModel() model.Message {
	return model.Message{
		TenantID:   mdi.TenantID,
		ExternalID: mdi.ID,
		Message:    mdi.Message,
		Timestamp:  mdi.Timestamp,
	}
}

const (
	maxRequestBodySize = 1048576
)

func ParseMessagesData(writer http.ResponseWriter, request *http.Request) (*MessagesData, error) {
	var messagesData *MessagesData
	request.Body = http.MaxBytesReader(writer, request.Body, maxRequestBodySize)
	dec := json.NewDecoder(request.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&messagesData)
	if err != nil {
		return nil, err
	}
	return messagesData, nil
}
