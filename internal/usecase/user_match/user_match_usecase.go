package usermatchusecase

import (
	"context"
	"date-apps-be/internal/constant"
	"date-apps-be/internal/model"
	userMatchRepo "date-apps-be/internal/repository/user_match"
	userusecase "date-apps-be/internal/usecase/user"
	"date-apps-be/internal/usecase/user_match/dto"
	"date-apps-be/pkg/derrors"
)

type (
	UserMatchUsecase interface {
		CreateUserMatch(ctx context.Context, userMatch *model.UserMatch) (err error)
		GetUserMatches(ctx context.Context, d dto.GetUserMatches) (userMatches []*model.UserMatch, err error)
		GetAvailableUsers(ctx context.Context, userUID string, page, limit uint64) (users []*model.User, quotaLeft int, err error)
		GetUserMatchTodayByUserUIDAndMatchUID(ctx context.Context, userUID, matchUID string) (userMatch *model.UserMatch, err error)
	}

	userMatchUsecase struct {
		repo        userMatchRepo.UserMatchRepository
		userUsecase userusecase.UserUsecase
	}
)

func NewUserMatchUsecase(repo userMatchRepo.UserMatchRepository, userUsecase userusecase.UserUsecase) UserMatchUsecase {
	return &userMatchUsecase{
		repo:        repo,
		userUsecase: userUsecase,
	}
}

// CreateUserMatch creates a new user match record in the database.
func (u *userMatchUsecase) CreateUserMatch(ctx context.Context, userMatch *model.UserMatch) (err error) {
	defer derrors.Wrap(&err, "CreateUserMatch(%q)", userMatch.UserUID)

	maxMatchPerDay := constant.MaxMatchPerDay

	userPackage, err := u.userUsecase.GetUserPackage(ctx, userMatch.UserUID)
	if err != nil {
		return
	}

	if userPackage == nil || userPackage.IsExpiredPackage() {
		// get total matches today for user
		total, err := u.repo.GetTotalUserMatchToday(ctx, userMatch.UserUID)
		if err != nil {
			return err
		}

		if total >= maxMatchPerDay {
			err = derrors.New(derrors.Forbidden, "Quota match per day reached")
			return err
		}
	} else if userPackage.Quota > 0 {
		maxMatchPerDay = int(userPackage.Quota)
		total, err := u.repo.GetTotalUserMatchToday(ctx, userMatch.UserUID)
		if err != nil {
			return err
		}

		if total >= maxMatchPerDay {
			err = derrors.New(derrors.Forbidden, "Quota match per day reached")
			return err
		}
	}

	return u.repo.CreateUserMatch(ctx, userMatch)
}

// GetUserMatches retrieves a list of user matches
// based on the user's UID
func (u *userMatchUsecase) GetUserMatches(ctx context.Context, d dto.GetUserMatches) (userMatches []*model.UserMatch, err error) {
	return u.repo.GetUserMatches(ctx, d)
}

// GetAvailableUsers retrieves a list of users that the current user has not matched with today
func (u *userMatchUsecase) GetAvailableUsers(ctx context.Context, userUID string, page, limit uint64) (users []*model.User, quotaLeft int, err error) {
	userPackage, err := u.userUsecase.GetUserPackage(ctx, userUID)
	if err != nil {
		return
	}

	total, err := u.repo.GetTotalUserMatchToday(ctx, userUID)
	if err != nil {
		return nil, 0, err
	}

	maxMatchPerDay := constant.MaxMatchPerDay
	quotaLeft = maxMatchPerDay - total

	if userPackage != nil && !userPackage.IsExpiredPackage() {
		if userPackage.Quota == 0 {
			// Unlimited matches for premium users with quota=0
			users, err = u.repo.GetAvailableUsers(ctx, userUID, page, limit)
			if err != nil {
				return
			}
			return users, 9999, nil
		}

		maxMatchPerDay = int(userPackage.Quota)
		quotaLeft = maxMatchPerDay - total
	}

	if total >= maxMatchPerDay {
		err = derrors.New(derrors.Forbidden, "Quota match per day reached")
		return nil, 0, err
	}

	users, err = u.repo.GetAvailableUsers(ctx, userUID, page, limit)
	if err != nil {
		return
	}
	return
}

func (u *userMatchUsecase) GetUserMatchTodayByUserUIDAndMatchUID(ctx context.Context, userUID, matchUID string) (userMatch *model.UserMatch, err error) {
	return u.repo.GetUserMatchTodayByUserUIDAndMatchUID(ctx, userUID, matchUID)
}
