package userpremiumrepository

import (
	"context"
	"database/sql"
	"date-apps-be/internal/model"
	repository "date-apps-be/internal/repository/common"
	"date-apps-be/pkg/derrors"
)

type UserPremiumRepository interface {
	repository.Repository
	CreateUserPackage(ctx context.Context, tx *sql.Tx, userPackage *model.UserPackage) (err error)
	GetUserPackage(ctx context.Context, userUID string) (userPackage *model.UserPackage, err error)
}

type userPremiumRepository struct {
	repository.Repository
}

func NewUserPremiumRepository(repo repository.Repository) UserPremiumRepository {
	return &userPremiumRepository{
		Repository: repo,
	}
}

func (u *userPremiumRepository) getDest(userPremium *model.UserPackage) []interface{} {
	return []interface{}{
		&userPremium.UID,
		&userPremium.UserUID,
		&userPremium.PremiumConfig.UID,
		&userPremium.PremiumConfig.Name,
		&userPremium.PremiumConfig.Description,
		&userPremium.PremiumConfig.Price,
		&userPremium.PremiumConfig.Quota,
		&userPremium.PremiumConfig.ExpiredDay,
		&userPremium.StartedAt,
		&userPremium.EndedAt,
		&userPremium.Quota,
	}
}
func (u *userPremiumRepository) GetUserPackage(ctx context.Context, userUID string) (userPackage *model.UserPackage, err error) {
	defer derrors.Wrap(&err, "GetUserPackage(%q)", userUID)

	query := `SELECT 
		up.uid, 
		up.user_uid, 
		up.premium_config_uid, 
		pc.name, 
		pc.description, 
		pc.price, 
		pc.quota, 
		pc.expired_day, 
		up.started_at, 
		up.ended_at, 
		up.quota
	FROM 
		user_premium up
	JOIN 
		premium_config pc ON up.premium_config_uid = pc.uid
	WHERE 
		up.user_uid = ?`

	args := []interface{}{
		userUID,
	}
	userPackage = &model.UserPackage{
		PremiumConfig: &model.PremiumConfig{}, // Ensure PremiumConfig is initialized
	}

	err = u.Query(ctx, query, u.getDest(userPackage), args)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, derrors.HandleSQLError(err, "r.Query")
	}

	return userPackage, nil
}

func (u *userPremiumRepository) CreateUserPackage(ctx context.Context, tx *sql.Tx, userPackage *model.UserPackage) (err error) {
	defer derrors.Wrap(&err, "CreateUserPackage(%v)", userPackage)

	query := `INSERT INTO user_premium (uid, user_uid, premium_config_uid, started_at, ended_at, quota) VALUES (?,?,?,?,?,?)`

	args := []interface{}{
		userPackage.UID,
		userPackage.UserUID,
		userPackage.PremiumConfigUID,
		userPackage.StartedAt,
		userPackage.EndedAt,
		userPackage.Quota,
	}

	_, err = u.Exec(ctx, tx, query, args)
	if err != nil {
		return derrors.WrapStack(err, derrors.Unknown, "r.Exec")
	}

	return nil
}
