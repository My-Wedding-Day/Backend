package database

import (
	"alta-wedding/config"
	"alta-wedding/models"
)

// Fungsi untuk membuat data booking
func CreateReservation(reservation *models.Reservation) (*models.Reservation, error) {
	tx := config.DB.Create(&reservation)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return reservation, nil
}
