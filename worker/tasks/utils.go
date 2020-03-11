package tasks

import (
	"fmt"
	"time"
)

const (
	layoutISO = "2006-01-02"
)

func GetNextDay(date string) (string, error) {
	t, err := time.Parse(layoutISO, date);
	if err != nil {
		return  date, fmt.Errorf("Couldn't calculate next day for date : %s", date);
	}

	nextDate := t.AddDate(0, 0, +1)
	return nextDate.Format(layoutISO), nil
}

func NormalizeDataMap(input map[string]interface{}) map[string]string{
	output := make(map[string]string)
	for k, v := range input {
		switch v.(type) {
		case int:
			output[k] = fmt.Sprintf("%d", v)
		case uint:
			output[k] = fmt.Sprintf("%d", v)
		case uint64:
			output[k] = fmt.Sprintf("%d", v)
		case float32:
			output[k] = fmt.Sprintf("%f", v)
		case float64:
			output[k] = fmt.Sprintf("%f", v)
		default:
			output[k] = v.(string)
		}
	}
	return output
}