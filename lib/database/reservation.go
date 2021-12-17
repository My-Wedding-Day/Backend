package database

import (
	"alta-wedding/config"
	"alta-wedding/models"
)

// Fungsi untuk membuat data booking
func CreateReservation(reservation *models.Reservation) (*models.Reservation, error) {
	// CHECK DATABASE ALREADY RESERVE OR NOT
	tx := config.DB.Where("date = ? AND package_id = ?", reservation.Date, reservation.Package_ID).Find(&models.Reservation{})
	// IF ERROR
	if tx.Error != nil {
		return nil, tx.Error
	}
	// IF DATA ALREADY
	if tx.RowsAffected > 0 {
		return nil, nil
	}
	// IF DIDN'T RESERVE CHECK
	err := config.DB.Create(&reservation).Error
	if err != nil {
		return nil, err
	}
	// SUCCESS RESERVE
	return reservation, nil
}
