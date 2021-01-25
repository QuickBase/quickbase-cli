package qbclient

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"

	"github.com/QuickBase/quickbase-cli/qberrors"
	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/viper"
)

// Client makes requests to the Quick Base API.
type Client struct {
	HTTPClient     *http.Client
	Plugins        []Plugin
	ReamlHostname  string
	TemporaryToken string
	URL            string
	UserAgent      string
	UserToken      string
}

// New returns a new Client.
func New(cfg ConfigIface) *Client {
	c := &Client{
		ReamlHostname:  cfg.RealmHostname(),
		TemporaryToken: cfg.TemporaryToken(),
		URL:            "https://api.quickbase.com/v1",
		UserAgent:      userAgent(),
		UserToken:      cfg.UserToken(),
	}

	// Configure and set the retry handler.
	rh := retryablehttp.NewClient()
	rh.RetryMax = 2
	rh.Logger = nil
	rh.ErrorHandler = c.errorHandler
	c.HTTPClient = rh.StandardClient()

	return c
}

// NewFromProfile returns a new Client, initializing the config from the
// passed profile.
func NewFromProfile(profile string) (client *Client, err error) {
	cfg := viper.New()
	cfg.SetDefault(OptionProfile, profile)
	if err = ReadInConfig(cfg); err == nil {
		client = New(Config{cfg: cfg})
	}
	return
}

func userAgent() string {
	return fmt.Sprintf("quickbase-cli/%s, (%s %s)", Version, runtime.GOOS, runtime.GOARCH)
}

// AddPlugin adds a Plugin to the stack.
func (c *Client) AddPlugin(p Plugin) {
	c.Plugins = append(c.Plugins, p)
}

// Do sends an arbitrary request to the Quick Base API.
// TODO Improve the error handling.
func (c *Client) Do(input Input, output Output) error {

	// Validate the input.
	if err := validator.New().Struct(input); err != nil {
		return qberrors.HandleErrorValidation(err)
	}

	// Marshal marshals the request body using Input.marshal.
	b, err := input.encode()
	if err != nil {
		return qberrors.Client(err).Safef(qberrors.InvalidInput, "error encoding input")
	}

	// Create the request, using the marshalled input as the body.
	req, err := http.NewRequest(input.method(), input.url(), bytes.NewBuffer(b))
	if err != nil {
		serr := qberrors.ErrSafe{Message: "error creating request"}
		return qberrors.Internal(err).Safe(serr)
	}

	// Add HTTP headers using Input.addHeaders.
	input.addHeaders(req)

	// Invoke each plugin's PreRequest hook.
	c.invokePreRequest(req)

	// Do the HTTP request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		serr := qberrors.ErrSafe{Message: "error executing request"}
		return qberrors.Service(err).Safe(serr)
	}

	// Invoke each plugin's PostResponse hook.
	c.invokePostResponse(resp)

	// Parse the response body. We do our best to handle this gracefully if
	// an error is thrown outside of the API's control plane, e.g., from
	// Cloudflare, which might not produce parsable output.
	if err := output.decode(resp.Body); err != nil {
		fmt.Println(err)
		switch true {
		case resp.StatusCode >= 200 && resp.StatusCode < 300:
			serr := qberrors.ErrSafe{Message: "error decoding response"}
			return qberrors.Internal(err).Safe(serr)
		case resp.StatusCode >= 400 && resp.StatusCode < 500:
			serr := qberrors.ErrSafe{Message: http.StatusText(resp.StatusCode), StatusCode: resp.StatusCode}
			return qberrors.Client(serr).Safe(serr)
		default:
			serr := qberrors.ErrSafe{Message: http.StatusText(resp.StatusCode), StatusCode: resp.StatusCode}
			return qberrors.Service(serr).Safe(serr)
		}
	}

	// Handle any errors, the logic of which will depend on whether we are
	// consuming the XML or RESTful API.
	return output.handleError(output, resp)
}

// errorHandler implements retryablehttp.ErrorHandler by invoking post-response
// plugins after the error occurs. It then closes the response body and returns
// the same error message as retryablehttp.Do.
func (c *Client) errorHandler(resp *http.Response, err error, numTries int) (*http.Response, error) {
	c.invokePostResponse(resp)

	if resp != nil {
		resp.Body.Close()
	}

	s := fmt.Sprintf("giving up after %d attempt", numTries)
	if numTries > 1 {
		s += "s"
	}

	serr := qberrors.ErrSafe{Message: s}
	return nil, qberrors.Service(err).Safe(serr)
}

func (c *Client) invokePreRequest(req *http.Request) {
	for _, plugin := range c.Plugins {
		plugin.PreRequest(req)
	}
}

func (c *Client) invokePostResponse(resp *http.Response) {
	for _, plugin := range c.Plugins {
		plugin.PostResponse(resp)
	}
}
