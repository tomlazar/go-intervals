package intervals

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/google/go-querystring/query"
)

// Client contains all the core logic for making requests
type Client struct {
	clientMu sync.Mutex
	client   *http.Client
	config   Config

	// BaseURL and APIKey are constant set at client start
	BaseURL *url.URL

	common         service
	PersonService  *PersonService
	TimeService    *TimeService
	ProjectService *ProjectService
}

// Config options for the client
type Config struct {
	IntervalsURL string
	APIKey       string
}

type baseResponse struct {
	PersonID  int    `json:"personid"`
	Status    string `json:"status"`
	Code      int    `json:"code"`
	ListCount int    `json:"listcount"`
}

type service struct {
	client *Client
}

// NewClient returns a new Intervals API client. If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(config Config, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseurl, err := url.Parse(config.IntervalsURL)

	if err != nil {
		log.Fatal(err)
	}

	c := &Client{client: httpClient}
	c.common.client = c
	c.PersonService = (*PersonService)(&c.common)
	c.TimeService = (*TimeService)(&c.common)
	c.ProjectService = (*ProjectService)(&c.common)
	c.BaseURL = baseurl
	c.config = config

	return c
}

// NewRequest creates an API request. A relative url can be provided in urlStr,
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
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
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", c.config.auth())

	return req, nil
}

// Do runs the request and sends it to the API
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil
			}
		}
	}

	return resp, nil
}

func (c Config) auth() string {
	bytes := []byte(fmt.Sprintf("%v:X", c.APIKey))
	tok := base64.StdEncoding.EncodeToString(bytes)
	return fmt.Sprintf("Basic %v", tok)
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
