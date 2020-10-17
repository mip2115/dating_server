package matchservice

import (
	"errors"
	"time"

	"code.mine/dating_server/repo"
	"code.mine/dating_server/types"
)

// MatchController -
type MatchController struct {
	repo repo.Repo
}

// TODO - error check to make sure youre not matching people who recently matched
// or who should be blocked

// CreateMatch -
func (c *MatchController) CreateMatch(m *types.Match) (*types.Match, error) {
	if m.UserOneUUID == nil {
		return nil, errors.New("userOneUUID missing")
	}
	if m.UserTwoUUID == nil {
		return nil, errors.New("userTwoUUID missing")
	}
	t := time.Now()
	m.DateCreated = &t
	m.DateUpdated = &t
	createdMatch, err := c.repo.CreateMatch(m)
	if err != nil {
		return nil, err
	}
	return createdMatch, nil
}

// DeleteMatch -
func (c *MatchController) DeleteMatch(matchUUID *string) error {
	err := c.repo.DeleteMatchByMatchUUID(matchUUID)
	if err != nil {
		return err
	}
	err = c.repo.DeleteTrackedLikeByMatchUUID(matchUUID)
	if err != nil {
		return err
	}
	return nil
}
