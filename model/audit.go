package model

import (
	"database/sql"
	"gorm.io/gorm"
)

type Audit struct {
	CreatedAt sql.NullTime
	ModifiedAt sql.NullTime
	DeletedAt gorm.DeletedAt
}
