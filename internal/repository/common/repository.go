package repository

import (
	"context"
	"database/sql"
	"date-apps-be/infrastructure/database"
	"date-apps-be/pkg/derrors"
	"date-apps-be/pkg/util"
	"fmt"
	"strings"
)

type repository struct {
	db *database.DB
}

// Repository general repository
type Repository interface {
	Master() *sql.DB
	Slave() *sql.DB
	Begin() (tx *sql.Tx, err error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
	Exec(ctx context.Context, tx *sql.Tx, query string, args []interface{}) (result sql.Result, err error)
	Query(ctx context.Context, query string, dest []interface{}, args []interface{}) error
	GetOffset(page uint64, limit uint64) uint64
	AddSortQuery(query string, allowedFields []string, sortBy string) (string, error)
	AddSortQueryWithPrefix(query string, allowedFields map[string]string, sortBy string) (string, error)
	NewNullString(str *string) sql.NullString
}

// NewRepository init general repository
func NewRepository(db *database.DB) Repository {
	return &repository{
		db: db,
	}
}

// Master return master
func (s *repository) Master() *sql.DB {
	return s.db.Master
}

// Slave return slave
func (s *repository) Slave() *sql.DB {
	return s.db.Slave
}

// Begin begin transaction
func (s *repository) Begin() (tx *sql.Tx, err error) {
	return s.db.Master.Begin()
}

// Commit commit transaction
func (s *repository) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

// Rollback rollback transaction
func (s *repository) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

// Insert Data
func (s *repository) Exec(ctx context.Context, tx *sql.Tx, query string, args []interface{}) (result sql.Result, err error) {
	if tx != nil {
		result, err = tx.ExecContext(ctx, query, args...)
	} else {
		result, err = s.Master().ExecContext(ctx, query, args...)
	}

	if err != nil {
		return result, err
	}

	return result, nil
}

// Select Data
func (s *repository) Query(ctx context.Context, query string, dest []interface{}, args []interface{}) (err error) {
	err = s.Master().QueryRowContext(ctx, query, args...).Scan(dest...)
	if err != nil {
		return err
	}

	return nil
}

func (s *repository) GetOffset(page uint64, limit uint64) uint64 {
	offset := (page - 1) * limit
	return offset
}

func (s *repository) validateAndReturnSortQuery(allowedFields []string, sortBy string) (string, error) {
	splits := strings.Split(sortBy, ".")
	if len(splits) != 2 {
		return "", derrors.New(derrors.InvalidArgument, "malformed sortBy query parameter, should be field.orderdirection")
	}
	field, order := splits[0], splits[1]
	if order != "desc" && order != "asc" {
		return "", derrors.New(derrors.InvalidArgument, "malformed orderdirection in sortBy query parameter, should be asc or desc")
	}
	if !util.StringInSlice(allowedFields, field) {
		return "", derrors.New(derrors.InvalidArgument, "unknown field in sortBy query parameter")
	}
	return fmt.Sprintf("%s %s", field, strings.ToUpper(order)), nil
}

func (s *repository) validateAndReturnSortQueryWithPrefix(allowedFields map[string]string, sortBy string) (string, error) {
	splits := strings.Split(sortBy, ".")
	if len(splits) != 2 {
		return "", derrors.New(derrors.InvalidArgument, "malformed sortBy query parameter, should be field.orderdirection")
	}
	field, order := splits[0], splits[1]
	if order != "desc" && order != "asc" {
		return "", derrors.New(derrors.InvalidArgument, "malformed orderdirection in sortBy query parameter, should be asc or desc")
	}
	if allowedFields[field] == "" {
		return "", derrors.New(derrors.InvalidArgument, "unknown field in sortBy query parameter")
	}
	return fmt.Sprintf("%s %s", allowedFields[field], strings.ToUpper(order)), nil
}

func (s *repository) AddSortQuery(query string, allowedFields []string, sortBy string) (string, error) {
	sortQuery, err := s.validateAndReturnSortQuery(allowedFields, sortBy)
	if err != nil {
		return "", err
	}
	return query + " ORDER BY " + sortQuery, nil
}

func (s *repository) AddSortQueryWithPrefix(query string, allowedFields map[string]string, sortBy string) (string, error) {
	sortQuery, err := s.validateAndReturnSortQueryWithPrefix(allowedFields, sortBy)
	if err != nil {
		return "", err
	}
	return query + " ORDER BY " + sortQuery, nil
}

func (s *repository) NewNullString(str *string) sql.NullString {
	if str == nil {
		return sql.NullString{}
	}

	if len(*str) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: *str,
		Valid:  true,
	}
}
