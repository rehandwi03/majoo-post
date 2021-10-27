package model

import (
	"database/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Merchant struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID          uuid.UUID `gorm:"type:uuid"`
	Name            string    `gorm:"type:string;size:255"`
	InstitutionName string    `gorm:"type:string;size:255"`
	PhoneNumber     string    `gorm:"type:string;size:13"`
	Audit
}

func (m *Merchant) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()

	m.Audit.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	m.Audit.ModifiedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return err
}

func (m *Merchant) BeforeUpdate(tx *gorm.DB) (err error) {
	m.Audit.ModifiedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return err
}
