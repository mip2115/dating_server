package messageservice

import (
	"testing"

	"code.mine/dating_server/mapping"

	mockRepo "code.mine/dating_server/mocks/repo"
	"code.mine/dating_server/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type MessageServiceTestSuite struct {
	suite.Suite
}

func TestMessageServiceTestSuite(t *testing.T) {
	suite.Run(t, new(MessageServiceTestSuite))
}

func (suite *MessageServiceTestSuite) TestAddMessage() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	messageController := MessageController{
		repo: mockRepo,
	}

	// t := time.Now()
	// savedMessage := &types.Message{
	// 	From:        mapping.StrToPtr("some-uuid"),
	// 	To:          mapping.StrToPtr("some-uuid"),
	// 	Content:     mapping.StrToPtr("This is a message"),
	// 	MatchUUID:   mapping.StrToPtr("match-uuid"),
	// 	UUID:        mapping.StrToPtr("message-uuid"),
	// 	DateCreated: &t,
	// 	DateUpdated: &t,
	// }

	mockRepo.EXPECT().SaveMessage(gomock.Any()).Return(nil)

	messageRequest := &types.MessageRequest{
		From:      mapping.StrToPtr("some-uuid"),
		To:        mapping.StrToPtr("some-uuid"),
		Content:   mapping.StrToPtr("This is a message"),
		MatchUUID: mapping.StrToPtr("match-uuid"),
		UUID:      mapping.StrToPtr("message-uuid"),
		Page:      mapping.IntToPtr(0),
	}
	msg, err := messageController.AddMessage(messageRequest)
	suite.Require().NoError(err)
	suite.Require().NotNil(msg)
}

func (suite *MessageServiceTestSuite) TestAddMessage_Error() {
	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	messageController := MessageController{
		repo: mockRepo,
	}

	messageRequestNoFrom := &types.MessageRequest{
		To:        mapping.StrToPtr("some-uuid"),
		Content:   mapping.StrToPtr("This is a message"),
		MatchUUID: mapping.StrToPtr("match-uuid"),
		UUID:      mapping.StrToPtr("message-uuid"),
		Page:      mapping.IntToPtr(0),
	}
	messageRequestNoTo := &types.MessageRequest{
		From:      mapping.StrToPtr("some-uuid"),
		Content:   mapping.StrToPtr("This is a message"),
		MatchUUID: mapping.StrToPtr("match-uuid"),
		UUID:      mapping.StrToPtr("message-uuid"),
		Page:      mapping.IntToPtr(0),
	}
	messageRequestNoContent := &types.MessageRequest{
		From:      mapping.StrToPtr("some-uuid"),
		To:        mapping.StrToPtr("some-uuid"),
		MatchUUID: mapping.StrToPtr("match-uuid"),
		UUID:      mapping.StrToPtr("message-uuid"),
		Page:      mapping.IntToPtr(0),
	}
	messageRequestNoMatchUUID := &types.MessageRequest{
		From:    mapping.StrToPtr("some-uuid"),
		To:      mapping.StrToPtr("some-uuid"),
		Content: mapping.StrToPtr("This is a message"),
		UUID:    mapping.StrToPtr("message-uuid"),
		Page:    mapping.IntToPtr(0),
	}

	tests := []struct {
		req *types.MessageRequest
		msg string
	}{
		{
			req: messageRequestNoFrom,
			msg: "message request no from",
		},
		{
			req: messageRequestNoTo,
			msg: "message request no to",
		},
		{
			req: messageRequestNoContent,
			msg: "message request no content",
		},
		{
			req: messageRequestNoMatchUUID,
			msg: "message request no match uuid",
		},
	}

	for _, t := range tests {
		_, err := messageController.AddMessage(t.req)
		suite.Require().Error(err, t.msg)
	}

}

func (suite *MessageServiceTestSuite) TestGetMessages() {

	mockCtrl := gomock.NewController(suite.T())
	defer mockCtrl.Finish()

	mockRepo := mockRepo.NewMockRepo(mockCtrl)
	messageController := MessageController{
		repo: mockRepo,
	}

	messageRequestSuccess := &types.MessageRequest{
		MatchUUID: mapping.StrToPtr("match-uuid"),
		Page:      mapping.IntToPtr(0),
	}

	messageRequestNoMatchUUID := &types.MessageRequest{
		Page: mapping.IntToPtr(0),
	}
	messageRequestNoPage := &types.MessageRequest{
		MatchUUID: mapping.StrToPtr("match-uuid"),
	}

	tests := []struct {
		req         *types.MessageRequest
		msg         string
		shouldError bool
	}{
		{
			req:         messageRequestSuccess,
			msg:         "get messages success",
			shouldError: false,
		},
		{
			req:         messageRequestNoMatchUUID,
			msg:         "get messages missing page",
			shouldError: true,
		},
		{
			req:         messageRequestNoPage,
			msg:         "get messages missing match uuid",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		if !tt.shouldError {
			msgs := []*types.Message{}
			mockRepo.EXPECT().GetMessagesByMatchUUID(gomock.Any(), gomock.Any(), gomock.Any()).Return(msgs, nil)
			msgs, err := messageController.GetMessages(tt.req, 3)
			suite.Require().NoError(err)
		} else {
			_, err := messageController.GetMessages(tt.req, 3)
			suite.Require().Error(err)
		}
	}

}
