package models

import (
	"encoding/json"
	"time"
)

type EventState string

const (
	EventStateCreated   EventState = "created"
	EventStatusRunning  EventState = "running"
	EventStateSuccess   EventState = "success"
	EventStateFailed    EventState = "failed"
	EventStatusPending  EventState = "pending"
	EventStatusCanceled EventState = "canceled"
)

type Event struct {
	ID         uint64     `json:"id" gorm:"primaryKey"`
	ScheduleID uint64     `json:"schedule_id" gorm:"not null"`
	State      EventState `json:"state" gorm:"not null"`
	Response   string     `json:"response" gorm:"type:text"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
	Error      string     `json:"error" gorm:"type:text"`
	Schedule   Schedule   `json:"-" gorm:"foreignKey:ScheduleID"`
}

type ResponseData struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

func (e *Event) SetResponse(resp *ResponseData) error {
	jsonData, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	e.Response = string(jsonData)
	return nil
}

func (e *Event) GetResponse() (*ResponseData, error) {
	if e.Response == "" {
		return nil, nil
	}
	var resp ResponseData
	err := json.Unmarshal([]byte(e.Response), &resp)
	return &resp, err
}
