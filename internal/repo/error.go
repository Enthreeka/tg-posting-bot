package repo

import (
	"errors"
	customErr "github.com/Enthreeka/tg-posting-bot/pkg/bot_error"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrRecordNotFound = pgx.ErrNoRows

var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

func ErrorHandler(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return customErr.ErrNoRows
	}
	errCode := ErrorCode(err)
	if errCode == ForeignKeyViolation {
		return customErr.ErrForeignKeyViolation
	}
	if errCode == UniqueViolation {
		return customErr.ErrUniqueViolation
	}

	return err
}
