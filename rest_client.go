package zulip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type RESTClient interface {
	DoRequest(ctx context.Context, method, path string, data map[string]any, response APIResponse, opts ...ClientSendRequestOption) error
	DoFileRequest(ctx context.Context, method, path string, fileName string, file io.Reader, response APIResponse, opts ...ClientSendRequestOption) error
}

// Client is the main HTTP Client to interact with Zulip's API
type Client struct {
	baseURL          string
	userAgent        string
	userEmail        string
	userAPIKey       string
	httpClient       *http.Client
	printRequestData bool
	printRawResponse bool
}

const (
	RESTClientDefaultTimeout  = 5 * time.Second
	RESTClientLongPollTimeout = 10 * time.Minute
)

type clientOptions struct {
	httpClient       *http.Client
	printRequestData bool
	printRawResponse bool
}

type ClientOption func(*clientOptions) error

func WithHTTPClient(client *http.Client) ClientOption {
	return func(o *clientOptions) error {
		if client == nil {
			return errors.New("http client is nil")
		}
		o.httpClient = client
		return nil
	}
}

func WithPrintRequestData() ClientOption {
	return func(o *clientOptions) error {
		o.printRequestData = true
		return nil
	}
}

func WithPrintRawResponse() ClientOption {
	return func(o *clientOptions) error {
		o.printRawResponse = true
		return nil
	}
}

func NewClient(baseURL, email, apikey string, options ...ClientOption) (*Client, error) {
	opts := clientOptions{
		httpClient:       &http.Client{},
		printRequestData: false,
		printRawResponse: false,
	}
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, err
		}
	}

	return &Client{
		baseURL:          baseURL,
		userAgent:        AgentName + "/" + Version,
		userEmail:        email,
		userAPIKey:       apikey,
		httpClient:       opts.httpClient,
		printRequestData: opts.printRequestData,
		printRawResponse: opts.printRawResponse,
	}, nil
}

type clientSendRequestOptions struct {
	timeout time.Duration
}

type ClientSendRequestOption func(*clientSendRequestOptions)

func WithTimeout(duration time.Duration) ClientSendRequestOption {
	return func(o *clientSendRequestOptions) {
		o.timeout = duration
	}
}

// DoRequest is the main function to send requests to Zulip's API.
func (c *Client) DoRequest(ctx context.Context, method, path string, data map[string]any, response APIResponse, opts ...ClientSendRequestOption) error {
	options := clientSendRequestOptions{
		timeout: RESTClientDefaultTimeout,
	}
	for _, opt := range opts {
		opt(&options)
	}

	formData := url.Values{}
	for k, v := range data {
		formData.Set(k, fmt.Sprintf("%v", v))
	}

	formDataEncoded := formData.Encode()
	var body io.Reader
	if method != http.MethodGet {
		body = strings.NewReader(formDataEncoded)
	}

	fullURLPath := c.baseURL + path
	if c.printRequestData {
		log.Printf("DEBUG: [%s] %s - %s", method, fullURLPath, formDataEncoded)
	}

	if method == http.MethodGet && len(data) > 0 {
		fullURLPath += "?" + formDataEncoded
	}

	reqCtx, reqCancel := context.WithTimeout(ctx, options.timeout)
	defer reqCancel()

	req, err := http.NewRequestWithContext(reqCtx, method, fullURLPath, body)
	if err != nil {
		return fmt.Errorf("creating send request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.userEmail, c.userAPIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("cannot read response body: %s", err)
	}

	if c.printRawResponse {
		jsonResponse, _ := response.MarshalJSON()
		log.Printf("DEBUG: %s", jsonResponse)
	}

	response.SetHTTPCode(resp.StatusCode)
	response.SetHTTPHeaders(resp.Header)

	return nil
}

// DoFileRequest is the main function to send requests to Zulip's API with a file. For file and emoji uploads.
func (c *Client) DoFileRequest(ctx context.Context, method, path string, fileName string, file io.Reader, response APIResponse, opts ...ClientSendRequestOption) error {
	options := clientSendRequestOptions{
		timeout: RESTClientDefaultTimeout,
	}
	for _, opt := range opts {
		opt(&options)
	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			"filename",
			filepath.Base(fileName)))
	h.Set("Content-Type", mime.TypeByExtension(filepath.Ext(fileName)))

	part, err := writer.CreatePart(h)
	if err != nil {
		return fmt.Errorf("cannot create writer from file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("copying file content: %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("closing writer: %v", err)
	}

	reqCtx, reqCancel := context.WithTimeout(ctx, options.timeout)
	defer reqCancel()

	fullURLPath := c.baseURL + path
	if c.printRequestData {
		log.Printf("DEBUG: [%s] %s - %s", method, fullURLPath, writer.FormDataContentType())
	}

	req, err := http.NewRequestWithContext(reqCtx, method, fullURLPath, &requestBody)
	if err != nil {
		return fmt.Errorf("creating send request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(c.userEmail, c.userAPIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if c.printRawResponse {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("DEBUG: Raw response: %s", string(bodyBytes))
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}

	response.SetHTTPCode(resp.StatusCode)
	response.SetHTTPHeaders(resp.Header)

	return nil
}
