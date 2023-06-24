package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	rd "github.com/playground-pro-project/playground-pro-api/features/reservation/data"
	rh "github.com/playground-pro-project/playground-pro-api/features/reservation/handler"
	rs "github.com/playground-pro-project/playground-pro-api/features/reservation/service"
	ud "github.com/playground-pro-project/playground-pro-api/features/user/data"
	uh "github.com/playground-pro-project/playground-pro-api/features/user/handler"
	us "github.com/playground-pro-project/playground-pro-api/features/user/service"
	vd "github.com/playground-pro-project/playground-pro-api/features/venue/data"
	vh "github.com/playground-pro-project/playground-pro-api/features/venue/handler"
	vs "github.com/playground-pro-project/playground-pro-api/features/venue/service"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, e *echo.Echo) {
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))

	initUserRouter(db, e)
	initVanueRouter(db, e)
	initReservationRouter(db, e)
}

func initUserRouter(db *gorm.DB, e *echo.Echo) {
	userData := ud.New(db)
	userService := us.New(userData)
	userHandler := uh.New(userService)

	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	usersGroup := e.Group("/users")
	{
		usersGroup.GET("", userHandler.GetUserProfile, middlewares.JWTMiddleware())
		usersGroup.PUT("", userHandler.UpdateUserProfile, middlewares.JWTMiddleware())
		usersGroup.PUT("/password", userHandler.UpdatePassword, middlewares.JWTMiddleware())
		usersGroup.DELETE("", userHandler.DeleteUser, middlewares.JWTMiddleware())
		usersGroup.PUT("", userHandler.UploadProfilePicture, middlewares.JWTMiddleware())
		usersGroup.PUT("", userHandler.RemoveProfilePicture, middlewares.JWTMiddleware())
	}
}

func initVanueRouter(db *gorm.DB, e *echo.Echo) {
	vanueData := vd.New(db)
	vanueService := vs.New(vanueData)
	vanueHandler := vh.New(vanueService)

	e.GET("/venues", vanueHandler.SearchVenue())
}

func initReservationRouter(db *gorm.DB, e *echo.Echo) {
	reservationData := rd.New(db)
	reservationService := rs.New(reservationData)
	reservationHandler := rh.New(reservationService)

	e.POST("/reservations", reservationHandler.MakeReservation(), middlewares.JWTMiddleware())
}
