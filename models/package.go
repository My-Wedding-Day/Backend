package models

import (
	"time"

	"gorm.io/gorm"
)

type Package struct {
	ID          int    `gorm:"primarykey"`
	PackageName string `gorm:"type:varchar(255);NOT NULL" json:"packagename" form:"packagename"`
	Price       int    `gorm:"type:int(100);NOT NULL" json:"price" form:"price"`
	Pax         int    `gorm:"type:int(100);NOT NULL" json:"pax" form:"pax"`
	PackageDesc string `gorm:"type:varchar(1000);NOT NULL" json:"packagedesc" form:"packagedesc"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type PostRequestBodyPackage struct {
	PackageName string `json:"packagename" form:"packagename"`
	Price       int    `json:"price" form:"price"`
	Pax         int    `json:"pax" form:"pax"`
	PackageDesc string `json:"packagedesc" form:"packagedesc"`
	UrlPhoto    string `json:"urlphoto" form:"urlphoto"`
}
