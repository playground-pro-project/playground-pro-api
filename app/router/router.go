package router

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	rsd "github.com/playground-pro-project/playground-pro-api/features/reservation/data"
	rsh "github.com/playground-pro-project/playground-pro-api/features/reservation/handler"
	rss "github.com/playground-pro-project/playground-pro-api/features/reservation/service"
	rd "github.com/playground-pro-project/playground-pro-api/features/review/data"
	rh "github.com/playground-pro-project/playground-pro-api/features/review/handler"
	rs "github.com/playground-pro-project/playground-pro-api/features/review/service"
	ud "github.com/playground-pro-project/playground-pro-api/features/user/data"
	uh "github.com/playground-pro-project/playground-pro-api/features/user/handler"
	us "github.com/playground-pro-project/playground-pro-api/features/user/service"
	vd "github.com/playground-pro-project/playground-pro-api/features/venue/data"
	vh "github.com/playground-pro-project/playground-pro-api/features/venue/handler"
	vs "github.com/playground-pro-project/playground-pro-api/features/venue/service"
	paymentgateway "github.com/playground-pro-project/playground-pro-api/utils/payment_gateway"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, e *echo.Echo) {
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))

	initUserRouter(db, e)
	initVenueRouter(db, e)
	initReservationRouter(db, e)
}

func initUserRouter(db *gorm.DB, e *echo.Echo) {
	userData := ud.New(db)
	validate := validator.New()
	userService := us.New(userData, validate)
	userHandler := uh.New(userService)

	e.POST("/register", userHandler.Register())
	e.POST("/login", userHandler.Login())
	e.POST("/resend-otp", userHandler.ReSendOTP())
	e.POST("/validation", userHandler.ValidateOTP())
	e.GET("/users", userHandler.GetUserProfile(), middlewares.JWTMiddleware())
	e.PUT("/users", userHandler.UpdateUserProfile(), middlewares.JWTMiddleware())
	e.PUT("/users/password", userHandler.UpdatePassword(), middlewares.JWTMiddleware())
	e.DELETE("/users", userHandler.DeleteUser(), middlewares.JWTMiddleware())
	e.POST("/users/upgrade", userHandler.UploadOwnerFile(), middlewares.JWTMiddleware())
	e.PUT("/users/profile-picture", userHandler.UploadProfilePicture(), middlewares.JWTMiddleware())
	e.DELETE("/users/profile-picture", userHandler.RemoveProfilePicture(), middlewares.JWTMiddleware())
}

func initVenueRouter(db *gorm.DB, e *echo.Echo) {
	venueData := vd.New(db)
	venueService := vs.New(venueData)
	venueHandler := vh.New(venueService)

	reviewData := rd.New(db)
	reviewService := rs.New(reviewData)
	reviewHandler := rh.New(reviewService)

	e.POST("/venues", venueHandler.RegisterVenue(), middlewares.JWTMiddleware())
	e.GET("/venues", venueHandler.SearchVenues())
	e.GET("/venues/:venue_id", venueHandler.SelectVenue(), middlewares.JWTMiddleware())
	e.GET("/venues/:venue_id/availability", venueHandler.VenueAvailability(), middlewares.JWTMiddleware())
	e.PUT("/venues/:venue_id", venueHandler.EditVenue(), middlewares.JWTMiddleware())
	e.DELETE("/venues/:venue_id", venueHandler.UnregisterVenue(), middlewares.JWTMiddleware())
	e.POST("/venues/:venue_id/reviews", reviewHandler.CreateReview, middlewares.JWTMiddleware())
	e.GET("/venues/:venue_id/reviews", reviewHandler.GetAllReview, middlewares.JWTMiddleware())
	e.DELETE("/reviews/:review_id", reviewHandler.DeleteReview)
	e.POST("/venues/:venue_id/images", venueHandler.CreateVenueImage(), middlewares.JWTMiddleware())
	e.DELETE("/venues/:venue_id/images/:image_id", venueHandler.DeleteVenueImage(), middlewares.JWTMiddleware())
	e.GET("/venues/:venue_id/images", venueHandler.GetAllVenueImage(), middlewares.JWTMiddleware())
}

func initReservationRouter(db *gorm.DB, e *echo.Echo) {
	reservationData := rsd.New(db)
	refund := &paymentgateway.MyRefund{}
	reservationService := rss.New(reservationData, refund)
	reservationHandler := rsh.New(reservationService)

	e.POST("/reservations", reservationHandler.MakeReservation(), middlewares.JWTMiddleware())
	e.POST("/reservations/status", reservationHandler.ReservationStatus())
	e.GET("/users/reservations", reservationHandler.ReservationHistory(), middlewares.JWTMiddleware())
}
