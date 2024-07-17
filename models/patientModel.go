package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FirstName  string    `gorm:"size:255;not null"`
	LastName   string    `gorm:"size:255;not null"`
	Email      string    `gorm:"size:254;unique;not null"`
	Phone      string    `gorm:"size:255"`
	BirthDate  time.Time `gorm:"type:date;not null"`
	HospitalID uint
	Hospital   *Hospital   `gorm:"foreignKey:HospitalID"`
	Address    *Address    `gorm:"foreignKey:PatientID;constraint:OnDelete:CASCADE"`
	CellTests  []*CellTest `gorm:"foreignKey:PatientID"`
}

type Address struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Street    string    `gorm:"size:255;not null"`
	City      string    `gorm:"size:255;not null"`
	PatientID uuid.UUID `gorm:"type:uuid;unique;not null"`
	Patient   *Patient  `gorm:"foreignKey:PatientID;constraint:OnDelete:CASCADE"`
}
