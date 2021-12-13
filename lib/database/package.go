package database

import (
	"alta-wedding/config"
	"alta-wedding/models"
)

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
