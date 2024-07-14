package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CellTest struct {
	gorm.Model
	ID              uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title           string           `gorm:"size:255;not null"`
	Description     string           `gorm:"type:text"`
	DetectionStatus string           `gorm:"size:255;not null"`
	PatientID       uuid.UUID        `gorm:"type:uuid;not null"`
	Patient         *Patient         `gorm:"foreignKey:PatientID"`
	Results         []*Result        `gorm:"foreignKey:CellTestID"`
	CellTestImages  []*CellTestImage `gorm:"foreignKey:CellTestID"`
}

type Result struct {
	gorm.Model
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Description  string         `gorm:"type:text"`
	CellTestID   uuid.UUID      `gorm:"type:uuid;not null"`
	CellTest     *CellTest      `gorm:"foreignKey:CellTestID"`
	ResultImages []*ResultImage `gorm:"foreignKey:ResultID"`
}

type CellTestImage struct {
	gorm.Model
	Image      string    `gorm:"size:100;not null"`
	CellTestID uuid.UUID `gorm:"type:uuid;not null"`
	CellTest   *CellTest `gorm:"foreignKey:CellTestID"`
}

type ResultImage struct {
	gorm.Model
	Image    string    `gorm:"size:100;not null"`
	ResultID uuid.UUID `gorm:"type:uuid;not null"`
	Result   *Result   `gorm:"foreignKey:ResultID"`
}
