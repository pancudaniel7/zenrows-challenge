package infra

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDatabase() *gorm.DB {
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	name := viper.GetString("database.database")
	user := viper.GetString("database.user")
	pass := viper.GetString("database.password")
	ssl := viper.GetString("database.sslmode")

	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", host, port, name, user, pass, ssl)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	return db
}
