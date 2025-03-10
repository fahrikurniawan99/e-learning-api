package config

import (
	"elearning_api/model"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() *gorm.DB {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Ambil variabel dari .env
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		Env.DbHost,
		Env.DbUser,
		Env.DbPassword,
		Env.DbName,
		Env.DbPort,
	)

	// Koneksi ke database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	fmt.Println("Database connected!")

	if Env.AppEnv != "production" {
		err = DB.AutoMigrate(
			&model.User{},
			&model.RegistrationMethod{},
			&model.AccessToken{},
			&model.RefreshToken{},
			&model.VerificationToken{},
			&model.AccessTokenLogs{},
			&model.RefreshTokenLogs{},
			// Tambahkan model lainnya yang perlu di-migrate
		)
		if err != nil {
			log.Fatal("Failed to migrate database:", err)
		}
		fmt.Println("Database migrated!")
	}

	return db
}
