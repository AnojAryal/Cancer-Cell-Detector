package models

import (
	"gorm.io/gorm"
)


type Hospital struct {
    gorm.Model
    Name     string    `gorm:"size:255;not null"`
    Address  string    `gorm:"size:255"`
    Phone    string    `gorm:"size:100"`
    Email    string    `gorm:"size:254"`
    Users    []*User   `gorm:"foreignKey:HospitalID"`
    Patients []*Patient `gorm:"foreignKey:HospitalID"`
}
