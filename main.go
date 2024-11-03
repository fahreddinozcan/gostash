package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gostash/handlers"
	"gostash/models"
)

func main() {
	dsn := "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.Schedule{})
	if err != nil {
		panic("failed to migrate database")
	}

	r := gin.Default()
	scheduleHandler := handlers.NewScheduleHandler(db)

	r.POST("/schedule", scheduleHandler.Create)
	r.GET("/schedule/:id", scheduleHandler.Get)
	r.GET("/schedules", scheduleHandler.List)
	r.DELETE("/schedule/:id", scheduleHandler.Delete)

	fmt.Println("Server is running on port 8080")
	err = r.Run()
	if err != nil {
		panic(err)
	}
}
