package models

import "database/sql"

type Models struct {
	UserModel UserModel
}

func New(db *sql.DB) *Models {
	return &Models{
		UserModel: NewUserModel(db),
	}
}
