package database

import (
	"github.com/google/uuid"
	"github.com/playground-pro-project/playground-pro-api/app/config"
	user "github.com/playground-pro-project/playground-pro-api/features/user/data"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
	"gorm.io/gorm"
)

func initSuperAdmin(db *gorm.DB) error {
	userID, err := uuid.NewUUID()
	if err != nil {
		log.Warn("error while create uuid for admin")
		return nil
	}

	// hashed, err := helper.HashPassword("secret")
	hashed, err := helper.HashPassword(config.ADMINPASSWORD)
	if err != nil {
		log.Warn("error while hashing password admin")
		return nil
	}

	admin := user.User{
		UserID:         userID.String(),
		Fullname:       "admin",
		Email:          "admin@gmail.com",
		Phone:          "081235288543",
		Password:       hashed,
		Bio:            "superadmin",
		Address:        "bikini-button",
		Role:           "admin",
		ProfilePicture: "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png",
	}

	var count int64
	db.Table("users").Where("role = 'admin'").Count(&count)
	if count > 0 {
		log.Warn("super admin already exists")
		return nil
	}

	result := db.Create(&admin)
	if result.Error != nil {
		log.Error("failed to create super admin")
		return result.Error
	}

	log.Info("super admin created successfully")
	return nil
}
