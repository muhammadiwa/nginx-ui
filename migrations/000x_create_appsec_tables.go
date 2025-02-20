package migrations

import (
	"github.com/muhammadiwa/nginx-ui/model"
	"gorm.io/gorm"
)

func init() {
	migrations = append(migrations, &AppSecMigration{})
}

type AppSecMigration struct{}

func (m *AppSecMigration) Up(db *gorm.DB) error {
	// Create AppSecConfig table
	if err := db.AutoMigrate(&model.AppSecConfig{}); err != nil {
		return err
	}

	// Create AppSecPolicy table
	if err := db.AutoMigrate(&model.AppSecPolicy{}); err != nil {
		return err
	}

	// Create AppSecRule table
	if err := db.AutoMigrate(&model.AppSecRule{}); err != nil {
		return err
	}

	// Create AppSecEvent table
	if err := db.AutoMigrate(&model.AppSecEvent{}); err != nil {
		return err
	}

	return nil
}

func (m *AppSecMigration) Down(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&model.AppSecEvent{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&model.AppSecRule{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&model.AppSecPolicy{}); err != nil {
		return err
	}
	if err := db.Migrator().DropTable(&model.AppSecConfig{}); err != nil {
		return err
	}
	return nil
}
