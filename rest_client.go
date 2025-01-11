package zulip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RESTClient interface {
	DoRequest(ctx context.Context, method, path string, data map[string]any, response APIResponse, opts ...DoRequestOption) error
	DoFileRequest(ctx context.Context, method, path string, fileName string, file io.Reader, response APIResponse, opts ...DoRequestOption) error
}

// Client is the main HTTP Client to interact with Zulip's API
type Client struct {
	baseURL    string
	userAgent  string
	userEmail  string
	userAPIKey string
	httpClient *http.Client
	logger     *slog.Logger
}

const (
	RESTClientDefaultTimeout  = 5 * time.Second
	RESTClientLongPollTimeout = 10 * time.Minute
)

type clientOptions struct {
	httpClient *http.Client
	userAgent  string
	logger     *slog.Logger
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

func WithCustomUserAgent(userAgent string) ClientOption {
	return func(o *clientOptions) error {
		o.userAgent = userAgent
		return nil
	}
}

func WithLogger(logger *slog.Logger) ClientOption {
	return func(o *clientOptions) error {
		if logger == nil {
			return errors.New("logger is nil")
		}
		o.logger = logger
		return nil
	}
}

func NewClientFromConfig(file, section string) (*Client, error) {
	zuliprc, err := ParseZuliprc(file)
	if err != nil {
		return nil, err
	}

	apiSection, ok := zuliprc[section]
	if !ok {
		return nil, fmt.Errorf("no '%s' section found in zuliprc file '%s'", section, file)
	}

	return NewClient(apiSection.Site, apiSection.Email, apiSection.Key)
}

func NewClient(site, email, key string, options ...ClientOption) (*Client, error) {
	opts := clientOptions{
		httpClient: &http.Client{},
		userAgent:  DefaultUserAgentName + "/" + Version,
		logger:     slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})),
	}
	for _, opt := range options {
		if err := opt(&opts); err != nil {
			return nil, err
		}
	}

	return &Client{
		baseURL:    site,
		userAgent:  opts.userAgent,
		userEmail:  email,
		userAPIKey: key,
		httpClient: opts.httpClient,
		logger:     opts.logger,
	}, nil
}

type clientSendRequestOptions struct {
	timeout time.Duration
}

type DoRequestOption func(*clientSendRequestOptions)

func WithTimeout(duration time.Duration) DoRequestOption {
	return func(o *clientSendRequestOptions) {
		o.timeout = duration
	}
}

// DoRequest is the main function to send requests to Zulip's API.
func (c *Client) DoRequest(ctx context.Context, method, path string, data map[string]any, response APIResponse, opts ...DoRequestOption) error {
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

	requestID := uuid.New().String()
	reqLog := c.logger.With(slog.String("request_id", requestID))

	reqLog.DebugContext(ctx, "Sending request",
		slog.String("method", method),
		slog.String("url", fullURLPath),
		slog.String("data", formDataEncoded))

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

	headersGroup := []slog.Attr{}
	for k, v := range resp.Header {
		headersGroup = append(headersGroup, slog.String(k, strings.Join(v, ", ")))
	}

	reqLog.DebugContext(ctx, "Received response",
		slog.Any("headers", slog.GroupValue(headersGroup...)),
		slog.String("request_id", requestID),
		slog.Int("status_code", resp.StatusCode),
	)

	response.SetHTTPCode(resp.StatusCode)
	response.SetHTTPHeaders(resp.Header)

	return nil
}

// DoFileRequest is the main function to send requests to Zulip's API with a file. For file and emoji uploads.
func (c *Client) DoFileRequest(ctx context.Context, method, path string, fileName string, file io.Reader, response APIResponse, opts ...DoRequestOption) error {
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

	mimeType := mime.TypeByExtension(filepath.Ext(fileName))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	h.Set("Content-Type", mimeType)

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

	requestID := uuid.New().String()
	reqLog := c.logger.With(slog.String("request_id", requestID))

	reqLog.DebugContext(ctx, "Sending file request",
		slog.String("method", method),
		slog.String("url", fullURLPath),
		slog.String("filename", fileName),
		slog.String("mimetype", mimeType),
		slog.Int("content_length", requestBody.Len()))

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

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}

	headersGroup := []slog.Attr{}
	for k, v := range resp.Header {
		headersGroup = append(headersGroup, slog.String(k, strings.Join(v, ", ")))
	}

	reqLog.DebugContext(ctx, "Received response",
		slog.Any("headers", slog.GroupValue(headersGroup...)),
		slog.Int("status_code", resp.StatusCode),
	)

	response.SetHTTPCode(resp.StatusCode)
	response.SetHTTPHeaders(resp.Header)

	return nil
}
