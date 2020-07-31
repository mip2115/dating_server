package article_service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"code.mine/dating_server/DB"
	"code.mine/dating_server/mapping"
	"code.mine/dating_server/types"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"

	"fmt"
)

const (
	minWords = 300
	maxWords = 3000
)

func DeleteArticleRecordByUUID(uuid *string) error {
	c, err := DB.GetCollection("articles")
	if err != nil {
		return err
	}

	_, err = c.DeleteOne(context.Background(), bson.M{"uuid": uuid})
	if err != nil {
		return err
	}
	return nil
}

// if user uses a link, you should still end up calling this at a later point
func CreateArticleRecord(userUUID *string, text *string) error {
	wordCount := len(mapping.StrToV(text))
	if wordCount > maxWords {
		return fmt.Errorf("article has word count of %d which exceeds max words of %d", wordCount, maxWords)
	}
	if wordCount < minWords {
		return fmt.Errorf("article has word count of %d which is less than min words of %d", wordCount, minWords)
	}
	if userUUID == nil {
		return errors.New("user uuid cannot be nil")
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	article := &types.Article{
		UserUUID:    userUUID,
		UUID:        mapping.StrToPtr(uuid.String()),
		WordCount:   mapping.Int64ToPtr(int64(wordCount)),
		DateCreated: mapping.TimeToPtr(time.Now()),
		Text:        text,
	}

	c, err := DB.GetCollection("articles")
	if err != nil {
		return err
	}

	_, err = c.InsertOne(context.Background(), article)
	if err != nil {
		return err
	}

	return nil
}

// pass in a user uuid
// get that user's article and profile
// get a list of all users in that person's city and gender pref's
// and age range.
// then get all THEIR articles and that's what you'll match on.

// https://cloud.google.com/natural-language/docs/reference/libraries#client-libraries-install-go
// https://github.com/gopherdata/resources/blob/master/tooling/README.md#nlp
func GetArticleMatchesForUser(userUUID *string) ([]*types.Article, error) {
	var user *types.User
	var userCandidates []*types.User

	c, err := DB.GetCollection("users")
	if err != nil {
		return nil, err
	}
	res := c.FindOne(context.Background(), bson.M{"uuid": mapping.StrToV(userUUID)})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no user found for user uuid ", mapping.StrToV(userUUID))
		}
		return nil, res.Err()
	}
	err = res.Decode(user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user is nil for user uuid ", mapping.StrToV(userUUID))
	}

	// fetch all users by filter
	filter := []bson.M{
		bson.M{"city": mapping.StrToV(user.City)},
		bson.M{"gender": mapping.StrToV(user.PartnerGender)},
	}
	cursor, err := c.Find(context.Background(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no candidates found for user uuid ", mapping.StrToV(userUUID))
		}
		return nil, err
	}
	err = cursor.Decode(userCandidates)
	if err != nil {
		return nil, err
	}
	if len(userCandidates) == 0 {
		return nil, fmt.Errorf("no candidates found for user uuid ", mapping.StrToV(userUUID))
	}

	// get the users articles
	var userArticles []*types.Article
	var candidateArticles []*types.Article

	c, err = DB.GetCollection("articles")
	if err != nil {
		return nil, err
	}
	cursor, err = c.Find(context.Background(), bson.M{"user_uuid": mapping.StrToV(userUUID)})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no articles found for user uuid ", mapping.StrToV(userUUID))
		}
		return nil, err
	}
	err = res.Decode(userArticles)
	if err != nil {
		return nil, err
	}
	if len(userArticles) == 0 {
		return nil, fmt.Errorf("no articles found for user uuid ", mapping.StrToV(userUUID))
	}

	articleUserUUIDs := []bson.M{}
	for _, candidate := range userCandidates {
		articleUserUUIDs = append(articleUserUUIDs, bson.M{"user_uuid": mapping.StrToV(candidate.UUID)})
	}
	cursor, err = c.Find(context.Background(), bson.M{"$or": articleUserUUIDs})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no articles found for user candidates for user uuid ", mapping.StrToV(userUUID))
		}
		return nil, err
	}
	err = cursor.Decode(candidateArticles)
	if err != nil {
		return nil, err
	}
	return candidateArticles, nil
}
