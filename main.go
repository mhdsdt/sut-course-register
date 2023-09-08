package main

import (
	"fmt"
	"io"
	"os"
	"flag"
	"sync"
	"net/http"
	"strings"
	"time"
	"strconv"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/fatih/color"
)

const (
	registrationURL = "https://my.edu.sharif.edu/api/reg"
	wsURL           = "wss://my.edu.sharif.edu/api/ws?token=%s"
)

type UserStateMessage struct {
	Message struct {
		Favorites []string `json:"favorites"`
        RegistrationTime float64
	} `json:"message"`
}

type ListUpdateMessage struct {
	Message []map[string]interface{} `json:"message"`
}

type CourseRegistrationStatus struct {
    CourseID string
    Status   string
}

var done chan struct{}
var maxRetries int
var delaySeconds int
var infiniteRequests bool
var onTimeRegistration bool
var registrationTime float64

var authToken string
var favoriteCourses []string
var coursesData []map[string]interface{}
var registrationStatuses []CourseRegistrationStatus

var registrationHeaders = map[string]string{
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
	"Accept":          "application/json",
	"Accept-Language": "en-US,en;q=0.5",
	"Referer":         "https://my.edu.sharif.edu/courses/marked",
	"Content-Type":    "application/json",
	"Authorization":   "",
	"Origin":          "https://my.edu.sharif.edu",
	"Connection":      "keep-alive",
	"Sec-Fetch-Dest":  "empty",
	"Sec-Fetch-Mode":  "cors",
	"Sec-Fetch-Site":  "same-origin",
	"TE":              "trailers",
}

func init() {
	flag.IntVar(&delaySeconds, "d", 5, "Delay in seconds between registration attempts")
	flag.IntVar(&maxRetries, "r", 5, "Maximum number of registration retries")
	flag.BoolVar(&infiniteRequests, "i", false, "Request indefinitely until successful")
	flag.BoolVar(&onTimeRegistration, "on-time", false, "Enable on-time registration")
	flag.Parse()
}

func main() {
	err := readAuthTokenFromFile()
	if err != nil {
		fmt.Printf("Error reading auth token from file: %v\n", err)
		return
	}
	wsURLWithToken := fmt.Sprintf(wsURL, authToken)
	conn, _, err := websocket.DefaultDialer.Dial(wsURLWithToken, nil)
	if err != nil {
		fmt.Printf("Error establishing WebSocket connection: %v\n", err)
		return
	}
	defer conn.Close()

	done = make(chan struct{})

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Error reading WebSocket message: %v\n", err)
				close(done)
				return
			}
			handleMessage(message)
		}
	}()
	<-done
}

func readAuthTokenFromFile() error {
    file, err := os.Open("token.txt")
    if err != nil {
        return err
    }
    defer file.Close()

    tokenData, err := io.ReadAll(file)
    if err != nil {
        return err
    }

    authToken = strings.TrimSpace(string(tokenData))
    registrationHeaders["Authorization"] = authToken

    return nil
}

func handleMessage(message []byte) {
    var msg struct {
        Type string `json:"type"`
    }

    if err := json.Unmarshal(message, &msg); err != nil {
        fmt.Printf("Error parsing WebSocket message: %v\n", err)
        return
    }

    switch msg.Type {
    case "userState":
        handleUserState(message)
    case "listUpdate":
        handleListUpdate(message)
		registerCourses()
    default:
        fmt.Printf("Received unknown message type: %s\n", msg.Type)
    }
}

func handleUserState(message []byte) {
    var parsedMessage UserStateMessage
    err := json.Unmarshal(message, &parsedMessage)
    if err != nil {
        return
    }

    registrationTime = parsedMessage.Message.RegistrationTime
    favoriteCourses = parsedMessage.Message.Favorites
    fmt.Printf("favoriteCourses: %v\n", favoriteCourses)
}

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

func handleListUpdate(message []byte) {
	var parsedMessage ListUpdateMessage
    if err := json.Unmarshal(message, &parsedMessage); err != nil {
        fmt.Printf("Error parsing listUpdate message: %v\n", err)
        return
    }

    coursesData = parsedMessage.Message
    fmt.Println("Updated list of courses")

    timeRemaining := time.Until(time.Unix(int64(registrationTime)/1000, 0))
    fmt.Print("\rRegistration will start in ", formatDuration(timeRemaining))

    if onTimeRegistration {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
        for range ticker.C {
            timeRemaining := time.Until(time.Unix(int64(registrationTime)/1000, 0))
            fmt.Print("\rRegistration will start in ", formatDuration(timeRemaining))
            if timeRemaining <= 0 {
                break
            }
        }
    }
    fmt.Println()
}

func registerCourses() {
	registrationStatuses = make([]CourseRegistrationStatus, len(favoriteCourses))

	var wg sync.WaitGroup

	for i, courseID := range favoriteCourses {
		units := getCourseUnits(courseID)

		wg.Add(1)

		go func(index int, courseID, units string) {
			defer wg.Done()
			registrationStatus := registerCourse(courseID, units)
			registrationStatuses[index] = CourseRegistrationStatus{
				CourseID: courseID,
				Status:   registrationStatus,
			}

			if registrationStatus != "success" && infiniteRequests {
				for registrationStatus != "success" {
					registrationStatus = registerCourse(courseID, units)
				}
			}
		}(i, courseID, units)
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

func registerCourse(courseID, units string) string {
    for retries := 0; retries < maxRetries; retries++ {
        requestData := fmt.Sprintf(`{"action":"add","course":"%s","units":%s}`, courseID, units)
        req, err := http.NewRequest("POST", registrationURL, strings.NewReader(requestData))
        if err != nil {
            return fmt.Sprintf("Error creating request for %s: %v", courseID, err)
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
		fmt.Printf("Error %s", err)
    }
    return string(responseBody)
}
