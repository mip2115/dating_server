package search_service

import (
	"context"

	"code.mine/dating_server/DB"
	us "code.mine/dating_server/service/user_service"
	"code.mine/dating_server/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	limitValue    = 3
	usersToReturn = 3
)

func CalculateTopUsers(userUUID *string, skipValue int) ([]types.User, error) {
	user, err := us.GetUser(userUUID)
	if err != nil {
		return nil, err
	}
	u := types.User{}
	/*
		topMatchedUsers, err := calculateTopUsers(
			user.University, user.City, user.Drink, user.Smoke, user.PartnerGender,
			user.Religion, user.Politics, skipValue,
		)
		if err != nil {
			return nil, err
		}
	*/
	topMatchedUsers := []types.User{u}
	return topMatchedUsers, nil

}

func calculateTopUsers(
	university *string, city *string, drink *string, smoke *string, partnerGender *string,
	religion *string, politics *string, skipValue int,
) ([]types.User, error) {

	c, err := DB.GetCollection("users")
	if err != nil {
		return nil, err
	}

	/*
			matchStage := bson.D{{"$match", bson.D{{"podcast", id}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", "$podcast"}, {"total", bson.D{{"$sum", "$duration"}}}}}}

	*/
	//matchStage1 := bson.M{"$match": bson.M{"city": city}}
	matchStage1 := bson.D{
		{
			"$match", bson.D{
				{"city", *city},
			},
		},
	}
	// matchStage2 := bson.M{"$match": bson.M{"gender": partnerGender}}
	matchStage2 := bson.D{
		{
			"$match", bson.D{
				{"gender", *partnerGender},
			},
		},
	}

	projectStage1 := bson.D{{
		"$project", bson.D{
			{"uuid", 1}, {"religion", 1}, {"profile_image", 1},
			{"job", 1}, {"age", 1}, {"purpose", 1},
			{"first_name", 1},
			{"religionScore", bson.M{"$toInt": 0}},
			{"politics", 1}, {"politicsScore", bson.M{"$toInt": 0}},
			{"university", 1}, {"universityScore", bson.M{"$toInt": 0}},
			{"drink", 1}, {"drinkScore", bson.M{"$toInt": 0}},
			{"smoke", 1}, {"smokeScore", bson.M{"$toInt": 0}},
		},
	}}

	religionScore := bson.M{
		"$cond": []interface{}{
			bson.M{"$eq": []interface{}{"$religion", *religion}},
			4, 0,
		},
	}
	politicsScore := bson.M{
		"$cond": []interface{}{
			bson.M{"$eq": []interface{}{"$politics", *politics}},
			4, 0,
		},
	}
	universityScore := bson.M{
		"$cond": []interface{}{
			bson.M{"$eq": []interface{}{"$university", *university}},
			3, 0,
		},
	}
	smokeScore := bson.M{
		"$cond": []interface{}{
			bson.M{"$eq": []interface{}{"$smoke", *smoke}},
			1, 0,
		},
	}
	drinkScore := bson.M{
		"$cond": []interface{}{
			bson.M{"$eq": []interface{}{"$drink", *drink}},
			1, 0,
		},
	}
	projectStage2 := bson.D{{"$project", bson.D{
		{"uuid", 1}, {"religion", 1}, {"profile_image", 1},
		{"first_name", 1},
		{"job", 1}, {"age", 1}, {"purpose", 1}, {"university", 1},
		{"politics", 1}, {"smoke", 1}, {"drink", 1},
		{"religionScore", religionScore},
		{"politicsScore", politicsScore}, {"universityScore", universityScore},
		{"smokeScore", smokeScore}, {"drinkScore", drinkScore},
	}},
	}

	projectStage3 := bson.D{{"$project", bson.M{
		"uuid": 1, "religion": 1, "profile_image": 1,
		"first_name": 1,
		"job":        1, "age": 1, "purpose": 1, "university": 1,
		"politics": 1, "smoke": 1, "drink": 1,
		"score": bson.M{
			"$sum": []interface{}{"$religionScore", "$politicsScore", "$universityScore", "$drinkScore", "$smokeScore"},
		},
	}}}

	sortStage := bson.D{{"$sort", bson.M{
		"score": -1,
	}}}

	limitStage := bson.D{{"$limit", limitValue}}
	skipStage := bson.D{{"$skip", skipValue}}

	/*
		projectStage1 := bson.M{
			"$project": bson.M{
				"uuid": 1, "religion": 1, "profile_image": 1,
				"job": 1, "age": 1, "purpose": 1,
				"first_name":    1,
				"religionScore": bson.M{"$toInt": 0},
				"politics":      1, "politicsScore": bson.M{"$toInt": 0},
				"university": 1, "universityScore": bson.M{"$toInt": 0},
				"drink": 1, "drinkScore": bson.M{"$toInt": 0},
				"smoke": 1, "smokeScore": bson.M{"$toInt": 0},
			},
		}
		religionScore := bson.M{
			"$cond": []interface{}{
				bson.M{"$eq": []interface{}{"$religion", *religion}},
				4, 0,
			},
		}
		politicsScore := bson.M{
			"$cond": []interface{}{
				bson.M{"$eq": []interface{}{"$politics", *politics}},
				4, 0,
			},
		}
		universityScore := bson.M{
			"$cond": []interface{}{
				bson.M{"$eq": []interface{}{"$university", *university}},
				3, 0,
			},
		}
		smokeScore := bson.M{
			"$cond": []interface{}{
				bson.M{"$eq": []interface{}{"$smoke", *smoke}},
				1, 0,
			},
		}
		drinkScore := bson.M{
			"$cond": []interface{}{
				bson.M{"$eq": []interface{}{"$drink", *drink}},
				1, 0,
			},
		}
		projectStage2 := bson.M{"$project": bson.M{
			"uuid": 1, "religion": 1, "profile_image": 1,
			"first_name": 1,
			"job":        1, "age": 1, "purpose": 1, "university": 1,
			"politics": 1, "smoke": 1, "drink": 1,
			"religionScore": religionScore,
			"politicsScore": politicsScore, "universityScore": universityScore,
			"smokeScore": smokeScore, "drinkScore": drinkScore,
		},
		}
		projectStage3 := bson.M{"$project": bson.M{
			"uuid": 1, "religion": 1, "profile_image": 1,
			"first_name": 1,
			"job":        1, "age": 1, "purpose": 1, "university": 1,
			"politics": 1, "smoke": 1, "drink": 1,
			"score": bson.M{
				"$sum": []interface{}{"$religionScore", "$politicsScore", "$universityScore", "$drinkScore", "$smokeScore"},
			},
		}}
		sortStage := bson.M{"$sort": bson.M{
			"score": -1,
		}}
		limitStage := bson.M{"$limit": limitValue}
		skipStage := bson.M{"$skip": skipValue}
	*/
	var users []types.User
	//var res []bson.D
	showInfoCursor, err := c.Aggregate(context.Background(), mongo.Pipeline{matchStage1, matchStage2, projectStage1, projectStage2, projectStage3, sortStage, skipStage, limitStage})
	if err != nil {
		return nil, err
	}
	err = showInfoCursor.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}
	/*
			pipeline := c.Pipe([]bson.M{
				matchStage1, matchStage2, projectStage1, projectStage2, projectStage3, sortStage, skipStage, limitStage,
			})


		err = pipeline.All(&users)
		if err != nil {
			return nil, err
		}
	*/
	return users, nil
}
