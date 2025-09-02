package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeSourceImplementations(t *testing.T) {
	assert.Implements(t, new(TimeSource), new(RealTimeSource))
	assert.Implements(t, new(TimeSource), new(MockTimeSource))
}

func TestMockTimeControl(t *testing.T) {
	initial := time.Date(2025, 4, 6, 2, 0, 0, 0, time.UTC)
	mock := NewMockTimeSource(initial)

	assert.Equal(t, initial, mock.Now())

	newTime := initial.Add(2 * time.Hour)
	mock.SetTime(newTime)

	assert.Equal(t, newTime, mock.Now())
}

func TestRealTimeSource(t *testing.T) {
	clock := NewRealTimeSource()
	before := time.Now()
	now := clock.Now()
	time.Sleep(time.Millisecond * 10)
	after := time.Now()

	assert.True(t, now.After(before))
	assert.True(t, now.Before(after))
}

func TestMockZeroTime(t *testing.T) {
	zeroTime := time.Time{}
	mock := NewMockTimeSource(zeroTime)

	assert.True(t, mock.Now().IsZero())

	mock.SetTime(zeroTime.Add(1 * time.Nanosecond))

	assert.False(t, mock.Now().Equal(zeroTime))
}
