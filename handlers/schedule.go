package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gostash/models"
	"gostash/types"
	"net/http"
)

type ScheduleHandler struct {
	db *gorm.DB
}

func NewScheduleHandler(db *gorm.DB) *ScheduleHandler {
	return &ScheduleHandler{db: db}
}

func (h *ScheduleHandler) Create(ctx *gin.Context) {
	var req types.ScheduleRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	schedule := models.Schedule{
		Name:     req.Name,
		Endpoint: req.Endpoint,
		Cron:     req.Cron,
		Body:     req.Body,
	}

	err = schedule.SetHeaders(req.Headers)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "invalid headers format",
		})
	}

	err = h.validateSchedule(&schedule)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	err = h.db.Create(&schedule).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create new schedule",
		})
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":     http.StatusCreated,
		"scheduleId": schedule.Id,
	})
}

func (h *ScheduleHandler) Get(ctx *gin.Context) {
	id := ctx.Param("id")

	var schedule models.Schedule
	err := h.db.First(&schedule, id).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "schedule not found",
		})
	}
}

func (h *ScheduleHandler) List(ctx *gin.Context) {
	var schedules []models.Schedule

	err := h.db.Find(&schedules).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to fetch schedules",
		})
	}

	ctx.JSON(http.StatusOK, schedules)
}

func (h *ScheduleHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.db.Delete(&models.Schedule{}, id).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "failed to delete schedule",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "schedule deleted",
	})
}

func (h *ScheduleHandler) validateSchedule(s *models.Schedule) error {
	if s.Name == "" {
		return errors.New("name is required")
	}
	if s.Endpoint == "" {
		return errors.New("endpoint is required")
	}
	if s.Cron == "" {
		return errors.New("cron is required")
	}
	return nil
}
