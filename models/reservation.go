package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	ID             int    `gorm:"primarykey;AUTO_INCREMENT"`
	Package_ID     int    `json:"package_id" form:"package_id"`
	User_ID        int    `json:"user_id" form:"user_id"`
	Date           string `gorm:"type:varchar(255)" json:"date" form:"date"`
	Additional     string `gorm:"type:varchar(255)" json:"additional" form:"additional"`
	Total_Pax      int    `gorm:"type:int" json:"total_pax" form:"total_pax"`
	Status_Order   string `gorm:"type:varchar(50); default:waiting" json:"status_order" form:"status_order"`
	Status_Payment string `gorm:"type:varchar(50); default:unpaid" json:"status_payment" form:"status_payment"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
