package model

import (
	"database/sql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	FirstName   string    `gorm:"type:string;size:50"`
	LastName    string    `gorm:"type:string;size:50"`
	Email       string    `gorm:"type:string;size:50"`
	Password    string    `gorm:"type:string;size:255"`
	PhoneNumber string    `gorm:"type:string;size:13"`
	Audit
}

func (u *User) ComparePassword(password string) (status bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (u *User) EncryptPassword() (err error) {
	newPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(newPassword)

	return
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	newPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(newPassword)
	u.Audit.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	u.Audit.ModifiedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return err
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.Audit.ModifiedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return err
}
