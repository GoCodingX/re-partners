package db

import (
	"errors"

	"github.com/uptrace/bun/driver/pgdriver"
)

func IsUniqueViolation(err error) bool {
	var pgErr pgdriver.Error

	return errors.As(err, &pgErr) && pgErr.Field('C') == "23505"
}

func IsForeignKeyViolation(err error) bool {
	var pgErr pgdriver.Error

	return errors.As(err, &pgErr) && pgErr.Field('C') == "23503"
}
