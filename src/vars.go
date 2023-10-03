package src

const (
	registrationURL = "https://my.edu.sharif.edu/api/reg"
	wsURL           = "wss://my.edu.sharif.edu/api/ws?token=%s"
	RED				= "\033[91m"
	GREEN			= "\033[92m"
	RESET			= "\033[0m"
)

type UserStateMessage struct {
	Message struct {
		Favorites        []string `json:"favorites"`
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

type Response struct {
	Jobs []struct {
		Result string `json:"result"`
	} `json:"jobs"`
}

var done chan struct{}
var MaxRetries int
var DelaySeconds int
var InfiniteRequests bool
var OnTimeRegistration bool
var registrationTime float64
var ConfigFileName string
var Offset int

var authToken string
var favoriteCourses []string
var action string
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
