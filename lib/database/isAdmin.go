package database

import (
	"alta-wedding/config"
	"alta-wedding/models"
)

func GetInvoiceAdmin() ([]models.PaymentInvoice, error) {
	var paymentinvoice []models.PaymentInvoice
	query := config.DB.Table("payments").Select("payments.id, payments.reservation_id, payments.url_photo, reservations.user_id, reservations.total_pax, reservations.status_payment, packages.price, reservations.total_price").
		Joins("join reservations on reservations.id = payments.reservation_id").Joins("join packages on packages.id = reservations.package_id").
		Where("reservations.status_payment = 'unpaid' AND reservations.status_order = 'accepted' AND reservations.deleted_at is NULL").Find(&paymentinvoice)
	if query.Error != nil {
		return nil, query.Error
	}
	return paymentinvoice, nil
}
