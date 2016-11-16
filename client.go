package tumblrclient

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"errors"
	"strings"
	"github.com/dghubble/oauth1"
	"github.com/tumblr/tumblr.go"
	"golang.org/x/net/context"
)

const apiBase = "https://api.tumblr.com/v2/"

// The Tumblr API Client object
type Client struct {
	tumblr.ClientInterface
	consumer *oauth1.Config
	user *oauth1.Token
	client *http.Client
}

// Constructor with only the consumer key and secret
func NewClient(consumerKey string, consumerSecret string) *Client {
	c := Client{}
	c.SetConsumer(consumerKey, consumerSecret)
	return &c
}

// Constructor with consumer key/secret and user token/secret
func NewClientWithToken(consumerKey string, consumerSecret string, token string, tokenSecret string) *Client {
	c := NewClient(consumerKey, consumerSecret)
	c.SetToken(token, tokenSecret)
	return c
}

// Set consumer credentials, invalidates any previously cached client
func (c *Client) SetConsumer(consumerKey string, consumerSecret string) {
	c.consumer = oauth1.NewConfig(consumerKey, consumerSecret)
	c.client = nil
}

// Set user credentials, invalidates any previously cached client
func (c *Client) SetToken(token string, tokenSecret string) {
	c.user = oauth1.NewToken(token, tokenSecret)
	c.client = nil
}

// Issue GET request to Tumblr API
func (c *Client) Get(endpoint string) (tumblr.Response, error) {
	return c.GetWithParams(endpoint, url.Values{})
}

// Issue GET request to Tumblr API with param values
func (c *Client) GetWithParams(endpoint string, params url.Values) (tumblr.Response, error) {
	return getResponse(c.GetHttpClient().Get(createRequestURI(appendPath(apiBase,endpoint),params)))
}

// Issue POST request to Tumblr API
func (c *Client) Post(endpoint string) (tumblr.Response, error) {
	return c.PostWithParams(endpoint, url.Values{});
}

// Issue POST request to Tumblr API with param values
func (c *Client) PostWithParams(endpoint string, params url.Values) (tumblr.Response, error) {
	return getResponse(c.GetHttpClient().PostForm(appendPath(apiBase, endpoint), params))
}

// Issue PUT request to Tumblr API
func (c *Client) Put(endpoint string) (tumblr.Response, error) {
	return c.PutWithParams(endpoint, url.Values{});
}

// Issue PUT request to Tumblr API with param values
func (c *Client) PutWithParams(endpoint string, params url.Values) (tumblr.Response, error) {
	req, err := http.NewRequest("PUT", createRequestURI(appendPath(apiBase, endpoint), params), strings.NewReader(""))
	if err == nil {
		return getResponse(c.GetHttpClient().Do(req))
	}
	return tumblr.Response{}, err
}

// Issue DELETE request to Tumblr API
func (c *Client) Delete(endpoint string) (tumblr.Response, error) {
	return c.DeleteWithParams(endpoint, url.Values{});
}

// Issue DELETE request to Tumblr API with param values
func (c *Client) DeleteWithParams(endpoint string, params url.Values) (tumblr.Response, error) {
	req, err := http.NewRequest("DELETE", createRequestURI(appendPath(apiBase, endpoint), params), strings.NewReader(""))
	if err == nil {
		return getResponse(c.GetHttpClient().Do(req))
	}
	return tumblr.Response{}, err
}

// Retrieve the underlying HTTP client
func (c *Client) GetHttpClient() *http.Client {
	if c.consumer == nil {
		panic("Consumer credentials are not set")
	}
	if c.user == nil {
		c.SetToken("", "")
	}
	if c.client == nil {
		c.client = c.consumer.Client(context.TODO(), c.user)
		c.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	return c.client
}

// Helper function to ease appending path to a base URI
func appendPath(base string, path string) string {
	// if path starts with `/` shave it off
	if path[0] == '/' {
		path = path[1:]
	}
	return base + path
}

// Helper function to create a URI with query params
func createRequestURI(base string, params url.Values) string {
	if len(params) != 0 {
		base += "?" + params.Encode()
	}
	return base
}

// Standard way of receiving data from the API response
func getResponse(resp *http.Response, e error) (tumblr.Response, error) {
	response := tumblr.Response{}
	if e != nil {
		return response, e
	}
	defer resp.Body.Close()
	response.Headers = resp.Header
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return response, e
	}
	response = *tumblr.NewResponse(body, resp.Header)
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return response, errors.New(resp.Status)
	}
	return response, nil
}

// Creates a PostRef out of an id and blog name
func (c *Client) GetPost(id uint64, blogName string) (*tumblr.PostRef) {
	return tumblr.NewPostRef(c, &tumblr.MiniPost{
		Id: id,
		BlogName: blogName,
	})
}

// Creates a BlogRef out of the provided name
func (c *Client) GetBlog(name string) (*tumblr.BlogRef) {
	return tumblr.NewBlogRef(c, name)
}

// Makes a request for user info based on the client's user token/secret values
func (c *Client) GetUser() (*tumblr.User, error) {
	return tumblr.GetUserInfo(c)
}

// Makes a request for the user's dashboard
func (c *Client) GetDashboard() (*tumblr.Dashboard, error) {
	return c.GetDashboardWithParams(url.Values{})
}

// Makes a request for the user's dashboard with params
func (c *Client) GetDashboardWithParams(params url.Values) (*tumblr.Dashboard, error) {
	return tumblr.GetDashboard(c, params)
}

// Makes a request for
func (c *Client) GetLikes() (*tumblr.Likes, error) {
	return c.GetLikesWithParams(url.Values{})
}

// Retrieves the posts the current user has liked
func (c *Client) GetLikesWithParams(params url.Values) (*tumblr.Likes, error) {
	return tumblr.GetLikes(c, params)
}

// Performs a tagged serach with this client, returning the result
func (c *Client) TaggedSearch(tag string) (*tumblr.SearchResults, error) {
	return tumblr.TaggedSearch(c, tag, url.Values{})
}

// Performs a tagged serach with this client, returning the result
func (c *Client) TaggedSearchWithParams(tag string, params url.Values) (*tumblr.SearchResults, error) {
	return tumblr.TaggedSearch(c, tag, params)
}