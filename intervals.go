package intervals

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/google/go-querystring/query"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var lock *int

// Client contains all the core logic for making requests
type Client struct {
	clientMu sync.Mutex
	client   *http.Client

	// BaseURL and APIKey are constant set at client start
	BaseURL *url.URL
	APIKey  string

	common        service
	PersonService *PersonService
	TimeService   *TimeService
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

func init() {
	log.SetFormatter(&log.TextFormatter{})

	viper.SetConfigName("intervals")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./.config/")
	viper.AddConfigPath("~/.config/")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.WithField("ERROR", err).Fatal("Cound not read in config")
	}
}

// NewClient returns a new Intervals API client. If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseurl, err := url.Parse(viper.GetString("INTERVALS_URL"))

	if err != nil {
		log.Fatal(err)
	}

	c := &Client{client: httpClient}
	c.common.client = c
	c.PersonService = (*PersonService)(&c.common)
	c.TimeService = (*TimeService)(&c.common)
	c.BaseURL = baseurl
	c.APIKey = viper.GetString("INTERVALS_APIKEY")

	return c
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
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
	req.Header.Add("Authorization", "Basic "+basicAuth(c.APIKey, "X"))

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
