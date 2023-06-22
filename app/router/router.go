package router

import (
	"github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	ud "github.com/playground-pro-project/playground-pro-api/features/user/data"
	uh "github.com/playground-pro-project/playground-pro-api/features/user/handler"
	us "github.com/playground-pro-project/playground-pro-api/features/user/service"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, e *echo.Echo) {
	userData := ud.New(db)
	userservice := us.New(userData)
	userHandler := uh.New(userservice)

	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	usersGroup := e.Group("/users")
	{
		usersGroup.GET("", userHandler.GetUserProfile, middlewares.JWTMiddlewareFunc())
		usersGroup.PUT("", userHandler.UpdateUserProfile, middlewares.JWTMiddlewareFunc())
		usersGroup.PUT("/password", userHandler.UpdatePassword, middlewares.JWTMiddlewareFunc())
		usersGroup.DELETE("", userHandler.DeleteUser, middlewares.JWTMiddlewareFunc())
		usersGroup.PUT("", userHandler.RemoveProfilePicture, middlewares.JWTMiddlewareFunc())
	}
}
