package helper

import (
	"time"

	"github.com/calculi-corp/reports-service/constants"
	db "github.com/calculi-corp/reports-service/db"
)

func AddTermFilter(filterName string, filterValue string) map[string]interface{} {
	newElement := map[string]interface{}{
		"term": map[string]interface{}{
			filterName: filterValue,
		},
	}
	return newElement
}

func AddTermsFilter(filterName string, filterValues []string) map[string]interface{} {
	newElement := map[string]interface{}{
		"terms": map[string]interface{}{
			filterName: filterValues,
		},
	}
	return newElement
}

func CalculatePreviousDates(inputDate string, duration string) (string, string, error) {
	t, err := time.Parse(constants.DATE_FORMAT_WITH_HYPHEN, inputDate)
	if err != nil {
		return "", "", err
	}

	switch duration {
	case "month":
		// subtract a month
		startDate := time.Date(t.Year(), t.Month()-1, 1, 0, 0, 0, 0, t.Location())
		endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

		//format
		s := startDate.Format(constants.DATE_FORMAT_WITH_HYPHEN_WITHOUT_TIME)
		e := endDate.Format(constants.DATE_FORMAT_WITH_HYPHEN_WITHOUT_TIME)
		return s, e, nil
	case "week":
		// Calculate the day of the week
		dayOfWeek := t.Weekday()

		// no of days elapsed
		diff := int(dayOfWeek+6) % 7

		// start date of the prev. week
		startDate := t.AddDate(0, 0, -(diff + 7))

		// end date of the previous week
		endDate := startDate.AddDate(0, 0, +6)

		//format
		s := startDate.Format(constants.DATE_FORMAT_WITH_HYPHEN_WITHOUT_TIME)
		e := endDate.Format(constants.DATE_FORMAT_WITH_HYPHEN_WITHOUT_TIME)
		return s, e, nil
	default:
		return "", "", db.ErrInvalidDuration
	}

}
