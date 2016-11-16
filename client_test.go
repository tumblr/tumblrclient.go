package tumblrclient

import (
	"testing"
	"net/http"
	"net/url"
	"errors"
	"strings"
	"io/ioutil"
	"fmt"
)

// Test basics
func TestNewClient(t *testing.T) {
	key := "some key"
	secret := "some secret"
	client := NewClient(key, secret)
	if client.consumer == nil {
		t.Fatal("Setting consumer key and secret should generate the consumer property")
	}
}

// Test http client is cleared any time consumer key/secret changes
func TestClient_SetConsumer(t *testing.T) {
	key := "some key"
	secret := "some secret"
	client := Client{}
	client.client = http.DefaultClient
	client.SetConsumer(key, secret)
	if client.consumer == nil {
		t.Fatal("Setting consumer key and secret should generate the consumer property")
	}
	if client.client != nil {
		t.Fatal("Changing consumer key/secret should clear the current http client if it exists")
	}
}

func TestNewClientWithToken(t *testing.T) {
	key := "some key"
	secret := "some secret"
	token := "some token"
	tokenSecret := "some token secret"
	client := NewClientWithToken(key, secret, token, tokenSecret)
	if client.consumer == nil {
		t.Fatal("Setting consumer key and secret should generate the consumer property")
	}
	if client.user == nil {
		t.Fatal("Setting user token/secret should generate the user property")
	}
}

func TestClient_SetToken(t *testing.T) {
	token := "some token"
	tokenSecret := "some token secret"
	client := Client{}
	client.client = http.DefaultClient
	client.SetToken(token, tokenSecret)
	if client.user == nil {
		t.Fatal("Setting user token/secret should generate the user property")
	}
	if client.client != nil {
		t.Fatal("Changing consumer key/secret should clear the current http client if it exists")
	}
}

func TestClient_GetHttpClient(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Client should panic if attempting to get HTTP Client without credentials")
		}
	}()
	client := Client{}
	client.GetHttpClient()

}

func TestClient_GetHttpClient2(t *testing.T) {
	key := "some key"
	secret := "some secret"
	client := NewClient(key, secret)
	client.GetHttpClient()
	if client.user == nil {
		t.Error("Getting HTTP client without tokens should generate user with empty values")
	}
}

func TestAppendPath(t *testing.T) {
	base := "base"
	path := "path"
	if result := appendPath(base, path); result != base + path {
		t.Error("Basic append path failed")
	}
}

func TestAppendPath2(t *testing.T) {
	base := "base"
	path := "path"
	if result := appendPath(base, "/"+path); result != base + path {
		t.Error("Slash prefix append path failed")
	}
}

func TestCreateRequestUri(t *testing.T) {
	type requestUriTestCase struct {
		base string
		params url.Values
		expected string
	}
	cases := []requestUriTestCase{
		requestUriTestCase{"base", url.Values{}, "base"},
		requestUriTestCase{"base", url.Values{"test": []string{""}}, "base?test="},
		requestUriTestCase{"base", url.Values{"test": []string{"value"}}, "base?test=value"},
		requestUriTestCase{"base", url.Values{"test": []string{"value with space"}}, "base?test=value+with+space"},
		requestUriTestCase{"base", url.Values{"test1": []string{"value1"}, "test2": []string{"value2"}}, "base?test1=value1&test2=value2"},
	}
	for _,c := range cases {
		if result := createRequestURI(c.base, c.params); result != c.expected {
			t.Errorf("Failed to create request URI. Got `%s` Expected `%s` from `%s` + %v", result, c.expected, c.base, c.params)
		}

	}
}

func TestGetResponse(t *testing.T) {
	resp := http.Response{}
	err := errors.New("Response error")
	if _,e := getResponse(&resp, err); e != err {
		t.Fatal("Response error should be returned")
	}
}

type failReader struct {
	strings.Reader
}

func (r *failReader) Read(b []byte) (n int, err error) {
	return 0, errors.New("Always fail")
}

func TestGetResponse2(t *testing.T) {
	body := "some string"
	successCodes := []int{200, 201, 300, 301, 302}
	resp := http.Response{
		Header: http.Header{},
		Body: ioutil.NopCloser(strings.NewReader(body)),
	}
	for _,code := range successCodes {
		resp.StatusCode = code
		resp.Status = fmt.Sprintf("%d Message", code)
		if _,e := getResponse(&resp, nil); e != nil {
			t.Fatalf("Response should succeed on %d response", code)
		}
	}
	// not comprehensive but you get the idea
	failCodes := []int{199, 400, 401, 403, 404, 429, 500, 501, 502, 503}
	for _,code := range failCodes {
		resp.StatusCode = code
		resp.Status = fmt.Sprintf("%d Message", code)
		if _,e := getResponse(&resp, nil); e == nil {
			t.Fatalf("Response should succeed on %d response", code)
		}
	}
}

func TestGetResponse3(t *testing.T) {
	body := "some string"
	resp := http.Response{
		StatusCode: 200,
		Status: "200 OK",
		Header: http.Header{},
		Body: ioutil.NopCloser(&failReader{Reader: *strings.NewReader(body)}),
	}
	if _,e := getResponse(&resp, nil); e == nil {
		t.Fatal("Expected read error to be returned")
	}
}