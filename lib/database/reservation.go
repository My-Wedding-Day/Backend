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

// Fungsi untuk mendapatkan reservasi by reservasi id
func GetReservation(id int) (interface{}, error) {
	var reservation models.GetReservationRespon
	query := config.DB.Table("reservations").Select("reservations.id, reservations.package_id, packages.package_name, organizers.wo_name, reservations.date, reservations.additional, reservations.total_pax, reservations.status_order, reservations.status_payment").
		Joins("join packages on packages.id = reservations.package_id").Joins("join organizers on organizers.id = packages.organizer_id").
		Where("reservations.user_id = ? AND reservations.deleted_at is NULL", id).Find(&reservation)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected < 1 {
		return nil, nil
	}
	return reservation, nil
}
