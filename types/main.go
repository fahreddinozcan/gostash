package types

type ScheduleRequest struct {
	Name     string            `json:"name"`
	Endpoint string            `json:"endpoint"`
	Cron     string            `json:"cron"`
	Body     string            `json:"body"`
	Headers  map[string]string `json:"headers"`
}
