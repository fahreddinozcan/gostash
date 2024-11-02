package models

import (
	"encoding/json"
	"time"
)

type Schedule struct {
	Id        uint64    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Endpoint  string    `json:"endpoint" gorm:"not null"`
	Cron      string    `json:"cron" gorm:"not null"`
	Body      string    `json:"body"`
	Headers   string    `json:"headers" gorm:"type:jsonb;default:'{}'"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	LastRunAt time.Time `json:"last_run_at"`
}

func (s *Schedule) SetHeaders(headers map[string]string) error {
	if headers == nil {
		s.Headers = "{}"
		return nil
	}

	jsonBytes, err := json.Marshal(headers)
	if err != nil {
		return err
	}
	s.Headers = string(jsonBytes)
	return nil
}

func (s *Schedule) GetHeaders() (map[string]string, error) {
	if s.Headers == "" {
		return make(map[string]string), nil
	}

	var headers map[string]string

	err := json.Unmarshal([]byte(s.Headers), &headers)
	if err != nil {
		return nil, err
	}
	return headers, nil
}
