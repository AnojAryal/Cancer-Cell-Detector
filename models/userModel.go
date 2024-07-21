package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username        string `gorm:"uniqueIndex"`
	Email           string `gorm:"uniqueIndex"`
	FullName        string
	Address         string
	BloodGroup      string
	Gender          string
	ContactNo       string
	Password        string `gorm:"not null"`
	IsVerified      bool   `gorm:"default:false"`
	IsAdmin         bool   `gorm:"default:false"`
	IsHospitalAdmin bool   `gorm:"default:false"`
	HospitalID      uint
	Hospital        *Hospital `gorm:"foreignKey:HospitalID"`
}

type PasswordResetToken struct {
	gorm.Model
	Email string `gorm:"index"`
	Token string `gorm:"uniqueIndex"`
	Used  bool   `gorm:"default:false"`
}
