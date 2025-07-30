package clock

import (
	"github.com/google/wire"
	"time"
)

type TimeSource interface {
	Now() time.Time
}

type RealTimeSource struct{}

func NewRealTimeSource() *RealTimeSource {
	return &RealTimeSource{}
}

func (r RealTimeSource) Now() time.Time {
	return time.Now()
}

type MockTimeSource struct {
	initialTime time.Time
	currentTime time.Time
}

func NewMockTimeSource(initialTime time.Time) *MockTimeSource {
	return &MockTimeSource{initialTime: initialTime, currentTime: initialTime}
}

func (m *MockTimeSource) Now() time.Time {
	return m.currentTime
}

func (m *MockTimeSource) SetTime(t time.Time) {
	m.currentTime = t
}

func (m *MockTimeSource) Reset() {
	m.currentTime = m.initialTime
}

var RealClockProvider = wire.NewSet(wire.Bind(new(TimeSource), new(*RealTimeSource)), NewRealTimeSource)
