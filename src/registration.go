package src

import (
	"fmt"
	"io"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

func registerCourses() {
	if Offset > 0 {
		offsetDuration := time.Duration(Offset) * time.Millisecond
		time.Sleep(offsetDuration)
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

			if registrationStatus != "success" && InfiniteRequests {
				for registrationStatus != "success" {
					registrationStatus = registerCourse(courseID, action, units)
				}
			}
		}(i, courseID, action, units)
	}

	wg.Wait()

	for _, status := range registrationStatuses {
		if status.Status == "success" {
			fmt.Printf("✅ %s. Successfully registered.\n", status.CourseID)
		} else {
			fmt.Printf("❌ %s. Failed to register. Reason: %s\n", status.CourseID, status.Status)
		}
	}

	close(done)
}

func registerCourse(courseID, action string, units string) string {
	for retries := 0; retries < MaxRetries; retries++ {
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

		var jsonResponse Response
		err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
		if err != nil {
			return fmt.Sprintf("Error decoding JSON response for %s: %v", courseID, err)
		}

		if len(jsonResponse.Jobs) > 0 {
			result := jsonResponse.Jobs[0].Result
			if result == "OK" || result == "COURSE_DUPLICATE" {
				fmt.Printf("%s✅ %s.%s\n", GREEN, courseID, RESET)
				return "success"
			}
			reason := getRegistrationFailureReason(body)
			fmt.Printf("%s❌ %s. Reason: %s.%s\n", RED, courseID, reason, RESET)
		}
		time.Sleep(time.Duration(DelaySeconds) * time.Second)
	}
	return "Max retries reached"
}
