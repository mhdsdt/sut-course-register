package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

func registerCourses() {
	if offset > 0 {
		offsetInMiliSeconds := time.Duration(offset) * time.Millisecond
		time.Sleep(offsetInMiliSeconds)
	}

	registrationStatuses = make([]CourseRegistrationStatus, len(favoriteCourses))

	var wg sync.WaitGroup

	for i, courseID := range favoriteCourses {
		units := getCourseUnits(courseID)

		wg.Add(1)

		go func(index int, courseID, action, units string) {
			defer wg.Done()
			registrationStatus := registerCourse(courseID, action, units)
			registrationStatuses[index] = CourseRegistrationStatus{
				CourseID: courseID,
				Status:   registrationStatus,
			}

			if registrationStatus != "success" && infiniteRequests {
				for registrationStatus != "success" {
					registrationStatus = registerCourse(courseID, action, units)
				}
			}
		}(i, courseID, action, units)
	}

	wg.Wait()

	for _, status := range registrationStatuses {
		if status.Status == "success" {
			color.Green("✅ %s. Successfully registered.\n", status.CourseID)
		} else {
			color.Red("❌ %s. Failed to register. Reason: %s\n", status.CourseID, status.Status)
		}
	}

	close(done)
}

func registerCourse(courseID, action string, units string) string {
	for retries := 0; retries < maxRetries; retries++ {
		requestData := fmt.Sprintf(`{"action":"%s","course":"%s","units":%s}`, action, courseID, units)
		req, err := http.NewRequest("POST", registrationURL, strings.NewReader(requestData))
		if err != nil {
			return fmt.Sprintf("Error creating request for %s: %v", courseID, err)
		}
		for key, value := range registrationHeaders {
			req.Header.Set(key, value)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Sprintf("Error sending request for %s: %v", courseID, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Sprintf("❌ %s. Status Code: %d", courseID, resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Sprintf("Error reading response for %s: %v", courseID, err)
		}
		if strings.Contains(string(body), `"result":"OK"`) || strings.Contains(string(body), `"result":"COURSE_DUPLICATE"`) {
			color.Green("✅ %s.\n", courseID)
			return "success"
		} else {
			reason := getRegistrationFailureReason(body)
			color.Red("❌ %s. Reason: %s.\n", courseID, reason)
			time.Sleep(time.Duration(delaySeconds) * time.Second)
		}
	}
	return "Max retries reached"
}
