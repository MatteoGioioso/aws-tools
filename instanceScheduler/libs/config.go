package libs

import (
	"errors"
	"log"
	"time"
)

type pattern struct {
	hoursOn map[int]bool
	daysOn  map[string]bool
}

type Period struct {
	Pattern string `json:"pattern"`
}

type Resource struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

type Resources map[string]Resource

type Report struct {
	SendReport bool `json:"sendReport"`
	Hour       int  `json:"hour"`
}

type SchedulerConfig struct {
	Report    Report    `json:"report"`
	Period    Period    `json:"period"`
	TimeZone  string    `json:"timeZone"`
	Resources Resources `json:"resources"`
}

type SchedulerConfigClient struct {
	now    func() time.Time
	Config *SchedulerConfig `json:"config"`
}

func NewSchedulerConfigClient(config *SchedulerConfig) *SchedulerConfigClient {
	return &SchedulerConfigClient{Config: config, now: time.Now}
}

func (s SchedulerConfigClient) ShouldWakeup() (bool, error) {
	tz, err := s.getCurrentTimeFromTZ()
	if err != nil {
		return false, err
	}

	p := s.Config.Period.Pattern
	patternConfig, ok := allowedPatterns[p]
	if !ok {
		return false, errors.New("invalid pattern")
	}

	log.Printf("Current weekday: %v", tz.Weekday().String())
	log.Printf("Current hour: %v", tz.Hour())
	log.Printf("Pattern: %v", p)

	if p == permanentOn {
		return true, nil
	}

	if s.isWakeupDay(patternConfig, tz) && s.isWakeupHour(patternConfig, tz) {
		return true, nil
	}

	return false, nil
}

func (s SchedulerConfigClient) isWakeupDay(pattern pattern, time time.Time) bool {
	return pattern.daysOn[time.Weekday().String()]
}

func (s SchedulerConfigClient) isWakeupHour(pattern pattern, time time.Time) bool {
	return pattern.hoursOn[time.Hour()]
}

func (s SchedulerConfigClient) getCurrentTimeFromTZ() (time.Time, error) {
	location, err := time.LoadLocation(s.Config.TimeZone)
	if err != nil {
		return time.Time{}, err
	}

	timeWithTimeZone := s.now().In(location)
	return timeWithTimeZone, nil
}

func (s SchedulerConfigClient) ShouldSendReport() (bool, error) {
	if !s.Config.Report.SendReport {
		return false, nil
	}

	tz, err := s.getCurrentTimeFromTZ()
	if err != nil {
		return false, err
	}

	if tz.Hour() == s.Config.Report.Hour {
		return true, nil
	}

	return false, nil
}
