package model

import (
	"database/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	OutletID    uuid.UUID `gorm:"type:uuid"`
	Name        string    `gorm:"type:string;size:255"`
	Description string    `gorm:"type:string;size:255"`
	Stock       int64
	Price       float64
	Image       string `gorm:"type:string;size:255"`
	Audit
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()

	p.Audit.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	p.Audit.ModifiedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return err
}

func (p *Product) BeforeUpdate(tx *gorm.DB) (err error) {
	p.Audit.ModifiedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return err
}
