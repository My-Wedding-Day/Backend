package database

import (
	"alta-wedding/config"
	"alta-wedding/models"
)

type GetPackageStruct struct {
	ID           int
	Organizer_ID int
	Wo_Name      string
	City         string
	Address      string
	PackageName  string
	Price        int
	Pax          int
	PackageDesc  string
	UrlPhoto     string
}

type GetPackageAllStruct struct {
	ID           int
	Organizer_ID int
	PackageName  string
	Price        int
	Pax          int
	PackageDesc  string
	UrlPhoto     string
}

func InsertPackage(Package models.Package) (models.Package, error) {
	tx := config.DB.Save(&Package)
	if tx.Error != nil {
		return Package, tx.Error
	}
	return Package, nil
}

func GetPackageByName(PackageName string) (int64, error) {
	tx := config.DB.Where("package_name = ?", PackageName).Find(&models.Package{})
	if tx.Error != nil {
		return 0, tx.Error
	}
	if tx.RowsAffected > 0 {
		return tx.RowsAffected, nil
	}
	return 0, nil
}

// Fungsi untuk mendapatkan seluruh data packages
func GetPackages() (interface{}, error) {
	var paket []GetPackageAllStruct

	query := config.DB.Table("packages").Select(
		"photos.url_photo, packages.package_desc, packages.pax, packages.price, packages.package_name, packages.organizer_id, packages.id").Joins(
		"join photos on packages.id = photos.package_id").Where(
		"packages.deleted_at is NULL").Find(&paket)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	return paket, nil
}

// Fungsi untuk mendapatkan seluruh data packages by id organizer
func GetPackagesByToken(id int) (interface{}, error) {
	var paket []GetPackageAllStruct

	query := config.DB.Table("packages").Select(
		"photos.url_photo, packages.package_desc, packages.pax, packages.price, packages.package_name, packages.organizer_id, packages.id").Joins(
		"join photos on packages.id = photos.package_id").Where(
		"packages.organizer_id = ? AND packages.deleted_at is NULL", id).Find(&paket)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	return paket, nil
}

// Fungsi untuk mendapatkan seluruh data packages by id
func GetPackagesByID(id int) (interface{}, error) {
	var paket []GetPackageStruct

	query := config.DB.Table("packages").Select(
		"organizers.wo_name, organizers.city, organizers.address, photos.url_photo, packages.package_desc, packages.pax, packages.price, packages.package_name, packages.organizer_id, packages.id").Joins(
		"join photos on packages.id = photos.package_id").Joins(
		"join organizers on organizers.id = packages.organizer_id").Where(
		"packages.id = ? AND packages.deleted_at is NULL", id).Find(&paket)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return 0, query.Error
	}
	return paket, nil
}
