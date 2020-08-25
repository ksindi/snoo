package main

import (
	"time"
)

// Session represents a SNOO session
type Session struct {
	AsleepDuration   int
	EndTime          time.Time
	ID               string
	SoothingDuration int
	StartTime        time.Time

	levels []Level
}

// NewSession initializes a new `Session`
func NewSession(sessionID string, levels []Level) Session {
	// https://stackoverflow.com/questions/25065055/what-is-the-maximum-time-time-in-go
	MaxTime := time.Unix(1<<63-62135596801, 999999999)

	var asleepDuration, soothingDuration int
	minStartTime := MaxTime
	maxEndTime := time.Time{} // initialize to min time

	for _, level := range levels {
		if level.StartTime.Before(minStartTime) {
			minStartTime = level.StartTime.Time
		}

		if level.EndTime().After(maxEndTime) {
			maxEndTime = level.EndTime()
		}

		if level.Type == "asleep" {
			asleepDuration += level.StateDuration
		}

		if level.Type == "soothing" {
			soothingDuration += level.StateDuration
		}

	}

	return Session{
		AsleepDuration:   asleepDuration,
		EndTime:          maxEndTime,
		ID:               sessionID,
		SoothingDuration: soothingDuration,
		StartTime:        minStartTime,

		levels: levels,
	}
}

// EndTime calculates the end time of a level
func (l *Level) EndTime() time.Time {
	return l.StartTime.Add(time.Second * time.Duration(l.StateDuration))
}

// TotalDuration calculates the total session duration
func (s *Session) TotalDuration() int {
	return s.SoothingDuration + s.AsleepDuration
}
