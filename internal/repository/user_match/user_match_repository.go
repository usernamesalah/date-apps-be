package usermatchrepository

import (
	"context"
	"date-apps-be/internal/model"
	repository "date-apps-be/internal/repository/common"
	"date-apps-be/internal/usecase/user_match/dto"
	"date-apps-be/pkg/derrors"
)

type UserMatchRepository interface {
	repository.Repository
	CreateUserMatch(ctx context.Context, userMatch *model.UserMatch) (err error)
	GetUserMatches(ctx context.Context, d dto.GetUserMatches) (userMatches []*model.UserMatch, err error)
	GetTotalUserMatchToday(ctx context.Context, userUID string) (total int, err error)
	GetAvailableUsers(ctx context.Context, userUID string, page, limit uint64) (users []*model.User, err error)
}

type userMatchRepository struct {
	repository.Repository
}

func NewUserMatchRepository(store repository.Repository) UserMatchRepository {
	return &userMatchRepository{Repository: store}
}

func (u *userMatchRepository) getDest(userMatch *model.UserMatch) []interface{} {
	return []interface{}{
		&userMatch.UserUID,
		&userMatch.Match.UID,
		&userMatch.Match.Name,
		&userMatch.MatchType,
		&userMatch.CreatedAt,
	}
}

func (u *userMatchRepository) GetUserMatches(ctx context.Context, d dto.GetUserMatches) (userMatches []*model.UserMatch, err error) {
	defer derrors.Wrap(&err, "GetUserMatches(%q)", d.UserUID)

	query := `SELECT user_uid , match_uid, name, match_type FROM user_matches 
			LEFT JOIN users ON user_matches.match_uid = users.uid
			WHERE user_uid = ?`

	var args = []interface{}{
		d.UserUID,
	}
	if d.MatchType.IsValid() {
		query += ` AND match_type = ? `
		args = append(args, d.MatchType.String())
	}

	query += ` LIMIT ?,?`

	args = append(args, u.GetOffset(d.Page, d.Limit), d.Limit)

	userMatches = []*model.UserMatch{}

	rows, err := u.Slave().QueryContext(ctx, query, args...)
	if err != nil {
		err = derrors.HandleSQLError(err, "QueryContext")
		return
	}

	for rows.Next() {
		wc := &model.UserMatch{}
		err = rows.Scan(u.getDest(wc)...)
		if err != nil {
			return nil, err
		}

		userMatches = append(userMatches, wc)
	}

	return userMatches, nil
}

func (u *userMatchRepository) CreateUserMatch(ctx context.Context, userMatch *model.UserMatch) (err error) {
	defer derrors.Wrap(&err, "CreateUserMatch(%v)", userMatch)

	query := `INSERT INTO user_matches (user_uid, match_uid, match_type) VALUES (?, ?, ?)`
	args := []interface{}{
		userMatch.UserUID,
		userMatch.Match.UID,
		userMatch.MatchType,
	}

	_, err = u.Master().ExecContext(ctx, query, args...)
	if err != nil {
		err = derrors.HandleSQLError(err, "ExecContext")
		return
	}

	return nil
}

func (u *userMatchRepository) GetAvailableUsers(ctx context.Context, userUID string, page, limit uint64) (users []*model.User, err error) {
	defer derrors.Wrap(&err, "GetAvailableUsers(%q)", userUID)

	query := `SELECT u.uid, u.name FROM users u
			WHERE u.uid NOT IN (
				SELECT match_uid FROM user_matches
				WHERE user_uid = ? AND DATE(created_at) = CURDATE()
			)
			ORDER BY RAND()
			LIMIT ?,?`

	args := []interface{}{
		userUID,
		u.GetOffset(page, limit), limit,
	}

	users = []*model.User{}

	rows, err := u.Slave().QueryContext(ctx, query, args...)
	if err != nil {
		err = derrors.HandleSQLError(err, "QueryContext")
		return
	}

	for rows.Next() {
		user := &model.User{}
		err = rows.Scan(&user.UID, &user.Name)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (u *userMatchRepository) GetTotalUserMatchToday(ctx context.Context, userUID string) (total int, err error) {
	defer derrors.Wrap(&err, "GetTotalUserMatchToday(%q)", userUID)

	query := `SELECT COUNT(*) FROM user_matches
			WHERE user_uid = ? AND DATE(created_at) = CURDATE()`

	err = u.Slave().QueryRowContext(ctx, query, userUID).Scan(&total)
	if err != nil {
		err = derrors.HandleSQLError(err, "QueryRowContext")
		return
	}

	return total, nil
}
