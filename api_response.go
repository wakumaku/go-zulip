package zulip

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// X-RateLimit-Remaining: The number of additional requests of this type
	// that the client can send before exceeding its limit.
	XRateLimitRemaining string = "X-RateLimit-Remaining"
	// X-RateLimit-Limit: The limit that would be applicable to a client that
	// had not made any recent requests of this type. This is useful for
	// designing a client's burst behavior so as to avoid ever reaching a
	// rate limit.
	XRateLimitLimit string = "X-RateLimit-Limit"
	// X-RateLimit-Reset: The time at which the client will no longer have any
	// rate limits applied to it (and thus could do a burst of
	// X-RateLimit-Limit requests).
	XRateLimitReset string = "X-RateLimit-Reset"

	// ResultSuccess is the string returned in the result field when there is
	// a success operation
	ResultSuccess string = "success"
	// ResultError is the string returned in the result field when there is
	// an error
	ResultError string = "error"
)

// APIResponse is the interface that wraps the basic methods of an API response.
type APIResponse interface {
	SetHTTPCode(httpCode int)
	SetHTTPHeaders(headers http.Header)
	json.Unmarshaler
	json.Marshaler
}

// APIResponseBase is the base struct for all API responses.
type APIResponseBase struct {
	httpCode    int
	httpHeaders http.Header

	// msg: an internationalized, human-readable error message string.
	msg string

	// result: either "error" or "success", which is redundant with the HTTP status code, but is convenient when print debugging.
	result string

	// code: a machine-readable error string, with a default value of
	// "BAD_REQUEST" for general errors.
	code string

	// allFields: a map of all the fields returned in the response.
	allFields map[string]any
}

// HTTPCode returns the HTTP status code of the response.
func (a APIResponseBase) HTTPCode() int {
	return a.httpCode
}

// HTTPHeaders returns the HTTP headers of the response.
func (a APIResponseBase) HTTPHeaders() http.Header {
	return a.httpHeaders
}

// SetHTTPCode sets the HTTP status code of the response.
func (a *APIResponseBase) SetHTTPCode(code int) {
	a.httpCode = code
}

// SetHTTPHeaders sets the HTTP headers of the response.
func (a *APIResponseBase) SetHTTPHeaders(headers http.Header) {
	a.httpHeaders = headers.Clone()
}

func (a APIResponseBase) XRateLimitRemaining() string {
	return a.httpHeaders.Get(XRateLimitRemaining)
}

func (a APIResponseBase) XRateLimitLimit() string {
	return a.httpHeaders.Get(XRateLimitLimit)
}

func (a APIResponseBase) XRateLimitReset() string {
	return a.httpHeaders.Get(XRateLimitReset)
}

// Msg returns the human-readable error message string.
func (a APIResponseBase) Msg() string {
	return a.msg
}

// Result returns either "error" or "success".
func (a APIResponseBase) Result() string {
	return a.result
}

// Code returns a machine-readable error string.
func (a APIResponseBase) Code() string {
	return a.code
}

// IsError returns true if the result is an error.
func (a APIResponseBase) IsError() bool {
	return a.result == ResultError
}

// IsSuccess returns true if the result is a success.
func (a APIResponseBase) IsSuccess() bool {
	return a.result == ResultSuccess
}

// FieldValue returns the value of a field in the response.
func (a APIResponseBase) FieldValue(field string) (any, error) {
	if v, found := a.allFields[field]; found {
		return v, nil
	}

	return nil, fmt.Errorf("field '%s' not found in extra response fields", field)
}

// AllFields returns all the fields in the response.
func (a APIResponseBase) AllFields() map[string]any {
	return a.allFields
}

func (a *APIResponseBase) UnmarshalJSON(b []byte) error {
	if string(b) == "null" || string(b) == `""` {
		return nil
	}

	var (
		code      string
		result    string
		msg       string
		allFields map[string]any
	)

	if err := json.Unmarshal(b, &allFields); err != nil {
		return err
	}

	if v, found := allFields["code"]; found {
		if val, ok := v.(string); ok {
			code = val
		}
	}

	if v, found := allFields["msg"]; found {
		if val, ok := v.(string); ok {
			msg = val
		}
	}

	if v, found := allFields["result"]; found {
		if val, ok := v.(string); ok {
			result = val
		}
	}

	*a = APIResponseBase{
		code:      code,
		msg:       msg,
		result:    result,
		allFields: allFields,
	}

	return nil
}

func (a APIResponseBase) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.allFields)
}
