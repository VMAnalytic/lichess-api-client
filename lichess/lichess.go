package lichess

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/VMAnalytic/lichess-api-client/lichess/decoders"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

const (
	defaultBaseURL        = "https://lichess.org/"
	userAgent             = "go-lichess-api-client"
	contentType           = "application/json"
	mediaTypeEnableNDJson = "application/x-ndjson"
)

type Client struct {
	client *http.Client

	baseURL   *url.URL
	UserAgent string
	apiKey    string

	rateLimiter *rate.Limiter

	common service

	// Services used for talking to different parts of the lichess API.
	Users   *UsersService
	Account *AccountService
	Games   *GamesService
}

func NewClient(apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = defaultHTTPClient()
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	rl := rate.NewLimiter(rate.Every(1*time.Second), 20)

	c := &Client{client: httpClient, baseURL: baseURL, apiKey: apiKey, UserAgent: userAgent}
	c.common.client = c
	c.rateLimiter = rl
	c.Users = (*UsersService)(&c.common)
	c.Account = (*AccountService)(&c.common)
	c.Games = (*GamesService)(&c.common)

	return c
}

func defaultHTTPClient() *http.Client {
	var transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Second * 5,
		}).DialContext}

	return &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}
}

func (c *Client) SetLimits(limit time.Duration, burst int) error {
	if burst <= 0 {
		return errors.New("burst should be > 0")
	}

	c.rateLimiter = rate.NewLimiter(rate.Every(limit), burst)

	return nil
}

type Response struct {
	*http.Response
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}

	return response
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("baseURL must have a trailing slash, but %q does not", c.baseURL)
	}

	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(true)
		err := enc.Encode(body)

		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("Accept", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.bareDo(ctx, req)
	if err != nil {
		return resp, errors.Wrap(err, "www")
	}
	defer resp.Body.Close()

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		if resp.Header.Get("Content-Type") == mediaTypeEnableNDJson {
			decErr := decoders.NewDecoder(resp.Body).Decode(v)
			return resp, decErr
		}

		decErr := json.NewDecoder(resp.Body).Decode(v)

		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}

		if decErr != nil {
			err = decErr
		}
	}

	return resp, err
}

func (c *Client) bareDo(ctx context.Context, req *http.Request) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}

	err := c.rateLimiter.Wait(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)

	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}

	response := newResponse(resp)

	err = c.checkResponse(resp)

	return response, err
}

func (c *Client) checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)

	if err == nil && data != nil {
		_ = json.Unmarshal(data, errorResponse)
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	switch {
	case r.StatusCode == http.StatusTooManyRequests:
		return &RateLimitError{
			Rate:     0,
			Response: errorResponse.Response,
			Message:  errorResponse.Message,
		}
	default:
		return errorResponse
	}
}

type service struct {
	client *Client
}
