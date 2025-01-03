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
func (aer APIResponseBase) HTTPCode() int {
	return aer.httpCode
}

// HTTPHeaders returns the HTTP headers of the response.
func (aer APIResponseBase) HTTPHeaders() http.Header {
	return aer.httpHeaders
}

// SetHTTPCode sets the HTTP status code of the response.
func (aer *APIResponseBase) SetHTTPCode(code int) {
	aer.httpCode = code
}

// SetHTTPHeaders sets the HTTP headers of the response.
func (aer *APIResponseBase) SetHTTPHeaders(headers http.Header) {
	aer.httpHeaders = headers.Clone()
}

func (aer APIResponseBase) XRateLimitRemaining() string {
	return aer.httpHeaders.Get(XRateLimitRemaining)
}

func (aer APIResponseBase) XRateLimitLimit() string {
	return aer.httpHeaders.Get(XRateLimitLimit)
}

func (aer APIResponseBase) XRateLimitReset() string {
	return aer.httpHeaders.Get(XRateLimitReset)
}

// Msg returns the human-readable error message string.
func (aer APIResponseBase) Msg() string {
	return aer.msg
}

// Result returns either "error" or "success".
func (aer APIResponseBase) Result() string {
	return aer.result
}

// Code returns a machine-readable error string.
func (aer APIResponseBase) Code() string {
	return aer.code
}

// IsError returns true if the result is an error.
func (aer APIResponseBase) IsError() bool {
	return aer.result == ResultError
}

// IsSuccess returns true if the result is a success.
func (aer APIResponseBase) IsSuccess() bool {
	return aer.result == ResultSuccess
}

// FieldValue returns the value of a field in the response.
func (aer APIResponseBase) FieldValue(field string) (any, error) {
	if v, found := aer.allFields[field]; found {
		return v, nil
	}

	return nil, fmt.Errorf("field '%s' not found in extra response fields", field)
}

// AllFields returns all the fields in the response.
func (aer APIResponseBase) AllFields() map[string]any {
	return aer.allFields
}

func (aer *APIResponseBase) UnmarshalJSON(b []byte) error {
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

	*aer = APIResponseBase{
		code:      code,
		msg:       msg,
		result:    result,
		allFields: allFields,
	}

	return nil
}

func (aer APIResponseBase) MarshalJSON() ([]byte, error) {
	return json.Marshal(aer.allFields)
}
