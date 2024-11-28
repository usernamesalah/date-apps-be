package userrepository

import (
	"context"
	"database/sql"
	"date-apps-be/internal/model"
	repository "date-apps-be/internal/repository/common"
	"date-apps-be/pkg/derrors"
)

type (
	userRepository struct {
		repository.Repository
	}

	UserRepository interface {
		repository.Repository
		CreateUser(ctx context.Context, tx *sql.Tx, user *model.User) (id int64, err error)
		DeleteUser(ctx context.Context, tx *sql.Tx, id string) error
		GetUserByUID(ctx context.Context, id string) (user *model.User, err error)
		GetUserByEmailOrPhoneNumber(ctx context.Context, email, phoneNumber string) (user *model.User, err error)
		UpdateUser(ctx context.Context, tx *sql.Tx, user *model.User) error
	}
)

func NewUserRepository(store repository.Repository) UserRepository {
	return &userRepository{
		Repository: store,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, tx *sql.Tx, user *model.User) (id int64, err error) {
	defer derrors.Wrap(&err, "CreateUser(%q)", user.Email)

	query := `INSERT INTO users (uid, name, email, phone_number, password) VALUES (?, ?, ?, ?, ?)`
	args := []interface{}{
		user.UID,
		user.Name,
		r.NewNullString(&user.Email),
		r.NewNullString(&user.PhoneNumber),
		user.Password,
	}

	result, err := r.Exec(ctx, tx, query, args)
	if err != nil {
		return 0, derrors.WrapStack(err, derrors.Unknown, "r.Exec")
	}

	return result.LastInsertId()
}

func (r *userRepository) GetUserByUID(ctx context.Context, uid string) (user *model.User, err error) {
	defer derrors.Wrap(&err, "GetUserByUID(%q)", uid)

	query := `SELECT uid, name, email, phone_number, password FROM users WHERE uid = ?`
	user = &model.User{}
	dest := []interface{}{
		&user.UID,
		&user.Name,
		&user.Email,
		&user.PhoneNumber,
		&user.Password,
	}

	args := []interface{}{
		uid,
	}

	err = r.Query(ctx, query, dest, args)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, derrors.HandleSQLError(err, "r.Query")
	}

	return user, nil
}

func (r *userRepository) GetUserByEmailOrPhoneNumber(ctx context.Context, email, phoneNumber string) (user *model.User, err error) {
	defer derrors.Wrap(&err, "GetUserByEmailOrPhoneNumber(%q, %q)", email, phoneNumber)

	query := `SELECT uid, name, email, phone_number, password FROM users WHERE email = ? OR phone_number = ?`
	user = &model.User{}
	dest := []interface{}{
		&user.UID,
		&user.Name,
		&user.Email,
		&user.PhoneNumber,
		&user.Password,
	}

	args := []interface{}{
		email,
		phoneNumber,
	}

	err = r.Query(ctx, query, dest, args)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, derrors.HandleSQLError(err, "r.Query")
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, tx *sql.Tx, user *model.User) (err error) {
	defer derrors.Wrap(&err, "UpdateUser(%q)", user.UID)

	query := `UPDATE users SET name = ?, email = ?, phone_number = ?, password = ? WHERE uid = ?`
	args := []interface{}{
		user.Name,
		user.Email,
		user.PhoneNumber,
		user.Password,
		user.UID,
	}

	_, err = r.Exec(ctx, tx, query, args)
	if err != nil {
		return derrors.WrapStack(err, derrors.Unknown, "r.Exec")
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, tx *sql.Tx, id string) (err error) {
	defer derrors.Wrap(&err, "DeleteUser(%q)", id)

	query := `DELETE FROM users WHERE uid = ?`
	args := []interface{}{
		id,
	}

	_, err = r.Exec(ctx, tx, query, args)
	if err != nil {
		return derrors.WrapStack(err, derrors.Unknown, "r.Exec")
	}

	return nil
}
