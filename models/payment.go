package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID             int    `gorm:"primarykey"`
	Reservation_ID int    `gorm:"type:int(100);NOT NULL; unique " json:"reservationid" form:"reservationid"`
	User_ID        int    `json:"userid" form:"userid"`
	UrlPhoto       string `gorm:"type:varchar(255);NOT NULL" json:"invoice" form:"invoice"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
