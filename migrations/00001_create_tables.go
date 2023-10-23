package migrations

import (
	"fmt"

	"github.com/go-bolo/bolo"
	"gorm.io/gorm"
)

func Get00001Migration() *bolo.Migration {
	return &bolo.Migration{
		Name: "create_tables",
		Up: func(app bolo.App) error {
			db := app.GetDB()
			return db.Transaction(func(tx *gorm.DB) error {

				err := tx.Exec(`CREATE TABLE IF NOT EXISTS system_settings (
					` + "`key`" + ` varchar(255) NOT NULL,
					` + "`value`" + ` varchar(255) NOT NULL,
					PRIMARY KEY (` + "`key`" + `)
				) `).Error
				if err != nil {
					return fmt.Errorf("failed to create system_settings table: %w", err)
				}

				return nil
			})
		},
		Down: func(app bolo.App) error {
			return nil
		},
	}
}
