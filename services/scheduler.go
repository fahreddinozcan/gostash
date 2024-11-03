package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"gostash/models"
	"net/http"
	"sync"
	"time"
)

type Scheduler struct {
	db           *gorm.DB
	schedules    map[uint64]*ScheduleJob
	schedulesMux sync.RWMutex
}

type ScheduleJob struct {
	Schedule   *models.Schedule
	Cron       *cron.Cron
	Context    context.Context
	CancelFunc context.CancelFunc
}

func NewScheduler(db *gorm.DB) *Scheduler {
	return &Scheduler{
		db:        db,
		schedules: make(map[uint64]*ScheduleJob),
	}
}

func (s *Scheduler) StartSchedule(schedule *models.Schedule) error {
	s.schedulesMux.Lock()
	defer s.schedulesMux.Unlock()

	_, exists := s.schedules[schedule.Id]
	if exists {
		return fmt.Errorf("schedule with id %d already running", schedule.Id)
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc(schedule.Cron, func() {
		s.executeJob(ctx, schedule)
	})
	if err != nil {
		cancel()
		return fmt.Errorf("invalid cron expr: %v", err)
	}

	s.schedules[schedule.Id] = &ScheduleJob{
		Schedule:   schedule,
		Cron:       c,
		Context:    ctx,
		CancelFunc: cancel,
	}

	c.Start()
	return nil
}

func (s *Scheduler) executeJob(ctx context.Context, schedule *models.Schedule) {
	done := make(chan bool)

	event := &models.Event{
		ScheduleID: schedule.Id,
		State:      models.EventStateCreated,
		CreatedAt:  time.Now(),
	}
	s.db.Create(event)

	go func() {
		event.State = models.EventStatusRunning
		s.db.Save(event)

		client := &http.Client{}

		var headers map[string]string
		err := json.Unmarshal([]byte(schedule.Headers), &headers)
		if err != nil {
			s.logEvent(schedule.Id, "failed", fmt.Sprintf("Error parsing headers: %v", err))
			done <- true
			return
		}

		req, err := http.NewRequest(schedule.Method, schedule.Endpoint, bytes.NewBufferString(schedule.Body))
		if err != nil {
			s.logEvent(schedule.Id, "failed", fmt.Sprintf("Error creating request: %v", err))
			done <- true
			return
		}

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		req = req.WithContext(ctx)
		resp, err := client.Do(req)
		if err != nil {
			s.logEvent(schedule.Id, "failed", fmt.Sprintf("Error executing request: %v", err))
			done <- true
			return
		}
		defer resp.Body.Close()

		s.db.Model(schedule).Update("last_run_at", time.Now())

		status := "success"
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			status = "failed"
		}

		s.logEvent(schedule.Id, status, fmt.Sprintf("Scheduled job with id %d with status %s", schedule.Id, status))
		done <- true
	}()

	select {
	case <-ctx.Done():
		s.logEvent(schedule.Id, "cancelled", "Schedule stopped")
		return
	case <-done:
		return
	}
}

func (s *Scheduler) logEvent(id uint64, status string, response string) {
	event := models.Event{
		ScheduleID: id,
		Status:     status,
		Response:   response,
	}

	s.db.Create(&event)
}
