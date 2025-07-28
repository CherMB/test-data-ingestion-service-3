package helper

import (
	"time"

	"github.com/calculi-corp/log"
)

func ConvertUTCtoTimeZone(utcTimeStr string, timeZone string) (string, error) {
	if utcTimeStr == "-" || utcTimeStr == "" {
		return utcTimeStr, nil
	}
	// Parse the input time string
	inputTime, err := time.Parse("2006/01/02 15:04:05", utcTimeStr)
	if err != nil {
		log.Errorf(err, "Error parsing the input time string : %s", utcTimeStr)
		return utcTimeStr, err
	}

	// Target timezone
	targetLocation, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Errorf(err, "Invalid time zone, error loading the time zone : %s", timeZone)
		return utcTimeStr, err
	}

	// Convert to target timezone
	localTime := inputTime.In(targetLocation)

	// Format the local time
	outputStr := localTime.Format("2006/01/02 15:04:05")

	return outputStr, nil
}
