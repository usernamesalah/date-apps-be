package usermatchusecase

import (
	"context"
	"date-apps-be/internal/constant"
	"date-apps-be/internal/model"
	userMatchRepo "date-apps-be/internal/repository/user_match"
	"date-apps-be/internal/usecase/user_match/dto"
	"date-apps-be/pkg/derrors"
)

type (
	UserMatchUsecase interface {
		CreateUserMatch(ctx context.Context, userMatch *model.UserMatch) (err error)
		GetUserMatches(ctx context.Context, d dto.GetUserMatches) (userMatches []*model.UserMatch, err error)
		GetAvailableUsers(ctx context.Context, userUID string, page, limit uint64) (users []*model.User, err error)
	}

	userMatchUsecase struct {
		repo userMatchRepo.UserMatchRepository
	}
)

func NewUserMatchUsecase(repo userMatchRepo.UserMatchRepository) UserMatchUsecase {
	return &userMatchUsecase{
		repo: repo,
	}
}

// CreateUserMatch creates a new user match record in the database.
// It first validates the user's premium package (TODO).
// Then, it retrieves the total number of matches the user has made today.
// If the user has reached the maximum number of matches allowed per day,
// it returns a forbidden error. Otherwise, it creates the user match record.
//
// Parameters:
//   - ctx: The context for managing request-scoped values, cancellation, and deadlines.
//   - userMatch: A pointer to the UserMatch model containing the match details.
//
// Returns:
//   - err: An error if the operation fails, otherwise nil.
func (u *userMatchUsecase) CreateUserMatch(ctx context.Context, userMatch *model.UserMatch) (err error) {
	defer derrors.Wrap(&err, "CreateUserMatch(%q)", userMatch.UserUID)

	// TODO: validate user premium package here

	// get total matches today for user
	total, err := u.repo.GetTotalUserMatchToday(ctx, userMatch.UserUID)
	if err != nil {
		return
	}

	if total >= constant.MaxMatchPerDay {
		err = derrors.New(derrors.Forbidden, "Quota match per day reached")
	}

	return u.repo.CreateUserMatch(ctx, userMatch)
}

// GetUserMatches retrieves a list of user matches
// based on the user's UID
func (u *userMatchUsecase) GetUserMatches(ctx context.Context, d dto.GetUserMatches) (userMatches []*model.UserMatch, err error) {
	return u.repo.GetUserMatches(ctx, d)
}

// GetAvailableUsers retrieves a list of users that the current user has not matched with today
func (u *userMatchUsecase) GetAvailableUsers(ctx context.Context, userUID string, page, limit uint64) (users []*model.User, err error) {
	return u.repo.GetAvailableUsers(ctx, userUID, page, limit)
}
