package model

import (
	"database/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Outlet struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	MerchantID  uuid.UUID `gorm:"type:uuid"`
	Name        string    `gorm:"type:string;size:255"`
	Location    string
	PhoneNumber string `gorm:"type:string;size:13"`
	Audit
}

func (o *Outlet) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()

	o.Audit.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	o.Audit.ModifiedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return err
}

func (o *Outlet) BeforeUpdate(tx *gorm.DB) (err error) {
	o.Audit.ModifiedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return err
}
