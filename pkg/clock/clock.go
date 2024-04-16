package clock

import "time"

type ClockI interface {
	Now() time.Time
}

type RealClock struct{}

func (c RealClock) Now() time.Time { return time.Now() }

type DumbClock struct{}

func (c DumbClock) Now() time.Time { return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) }

var instance ClockI = RealClock{}

func InitClock(isDumb bool) {
	if isDumb {
		instance = DumbClock{}
	}
}

func Now() time.Time { return instance.Now() }
