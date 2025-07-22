package utils

import (
	"strings"
	"time"
	_ "time/tzdata"

	"log"
)

func GetTimeLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Error loading location: %v", err)
	}

	return loc
}

func GetCurrentTime() time.Time {

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Fatalf("Error loading location: %v", err)
	}

	currentTime, err := time.Parse("2006-01-02 15:04:05", time.Now().In(loc).Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Fatalf("Error parsing time: %v", err)
	}

	return currentTime
}

func ParseTime(date string) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	formattedDate := strings.ReplaceAll(strings.ReplaceAll(date, "T", " "), "Z", "")
	parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", formattedDate, loc)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}
