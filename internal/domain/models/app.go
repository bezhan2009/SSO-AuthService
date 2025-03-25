package models

type App struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" gorm:"unique not null"`
	Secret string `json:"secret" gorm:"not null"`
}
