package util

import "time"

func LocationOrDefault(timezone string) *time.Location {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Local
	}
	return loc
}

func TimeNow(timezone string) time.Time {
	return time.Now().In(LocationOrDefault(timezone))
}

func CalculateNextBlastDate(offset time.Time, targetDay time.Weekday, targetTime time.Duration) time.Time {
	today := time.Date(offset.Year(), offset.Month(), offset.Day(), 0, 0, 0, 0, offset.Location())
	daysOut := time.Duration(int(time.Saturday) - int(offset.Weekday()) + int(targetDay) + 1)
	nextBlastDate := today.Add(daysOut * 24 * time.Hour).Add(targetTime)
	return nextBlastDate
}
