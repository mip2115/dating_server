package repo

import (
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"

	//"github.com/kama/server/types"

	"testing"

	"code.mine/dating_server/DB"
	"github.com/stretchr/testify/suite"
)

type MessagesTestSuite struct {
	suite.Suite
}

func (suite *MessagesTestSuite) SetupSuite() {

}
func (suite *MessagesTestSuite) SetupTest() {

}

func (suite *MessagesTestSuite) TearDownAllSuite() {

}

func (suite *MessagesTestSuite) TearDownTest() {
}

func TestMessagesTestSuite(t *testing.T) {

	defer DB.Disconnect()
	suite.Run(t, new(MessagesTestSuite))
}

func (suite *MessagesTestSuite) TestSaveMessage() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	msg := &types.Message{
		MatchUUID: mapping.StrToPtr("match-uuid"),
	}
	err = SaveMessage(msg)
	suite.Require().NoError(err)
}

func (suite *MessagesTestSuite) TestGetMessagesByMatchUUID() {
	_, err := DB.GetTestDB()
	suite.Require().NoError(err)

	for i := 0; i < 2; i++ {
		msg := &types.Message{
			MatchUUID: mapping.StrToPtr("match-uuid"),
		}
		err = SaveMessage(msg)
		suite.Require().NoError(err)

	}

	msgs, err := GetMessagesByMatchUUID(0, 10, mapping.StrToPtr("match-uuid"))
	suite.Require().NoError(err)
	suite.Require().NotNil(msgs)
	suite.Require().Equal(2, len(msgs))

	for i := 0; i < 15; i++ {
		msg := &types.Message{
			MatchUUID: mapping.StrToPtr("match-uuid-2"),
		}
		err = SaveMessage(msg)
		suite.Require().NoError(err)
	}
	msgs, err = GetMessagesByMatchUUID(0, 10, mapping.StrToPtr("match-uuid-2"))
	suite.Require().NoError(err)
	suite.Require().NotNil(msgs)
	suite.Require().Equal(10, len(msgs))

	msgs, err = GetMessagesByMatchUUID(1, 10, mapping.StrToPtr("match-uuid-2"))
	suite.Require().NoError(err)
	suite.Require().NotNil(msgs)
	suite.Require().Equal(5, len(msgs))

	msgs, err = GetMessagesByMatchUUID(0, 10, mapping.StrToPtr("not-found"))
	suite.Require().NoError(err)
	suite.Require().NotNil(msgs)
	suite.Require().Equal(0, len(msgs))
}
