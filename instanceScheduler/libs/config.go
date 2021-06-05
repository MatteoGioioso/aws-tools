package libs

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

type pattern struct {
	hoursOn map[int]bool
	daysOn  map[string]bool
}

var allowedTypes = [4]string{"EC2", "RDS", "ASG", "Fargate"}

type Period struct {
	Pattern string `json:"pattern"`
}

type Resource struct {
	Type       string `json:"type"`
	Identifier string `json:"identifier"`
}

type Resources map[string]Resource

type SchedulerConfig struct {
	Period    Period    `json:"period"`
	TimeZone  string    `json:"timeZone"`
	Resources Resources `json:"resources"`
	now       func() time.Time
}

func (s SchedulerConfig) ParseTime(resourceIdentifier string) ([]int, error) {
	period := s.Period.Pattern
	hours := strings.Split(period, "-")

	ints := make([]int, 0)
	for _, hour := range hours {
		atoi, err := strconv.Atoi(hour)
		if err != nil {
			return nil, err
		}

		ints = append(ints, atoi)
	}

	return ints, nil
}

func (s SchedulerConfig) shouldWakeup() (bool, error) {
	tz, err := s.getCurrentTimeFromTZ()
	if err != nil {
		return false, err
	}

	log.Printf("Current weekday: %v", tz.Weekday().String())
	log.Printf("Current hour: %v", tz.Hour())

	p := s.Period.Pattern
	patternConfig, ok := allowedPatterns[p]
	if !ok {
		return false, errors.New("invalid pattern")
	}

	if s.isWakeupDay(patternConfig, tz) && s.isWakeupHour(patternConfig, tz) {
		return true, nil
	}

	return false, nil
}

func (s SchedulerConfig) isWakeupDay(pattern pattern, time time.Time) bool {
	return pattern.daysOn[time.Weekday().String()]
}

func (s SchedulerConfig) isWakeupHour(pattern pattern, time time.Time) bool {
	return pattern.hoursOn[time.Hour()]
}

func (s SchedulerConfig) getCurrentTimeFromTZ() (time.Time, error) {
	location, err := time.LoadLocation(s.TimeZone)
	if err != nil {
		return time.Time{}, err
	}

	timeWithTimeZone := s.now().In(location)
	return timeWithTimeZone, nil
}
