package httputils

import (
	"net/http"
)

//type ErrNotFound interface {
//	error
//	NotFound()
//}
//
//type errNotFound struct {
//	Err error
//}
//
//func (errNotFound) NotFound() {
//
//}

type ErrHandlerFunc func(w http.ResponseWriter, r *http.Request, err HTTPError)

// HTTPError represents an HTTP error with HTTP status code and error message
type HTTPError interface {
	error
	// StatusCode returns the HTTP status code of the error
	StatusCode() int
	ErrHandlerFunc() ErrHandlerFunc
	//Headers() http.Header
	//AddHeader(key, value string)
	//SetHeader(key, value string)
	//DelHeader(key string)
	//GetHeader(key string) string
	WithError(err error) HTTPError
	InsideError() error
}

// apiError represents an error that can be sent in an error response.
type APPError struct {
	// Status represents the HTTP status code
	status int `json:"-"`
	// ErrorCode is the code uniquely identifying an error
	// ErrorCode string `json:"error_code"`
	// Message is the error message that may be displayed to end users
	Message      string `json:"message"`
	DebugMessage string `json:"debug_message,omitempty"`
	// Details specifies the additional error information
	Errors interface{} `json:"errors,omitempty"`
	err    error

	headers http.Header

	isJson bool
}

func (e *APPError) ErrHandlerFunc() ErrHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, err HTTPError) {
		if ae, ok := err.(*APPError); ok {
			for key, vals := range ae.Headers() {
				for _, val := range vals {
					w.Header().Add(key, val)
				}
			}
		}
		DefaultErrorHandleFunc(w, r, err)
	}
}

func (e *APPError) WithError(err error) HTTPError {
	e.err = err
	return e
}

func (e *APPError) Debug(msg string) HTTPError {
	e.DebugMessage = msg
	return e
}

func (e *APPError) InsideError() error {
	return e.err
}

// Error returns the error message.
func (e *APPError) Error() string {
	//if jsonData, err := e.ToJson(); err == nil {
	//	return string(jsonData)
	//} else {
	//	return err.Error()
	//}
	return e.Message
	//contentType := e.GetHeader("Content-Type")
	//// todo 这里可以缓存
	//if strings.Contains(contentType, "json") {
	//	// TODO: error handing
	//	b, _ := json.Marshal(e)
	//	return string(b)
	//} else {
	//	return e.Message
	//}
}

// StatusCode returns the HTTP status code.
func (e *APPError) StatusCode() int {
	return e.status
}

func (e *APPError) Headers() http.Header {
	//h := http.Header{}
	//h.Add("Content-Type", "application/json")
	return e.headers
}

func (e *APPError) AddHeader(key, value string) {
	e.headers.Add(key, value)
}

func (e *APPError) SetHeader(key, value string) {
	e.headers.Set(key, value)
}

func (e *APPError) DelHeader(key string) {
	e.headers.Del(key)
}

func (e *APPError) GetHeader(key string) string {
	return e.headers.Get(key)
}

func NewAPIError(status int, message, debugMessage string, errors ...interface{}) *APPError {
	apiError := &APPError{
		status:  status,
		Message: message,
		// ErrorCode: errorCode,
		DebugMessage: debugMessage,
		headers:      http.Header{},
	}
	if len(errors) > 0 {
		apiError.Errors = errors[0]
	}
	return apiError
}

// InternalServerError creates a new API error representing an internal server error (HTTP 500)
func InternalServerError(message string, debugMessage ...string) *APPError {

	if message == "" {
		message = http.StatusText(http.StatusInternalServerError)
	}

	debugMsg := ""

	if len(debugMessage) > 0 {
		debugMsg = debugMessage[0]
	}
	return NewAPIError(http.StatusInternalServerError, message, debugMsg)
}

// NotFound creates a new API error representing a resource-not-found error (HTTP 404)
func NotFound(message string, debugMessage ...string) *APPError {
	if message == "" {
		message = http.StatusText(http.StatusNotFound)
	}

	debugMsg := ""

	if len(debugMessage) > 0 {
		debugMsg = debugMessage[0]
	}
	return NewAPIError(http.StatusNotFound, message, debugMsg)
}

// Unauthorized creates a new API error representing an authentication failure (HTTP 401)
func Unauthorized(message ...string) *APPError {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = http.StatusText(http.StatusUnauthorized)
	}
	return NewAPIError(http.StatusUnauthorized, msg, "")
}

func Forbidden(message ...string) *APPError {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = http.StatusText(http.StatusForbidden)
	}
	return NewAPIError(http.StatusForbidden, msg, "")
}

func BadRequest(message ...string) *APPError {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	} else {
		msg = http.StatusText(http.StatusBadRequest)
	}
	return NewAPIError(http.StatusBadRequest, msg, "")
}
