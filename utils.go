package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func formatDuration(d time.Duration) string {
    days := d / (24 * time.Hour)
    d -= days * 24 * time.Hour
    hours := d / time.Hour
    d -= hours * time.Hour
    minutes := d / time.Minute
    d -= minutes * time.Minute
    seconds := d / time.Second

    return fmt.Sprintf("%dd %02dh %02dm %02ds", days, hours, minutes, seconds)
}

func getCourseUnits(courseID string) string {
    for _, course := range coursesData {
        id, ok := course["id"].(string)
        if !ok {
            continue
        }

        if id == courseID {
            units, ok := course["units"].(float64)
            if ok {
				unitsInt := int(units)
                return strconv.Itoa(unitsInt)
            }
        }
    }

    return "0"
}

func getRegistrationFailureReason(responseBody []byte) string {
	var responseJSON map[string]interface{}
	err := json.Unmarshal(responseBody, &responseJSON)
	if err == nil {
		err, ok := responseJSON["error"].(string)
		if ok {
			return err
		}
		jobs, ok := responseJSON["jobs"].([]interface{})
		if ok && len(jobs) > 0 {
			job, ok := jobs[0].(map[string]interface{})
			if ok {
				result, ok := job["result"].(string)
				if ok {
					return result
				}
			}
		}
	}
	return string(responseBody)
}