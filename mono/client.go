package mono

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/adelowo/go-mono/mono/util"
)

const (
	Version          = "0.1.0"
	defaultBaseURL   = "https://api.withmono.com/v2"
	defaultUserAgent = "go-mono" + "/" + Version
)

var errNonNilContext = errors.New("context must be non-nil")

type BaseMonoResponse struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type NoopRequestBody struct{}

// ToReader converts any struct into a io#Reader that can be used
func ToReader[T NoopRequestBody | any](t T) (io.Reader, error) {
	b := bytes.NewBuffer(nil)

	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	err := enc.Encode(t)

	return b, err
}

type service struct {
	client *Client
}

type Client struct {
	httpClient *http.Client
	userAgent  string
	apikey     string

	Account       *AccountService
	Authorisation *AuthorisationService
	Transaction   *TransactionService
}

func New(opts ...Option) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	if util.IsStringEmpty(c.apikey) {
		return nil, errors.New("please provide your API key")
	}

	if c.httpClient == nil {
		return nil, errors.New("please provide a useable HTTP client")
	}

	srv := &service{client: c}

	c.Account = (*AccountService)(srv)
	c.Authorisation = (*AuthorisationService)(srv)
	c.Transaction = (*TransactionService)(srv)

	return c, nil
}

func (c *Client) newRequest(method, resource string, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(resource, "/") {
		return nil, errors.New("resource must contain a / prefix")
	}

	req, err := http.NewRequest(method, defaultBaseURL+resource, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("mono-sec-key", c.apikey)

	return req, nil
}

type Response struct {
	*http.Response
}

func (c *Client) Do(ctx context.Context, req *http.Request, v any) (*Response, error) {
	if ctx == nil {
		return nil, errNonNilContext
	}

	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if e, ok := err.(*url.Error); ok {
			return nil, e
		}

		return nil, err
	}

	if resp.StatusCode > http.StatusCreated {

		var s struct {
			Message string `json:"message"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
			return nil, err
		}

		return nil, errors.New(s.Message)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	case *string:
		var s strings.Builder
		_, err = io.Copy(&s, resp.Body)
		if err == nil {
			*v = s.String()
		}
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)

		switch decErr {
		case io.EOF:
			err = nil
		default:
			err = decErr
		}
	}

	if err != nil {
		return nil, err
	}

	return &Response{resp}, err
}
