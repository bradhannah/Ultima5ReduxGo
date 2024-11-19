package datetime

import (
	"fmt"
	"log"
)

const (
	DaysInMonth    = 28
	MonthsPerYear  = 12
	MinutesPerHour = 60
	HoursPerDay    = 24
)

type TimeOfDay int

const (
	Morning TimeOfDay = iota
	Noon
	Evening
	Midnight
)

type UltimaDate struct {
	Year   uint16
	Month  byte
	Day    byte
	Hour   byte
	Minute byte
}

func (d *UltimaDate) GetDateAsString() string {
	return fmt.Sprintf("%d-%d-%d", d.Month, d.Day, d.Year)
}

func (d *UltimaDate) GetTimeAsString() string {
	hour := d.Hour
	am := true
	if hour >= 12 {
		hour -= 12
		am = false
	}

	if hour == 0 {
		hour = 12
	}

	if am {
		return fmt.Sprintf("%2d:%02dAM", hour, d.Minute)
	}

	return fmt.Sprintf("%2d:%02dPM", hour, d.Minute)
}

func (d *UltimaDate) Advance(nMinutes int) {

	// nMinute that time advancement does not exceed 9 hours (for time-saving assumptions)
	if nMinutes > MinutesPerHour*9 {
		log.Fatal("you cannot advance more than 9 hours at a time")
	}

	// Check if adding minutes moves to a new hour
	if int(d.Minute)+nMinutes > MinutesPerHour-1 {
		nHours := byte(nMinutes / MinutesPerHour)
		nExtraMinutes := nMinutes % MinutesPerHour

		newHour := d.Hour + nHours + 1
		d.Minute = byte((int(d.Minute) + nExtraMinutes) % MinutesPerHour)

		// Check if advancing hours moves to a new day
		if newHour <= HoursPerDay-1 {
			d.Hour = newHour
		} else {
			d.Hour = newHour % HoursPerDay
			// Increment day and handle end of month
			nDay := d.Day + 1
			if nDay > DaysInMonth {
				d.Day = 1
				nMonth := d.Month + 1
				// Increment month and handle end of year
				if nMonth > MonthsPerYear {
					d.Month = 1
					d.Year++
				} else {
					d.Month = nMonth
				}
			} else {
				d.Day = nDay
			}
		}
	} else {
		d.Minute += byte(nMinutes)
	}
}

func (d *UltimaDate) SetTimeOfDay(timeOfDay TimeOfDay) {
	d.Minute = 0

	switch timeOfDay {
	case Morning:
		d.Hour = 5
	case Noon:
		d.Hour = 12
	case Evening:
		d.Hour = 17
	case Midnight:
		d.Hour = 0
	}
}
