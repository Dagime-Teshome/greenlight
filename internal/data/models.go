package data

import (
	"database/sql"
	"errors"
)

type Models struct {
	Movies      MovieModel
	Users       UserModel
	Tokens      TokenModel
	Permissions PermissionsModel
}

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

func NewModel(db *sql.DB) Models {

	return Models{
		Movies:      MovieModel{DB: db},
		Permissions: PermissionsModel{DB: db},
		Users:       UserModel{DB: db},
		Tokens:      TokenModel{DB: db},
	}
}
