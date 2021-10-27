package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/rehandwi03/test-case-backend-majoo/handler/http"
	"github.com/rehandwi03/test-case-backend-majoo/model"
	"github.com/rehandwi03/test-case-backend-majoo/repository"
	"github.com/rehandwi03/test-case-backend-majoo/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	if err := godotenv.Load(); err != nil {
		log.Printf("error when load env: %v", err)
	}

	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUser,
		dbPass, dbName, dbPort,
	)
	db, err := gorm.Open(
		postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		},
	)
	if err != nil {
		log.Panicf("error when connecting to database: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Merchant{}, &model.Outlet{}, &model.Product{}); err != nil {
		log.Printf("error migrating table: %v", err)
	}

	app := fiber.New()
	app.Use(
		logger.New(
			logger.Config{
				Format:     "${pid} ${status} - ${method} ${path}\n",
				TimeFormat: "02-Jan-2006",
				TimeZone:   "Asia/Jakarta",
			},
		),
	)
	apiGroup := app.Group("/api")

	userRepo := repository.NewUserRepository(db)
	merchantRepo := repository.NewMerchantRepository(db)
	outletRepo := repository.NewOutletRepository(db)
	productRepo := repository.NewProductRepository(db)

	userSvc := service.NewUserService(userRepo)
	merchantSvc := service.NewMerchantService(merchantRepo, userRepo)
	outletSvc := service.NewOutletService(outletRepo, merchantRepo)
	productSvc := service.NewProductService(productRepo, outletRepo)
	authRepo := service.NewAuthService(userRepo)

	http.NewUserHandler(apiGroup, userSvc)
	http.NewMerchantHandler(apiGroup, merchantSvc)
	http.NewOutletHandler(apiGroup, outletSvc)
	http.NewProductHandler(apiGroup, productSvc)
	http.NewAuthHandler(apiGroup, authRepo)

	if err := app.Listen(":" + os.Getenv("APP_PORT")); err != nil {
		log.Fatalf("can't start applicaton: %v", err)
	}
}
