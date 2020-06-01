package scripts

import (
	"context"
	"github.com/kama/server/DB"
	"github.com/kama/server/mapping"
	"github.com/kama/server/types"
)

func CreateUsers() error {
	err := DropCollectionUsers()
	if err != nil {
		return err
	}
	u1 := types.User{
		UUID:          mapping.StrToPtr("eab85cb1-0a11-47d1-890d-93015dc1e6fz"),
		Religion:      mapping.StrToPtr("judaism"),
		City:          mapping.StrToPtr("New York"),
		University:    mapping.StrToPtr("New York"),
		Politics:      mapping.StrToPtr("conservative"),
		Smoke:         mapping.StrToPtr("sometimes"),
		Drink:         mapping.StrToPtr("sometimes"),
		Email:         mapping.StrToPtr("testuser8@test.com"),
		Password:      mapping.StrToPtr("testing123"),
		Mobile:        mapping.StrToPtr("0001112233"),
		Job:           mapping.StrToPtr("cleaner"),
		Age:           mapping.IntToPtr(20),
		PartnerGender: mapping.StrToPtr("male"),
		Gender:        mapping.StrToPtr("female"),
	}
	u2 := types.User{
		UUID:          mapping.StrToPtr("eab85cb1-0a11-47d1-890d-93015dc1e6fa"),
		Religion:      mapping.StrToPtr("judaism"),
		City:          mapping.StrToPtr("New York"),
		University:    mapping.StrToPtr("New York University"),
		Politics:      mapping.StrToPtr("conservative"),
		Smoke:         mapping.StrToPtr("sometimes"),
		Drink:         mapping.StrToPtr("sometimes"),
		Email:         mapping.StrToPtr("testuser8@test.com"),
		Password:      mapping.StrToPtr("testing123"),
		Mobile:        mapping.StrToPtr("0001112233"),
		Job:           mapping.StrToPtr("cleaner"),
		Age:           mapping.IntToPtr(20),
		PartnerGender: mapping.StrToPtr("female"),
		Gender:        mapping.StrToPtr("male"),
	}
	u3 := types.User{
		UUID:          mapping.StrToPtr("eab85cb1-0a11-47d1-890d-93015dc1e6fb"),
		Religion:      mapping.StrToPtr("hindu"),
		City:          mapping.StrToPtr("New York"),
		University:    mapping.StrToPtr("Columbia"),
		Politics:      mapping.StrToPtr("conservative"),
		Smoke:         mapping.StrToPtr("sometimes"),
		Drink:         mapping.StrToPtr("sometimes"),
		Email:         mapping.StrToPtr("testuser8@test.com"),
		Password:      mapping.StrToPtr("testing123"),
		Mobile:        mapping.StrToPtr("0001112233"),
		Job:           mapping.StrToPtr("cleaner"),
		Age:           mapping.IntToPtr(20),
		PartnerGender: mapping.StrToPtr("female"),
		Gender:        mapping.StrToPtr("male"),
	}
	u4 := types.User{
		UUID:          mapping.StrToPtr("eab85cb1-0a11-47d1-890d-93015dc1e6fc"),
		Religion:      mapping.StrToPtr("judaism"),
		City:          mapping.StrToPtr("New York"),
		University:    mapping.StrToPtr("Columbia"),
		Politics:      mapping.StrToPtr("conservative"),
		Smoke:         mapping.StrToPtr("sometimes"),
		Drink:         mapping.StrToPtr("sometimes"),
		Email:         mapping.StrToPtr("testuser8@test.com"),
		Password:      mapping.StrToPtr("testing123"),
		Mobile:        mapping.StrToPtr("0001112233"),
		Job:           mapping.StrToPtr("cleaner"),
		Age:           mapping.IntToPtr(20),
		PartnerGender: mapping.StrToPtr("female"),
		Gender:        mapping.StrToPtr("male"),
	}
	u5 := types.User{
		UUID:          mapping.StrToPtr("eab85cb1-0a11-47d1-890d-93015dc1e6fp"),
		Religion:      mapping.StrToPtr("christian"),
		City:          mapping.StrToPtr("New York"),
		University:    mapping.StrToPtr("Columbia"),
		Politics:      mapping.StrToPtr("liberal"),
		Smoke:         mapping.StrToPtr("never"),
		Drink:         mapping.StrToPtr("never"),
		Email:         mapping.StrToPtr("testuser9@test.com"),
		Password:      mapping.StrToPtr("testing123"),
		Mobile:        mapping.StrToPtr("0001112233"),
		Job:           mapping.StrToPtr("janitor"),
		Age:           mapping.IntToPtr(20),
		PartnerGender: mapping.StrToPtr("male"),
		Gender:        mapping.StrToPtr("female"),
	}
	c, err := DB.GetCollection("users")
	if err != nil {
		return err
	}
	_, err = c.InsertMany(
		context.Background(),
		[]interface{}{
			u1, u2, u3, u4, u5,
		})
	if err != nil {
		return err
	}
	return nil
}

func DropCollectionUsers() error {
	c, err := DB.GetCollection("users")
	if err != nil {
		return err
	}
	err = c.Drop(context.Background())
	if err != nil {
		return err
	}
	return nil
}
