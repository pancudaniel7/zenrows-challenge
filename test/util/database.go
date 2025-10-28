package util

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewTestDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.database"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.sslmode"),
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
