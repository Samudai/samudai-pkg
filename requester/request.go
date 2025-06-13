package requester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// RequestError - Generic error for requester
type RequestError struct {
	ErrorString string
}

// Error - returns the error string
func (re *RequestError) Error() string {
	return re.ErrorString
}

var (
	// BadRequestError - the client sent an invalid request
	BadRequestError = &RequestError{"bad request"}
	// ConflictError - the request created a data conflict
	ConflictError = &RequestError{"conflict"}
	// RequestFailedError - the request failed
	RequestFailedError = &RequestError{"request failed"}
	// InvalidResponseError - the response is invalid
	InvalidResponseError = &RequestError{"response invalid"}
	// ServerError - the server returned a 5xx error
	ServerError = &RequestError{"internal server error"}
	// MaxRetriesError - the request has hit the maximum retries
	MaxRetriesError = &RequestError{"maximum retries reached"}
	// ResourceNotFound - the requested resource is not found
	ResourceNotFound = &RequestError{"resource not found"}
	// UnprocessableError - the request was not processable
	UnprocessableError = &RequestError{"unprocessable"}
)

// Post - Send a POST HTTP request
func Post(url string, params interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(params)
	if err != nil {
		return []byte{}, err
	}

	return DoRequest(RequestParams{
		Body:   jsonBody,
		URL:    url,
		Method: "POST",
	})
}

// Put - Send a PUT HTTP request
func Put(url string, params interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(params)
	if err != nil {
		return []byte{}, err
	}

	return DoRequest(RequestParams{
		Body:   jsonBody,
		URL:    url,
		Method: "PUT",
	})
}

// Get - Send a GET HTTP request
func Get(url string) ([]byte, error) {
	return DoRequest(RequestParams{
		Body:   nil,
		URL:    url,
		Method: "GET",
	})
}

// Delete - Send a DELETE HTTP request
func Delete(url string) ([]byte, error) {
	return DoRequest(RequestParams{
		Body:   nil,
		URL:    url,
		Method: "DELETE",
	})
}

// RequestParams - paramters for the DoRequest function
type RequestParams struct {
	Body         []byte
	URL          string
	Method       string
	RequestToken *string

	MaxRetries    int
	PauseDuration time.Duration

	CacheKey  *string
	CacheTime *time.Duration

	Headers http.Header
}

// DoRequest - The underlying function for Get/Post/Put/Delete
func DoRequest(params RequestParams) ([]byte, error) {
	if params.PauseDuration == 0 {
		// set pause duration to 1/4 of a second by default
		params.PauseDuration = time.Duration(time.Millisecond * 250)
	}

	// original request
	b, err := request(params, 0)
	if err != nil && (err == RequestFailedError || err == ServerError || err == InvalidResponseError) {
		for i := 0; i < params.MaxRetries; i++ {
			time.Sleep(params.PauseDuration)
			b, retryErr := request(params, 0)
			if retryErr != nil && (retryErr == RequestFailedError || retryErr == ServerError || retryErr == InvalidResponseError) {
				// try again after waiting a little bit
				continue
			}
			return b, retryErr
		}
	}
	return b, err
}

// inner function that's looped based on retry limit
func request(params RequestParams, count int) ([]byte, error) {
	// request all the things
	b := make([]byte, 0)

	client := &http.Client{
		Timeout: time.Duration(25 * time.Second),
	}

	var reader io.Reader
	if len(params.Body) > 0 {
		reader = bytes.NewBuffer(params.Body)
	}

	req, err := http.NewRequest(params.Method, params.URL, reader)
	if err != nil {
		log.Println(fmt.Sprintf("Error creating request %s %s with body %s - %v", params.Method, params.URL, string(params.Body), err))
		return b, RequestFailedError
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Requester", "service")

	if params.RequestToken != nil {
		req.Header.Add("Request-Token", *params.RequestToken)
	}

	// send the request
	resp, err := client.Do(req)

	// if there was a request error, fail
	if err != nil {
		log.Println(fmt.Sprintf("Request error %s %s with body %s - %v", params.Method, params.URL, string(params.Body), err))
		return b, RequestFailedError
	}

	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(fmt.Sprintf("Invalid response error when requesting %s %s with body %s - %v", params.Method, params.URL, string(params.Body), err))
		return b, InvalidResponseError
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Println(fmt.Sprintf("Non-2xx response %s %s with body %s. Got response %s (%d)", params.Method, params.URL, string(params.Body), string(b), resp.StatusCode))
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return b, BadRequestError
		case http.StatusUnprocessableEntity:
			return b, UnprocessableError
		case http.StatusConflict:
			return b, ConflictError
		case http.StatusNotFound:
			return b, ResourceNotFound
		default:
			return b, ServerError
		}
	}

	return b, nil
}
