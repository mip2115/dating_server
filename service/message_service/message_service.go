package messageservice

import (
	"code.mine/dating_server/types"
	uuid "github.com/satori/go.uuid"

	"errors"
	"time"

	"code.mine/dating_server/mapping"
	"code.mine/dating_server/repo"
)

// MessageController -
type MessageController struct {
	repo repo.Repo
}

func (c *MessageController) AddMessage(
	messageRequest *types.MessageRequest,
) (*types.Message, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	if messageRequest.MatchUUID == nil {
		return nil, errors.New("need match uuid in message request")
	}
	if messageRequest.Content == nil {
		return nil, errors.New("need content in message request")
	}
	if messageRequest.From == nil {
		return nil, errors.New("need from in message request")
	}
	if messageRequest.To == nil {
		return nil, errors.New("need to in message request")
	}
	msg := &types.Message{}
	msg.From = messageRequest.From
	msg.To = messageRequest.To
	msg.Content = messageRequest.Content
	msg.UUID = mapping.StrToPtr(u.String())
	t := time.Now()
	msg.DateCreated = &t
	msg.DateUpdated = &t
	err = c.repo.SaveMessage(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// GetMessages -
func (c *MessageController) GetMessages(messageRequest *types.MessageRequest, nPerPage int) ([]*types.Message, error) {
	if messageRequest.MatchUUID == nil {
		return nil, errors.New("need match uuid in message request")
	}
	if messageRequest.Page == nil {
		return nil, errors.New("need page in message request")
	}
	pagesToSkip := mapping.IntToV(messageRequest.Page)
	msgs, err := c.repo.GetMessagesByMatchUUID(pagesToSkip, nPerPage, messageRequest.MatchUUID)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
