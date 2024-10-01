package game_state

import "fmt"

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
	if d.Hour >= 12 {
		return fmt.Sprintf("%2d:%02dPM", d.Hour-12, d.Minute)
	}

	return fmt.Sprintf("%2d:%02dAM", d.Hour, d.Minute)
}
