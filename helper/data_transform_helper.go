package helper

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/calculi-corp/reports-service/db"
	"github.com/calculi-corp/reports-service/exceptions"

	"github.com/calculi-corp/log"
)

// takes a set of dates sent as input and adds it to the response from Jolt
func AddDateBuckets(input json.RawMessage, newDates []map[string]string) (json.RawMessage, error) {
	var inputSlice []map[string]interface{}
	err := json.Unmarshal(input, &inputSlice)
	if log.CheckErrorf(err, "error Unmarshalling in helper.AddDateBuckets() : ") {
		return nil, db.ErrInternalServer
	}

	if len(inputSlice) != 0 && len(newDates) != 0 {
		for _, outerMap := range inputSlice {
			data, ok := outerMap["data"]
			if !ok {
				log.Errorf(errors.New("error in extracting data field in inputSlice - helper.AddDateBuckets()"), exceptions.ErrHelperAddBuckets)
				return nil, db.ErrInternalServer
			}
			// Check if data is nil or not of type []interface{}
			_, ok = data.([]interface{})
			if !ok || data == nil {
				log.Errorf(errors.New("error in type assertion or data is nil - helper.AddDateBuckets()"), exceptions.ErrHelperAddBuckets)
				return nil, db.ErrInternalServer
			}
			for _, item := range data.([]interface{}) {
				item2, ok := item.(map[string]interface{})
				if !ok {
					log.Errorf(errors.New("error in type assertion - helper.AddDateBuckets()"), exceptions.ErrHelperAddBuckets)
					return nil, db.ErrInternalServer
				}
				for _, m := range newDates {
					outerDate, ok := item2["x"].(string)
					if !ok {
						log.Errorf(errors.New("error in type assertion of date - helper.AddDateBuckets()"), exceptions.ErrHelperAddBuckets)
						return nil, db.ErrInternalServer
					}
					if outerDate == m["startDate"] {
						item2["date"] = map[string]interface{}{
							"startDate": m["startDate"],
							"endDate":   m["endDate"],
						}
					}
				}
			}
		}
	} else {
		log.Errorf(errors.New("error in helper.AddDateBuckets() "), "length of input/dates is 0 in helper.AddDateBuckets() ")
		return nil, db.ErrInternalServer
	}
	rawMsg, err := json.Marshal(inputSlice)
	if log.CheckErrorf(err, "error Marshalling in helper.AddDateBuckets() : ") {
		return nil, db.ErrInternalServer
	}
	return rawMsg, nil
}

// calculates the start and end dates to the date histogram buckets which are added to Jolt's response (for widgets with clickable date charts)
func GetBucketDates(aggrBy, startDate, endDate string) ([]map[string]string, error) {
	layout := "2006-01-02"
	start, err := time.Parse(layout, startDate)
	if log.CheckErrorf(err, "error parsing input start date in helper.GetBucketDates() ") {
		return nil, db.ErrInternalServer
	}

	end, err := time.Parse(layout, endDate)
	if log.CheckErrorf(err, "error parsing input end date in helper.GetBucketDates() ") {
		return nil, db.ErrInternalServer
	}

	var dateBuckets []map[string]string

	if aggrBy == "week" {
		for !start.After(end) {
			week := make(map[string]string)
			week["startDate"] = start.Format(layout)

			if start.Weekday() == time.Sunday {
				week["endDate"] = start.Format(layout)
			} else {
				// move to the subsequent Sunday
				start = start.AddDate(0, 0, int(7-start.Weekday()))
				if start.After(end) {
					week["endDate"] = end.Format(layout)
				} else {
					week["endDate"] = start.Format(layout)
				}
			}

			dateBuckets = append(dateBuckets, week)

			// move to Monday (which would be the start date for the subsequent week)
			start = start.AddDate(0, 0, 1)

		}
	} else if aggrBy == "day" {
		for !start.After(end) {
			day := make(map[string]string)
			day["startDate"] = start.Format(layout)
			day["endDate"] = start.Format(layout)

			dateBuckets = append(dateBuckets, day)

			// move to the next day
			start = start.AddDate(0, 0, 1)

		}
	} else if aggrBy == "month" {
		for !start.After(end) {
			month := make(map[string]string)
			month["startDate"] = start.Format(layout)

			// Calculate the end date as the last date of the month
			monthEndDate := time.Date(start.Year(), start.Month()+1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)
			if monthEndDate.After(end) {
				month["endDate"] = end.Format(layout)
			} else {
				month["endDate"] = monthEndDate.Format(layout)
			}

			dateBuckets = append(dateBuckets, month)

			// Move to the start date of the subsequent month
			start = time.Date(start.Year(), start.Month()+1, 1, 0, 0, 0, 0, time.UTC)
		}
	} else {
		log.Errorf(errors.New("error in helper.GetBucketDates(): "), "invalid aggrBy argument received in helper.GetBucketDates() ")
		return nil, db.ErrInternalServer
	}

	return dateBuckets, nil
}
