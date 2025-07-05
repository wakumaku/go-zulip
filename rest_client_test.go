package zulip_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wakumaku/go-zulip"
)

func TestRestClientDoRequest(t *testing.T) {
	email := "email@test"
	apiKey := "apikey"

	auth := email + ":" + apiKey
	base64Auth := base64.StdEncoding.EncodeToString([]byte(auth))

	expectedUserAgent := "test/user-agent"
	expectedContentType := "application/x-www-form-urlencoded"
	expectedAccept := "application/json"
	expectedMethod := http.MethodPost
	expectedEndpoint := "/endpoint"
	expectedBody := `key=value` // url.Values.Encode() format
	expectedHeaders := http.Header{
		"Content-Type": []string{expectedContentType},
		"User-Agent":   []string{expectedUserAgent},
		"Accept":       []string{expectedAccept},
		"Authorization": []string{
			"Basic " + base64Auth,
		},
		"Content-Length":  []string{fmt.Sprintf("%d", len(expectedBody))},
		"Accept-Encoding": []string{"gzip"},
	}

	requestRecorder := struct {
		headers http.Header
		method  string
		path    string
		body    []byte
	}{
		headers: make(http.Header),
		body:    make([]byte, 0),
	}

	handlerTest := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// store request data to check later
		requestRecorder.method = r.Method
		requestRecorder.path = r.URL.Path
		requestRecorder.headers = r.Header
		body, _ := io.ReadAll(r.Body)
		requestRecorder.body = body

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result": "success"}`))
	})

	mockServer := httptest.NewServer(handlerTest)
	defer mockServer.Close()

	baseURL := mockServer.URL

	client, err := zulip.NewClient(zulip.Credentials(baseURL, email, apiKey),
		zulip.WithCustomUserAgent(expectedUserAgent),
	)
	assert.NoError(t, err)

	msg := map[string]any{
		"key": "value", // will become url.Values.Encode() format
	}

	var resp zulip.APIResponseBase

	err = client.DoRequest(context.TODO(), expectedMethod, expectedEndpoint, msg, &resp,
		zulip.WithTimeout(10*time.Second),
	)

	assert.NoError(t, err)
	assert.Equal(t, zulip.ResultSuccess, resp.Result())

	assert.Equal(t, expectedMethod, requestRecorder.method)
	assert.Equal(t, expectedEndpoint, requestRecorder.path)
	assert.Equal(t, expectedHeaders, requestRecorder.headers)
	assert.Equal(t, expectedBody, string(requestRecorder.body))
}

func TestRestClientDoRequestFile(t *testing.T) {
	email := "email@test"
	apiKey := "apikey"

	auth := email + ":" + apiKey
	base64Auth := base64.StdEncoding.EncodeToString([]byte(auth))

	expectedUserAgent := "test/user-agent"
	expectedContentType := "multipart/form-data; boundary=a106139c086e0e50b2d83cbfc544d0326e7e7e2d19ffce4ae54f01c99ade"
	expectedAccept := "application/json"
	expectedMethod := http.MethodPost
	expectedEndpoint := "/endpoint"
	expectedBody := `--939e8e6271d5c6e29e77a3d261655eb34fcf1ca62384ca607c4671c0b3a3\r\nContent-Disposition: form-data; name=\"filename\"; filename=\"file.txt\"\r\nContent-Type: text/plain; charset=utf-8\r\n\r\nfile content\r\n--939e8e6271d5c6e29e77a3d261655eb34fcf1ca62384ca607c4671c0b3a3--\r\n`
	expectedHeaders := http.Header{
		"Content-Type": []string{expectedContentType},
		"User-Agent":   []string{expectedUserAgent},
		"Accept":       []string{expectedAccept},
		"Authorization": []string{
			"Basic " + base64Auth,
		},
		"Content-Length":  []string{fmt.Sprintf("%d", len(expectedBody))},
		"Accept-Encoding": []string{"gzip"},
	}

	requestRecorder := struct {
		headers http.Header
		method  string
		path    string
		body    []byte
	}{
		headers: make(http.Header),
		body:    make([]byte, 0),
	}

	handlerTest := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// store request data to check later
		requestRecorder.method = r.Method
		requestRecorder.path = r.URL.Path
		requestRecorder.headers = r.Header
		body, _ := io.ReadAll(r.Body)
		requestRecorder.body = body

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result": "success"}`))
	})

	mockServer := httptest.NewServer(handlerTest)
	defer mockServer.Close()

	baseURL := mockServer.URL

	client, err := zulip.NewClient(zulip.Credentials(baseURL, email, apiKey),
		zulip.WithCustomUserAgent(expectedUserAgent),
	)
	assert.NoError(t, err)

	fileName := "file.txt"
	msg := bytes.NewReader([]byte("file content"))

	var resp zulip.APIResponseBase

	err = client.DoFileRequest(context.TODO(), expectedMethod, expectedEndpoint, fileName, msg, &resp,
		zulip.WithTimeout(10*time.Second),
	)

	assert.NoError(t, err)
	assert.Equal(t, zulip.ResultSuccess, resp.Result())

	assert.Equal(t, expectedMethod, requestRecorder.method)
	assert.Equal(t, expectedEndpoint, requestRecorder.path)
	// We cannot check all the headers because the boundary is random (Mimetype multipart/form-data)
	// assert.Equal(t, expectedHeaders, requestRecorder.headers)
	// assert.Equal(t, expectedBody, string(requestRecorder.body))
	// Check only the headers that are not random
	assert.Equal(t, expectedHeaders.Get("User-Agent"), requestRecorder.headers.Get("User-Agent"))
	assert.Equal(t, expectedHeaders.Get("Accept"), requestRecorder.headers.Get("Accept"))
	assert.Equal(t, expectedHeaders.Get("Authorization"), requestRecorder.headers.Get("Authorization"))
	assert.Equal(t, expectedHeaders.Get("Accept-Encoding"), requestRecorder.headers.Get("Accept-Encoding"))
}
