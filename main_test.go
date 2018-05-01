package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/josephholsten/yasssd/store"
)

func init() {
	store.Truncate()
}

type RequestTestData struct {
	RequestMethod       string
	RequestURI          string
	RequestHeaders      map[string][]string
	RequestBody         io.Reader
	ResponseStatusCode  int
	ResponseContentType string
	ResponseBody        string
}

func TestRequests(t *testing.T) {
	Tests := []RequestTestData{
		{"POST", "/register", nil, strings.NewReader(`{"username": "pat", "password":"secretiv"}`), 204, "", ""},

		{"POST", "/register", nil, strings.NewReader(`{"username": "p", "password":"secretiv"}`), 400, "application/json", errorResp("username too short, must be at least 3 characters")},
		{"POST", "/register", nil, strings.NewReader(`{"username": "123456789012345678901", "password":"secretiv"}`), 400, "application/json", errorResp("username too long, must be no more than 20 characters")},
		{"POST", "/register", nil, strings.NewReader(`{"username": "pat*", "password":"secretiv"}`), 400, "application/json", errorResp("username contains invalid characters, must contain only alphanumeric characters")},
		{"POST", "/register", nil, strings.NewReader(`{"username": "pat", "password":"1234567"}`), 400, "application/json", errorResp("password too short, must be at least 8 characters")},
		{"POST", "/register", nil, nil, 400, "application/json", errorResp("request body was empty, expected json with username and password fields")},

		{"POST", "/login", nil, strings.NewReader(`{"username": "pat", "password":"secretiv"}`), 200, "application/json", `{"token":"1"}`},
		{"GET", "/login", nil, nil, 400, "application/json", errorResp("unexpected method:GET")},

		{"GET", "/files", map[string][]string{"X-Token": []string{"1"}}, nil, 200, "application/json", ""},
		{"GET", "/files", map[string][]string{"X-Token": []string{"non-existant"}}, nil, 403, "application/json", "{\n  \"error\": \"session token not recognized\"\n}"},

		{"GET", "/files/foo", map[string][]string{"X-Token": []string{"1"}}, nil, 200, "application/json", "foo-file-contents"},
		{"GET", "/files/foo", map[string][]string{"X-Token": []string{"non-existant"}}, nil, 403, "application/json", "{\n  \"error\": \"session token not recognized\"\n}"},
	}

	for _, test := range Tests {
		store.Truncate()
		ensureAccount("pat", "secretiv")
		ensureToken("1", "pat")
		ensureFile("pat", "foo", "foo-file-contents")

		method := test.RequestMethod
		uri := test.RequestURI
		headers := test.RequestHeaders
		body := test.RequestBody
		wantCode := test.ResponseStatusCode
		wantCt := test.ResponseContentType
		wantResp := test.ResponseBody

		resp := httptest.NewRecorder()
		req, err := http.NewRequest(method, uri, body)
		if err != nil {
			t.Errorf("Error: %v, TestData: %v", err, test)
		}
		req.Header = headers
		http.DefaultServeMux.ServeHTTP(resp, req)
		if ok, err := isValidResp(resp, wantCode, wantCt, wantResp); !ok {
			t.Errorf("Error: %v, TestData: %v", err, test)
		}
	}

}

func isValidResp(resp *httptest.ResponseRecorder, wantCode int, wantCt, wantBody string) (bool, error) {
	res := resp.Result()

	if res.StatusCode != wantCode {
		return false, fmt.Errorf("Response.StatusCode want:%d got:%d; resp:%v", wantCode, res.StatusCode, resp)
	}

	ct := res.Header.Get("Content-Type")
	if ct != wantCt {
		return false, fmt.Errorf("Response.ContentType want:%q got:%q", wantCt, ct)
	}

	gotBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	if string(gotBody) != string(wantBody) {
		return false, fmt.Errorf("Response.Body want:%q got:%q", wantBody, gotBody)
	}
	return true, nil
}

var dbPath = store.DBPath

func ensureAccount(username, password string) error {
	a := store.Account{
		Username: username,
		Password: password,
	}
	return a.Create()
}

func ensureToken(tokenID, accountID string) error {
	return store.CreateToken(tokenID, accountID)
}

func ensureFile(accountID, fileID, contents string) error {
	// TODO: implement
	return nil
}

func errorResp(msg string) string {
	return fmt.Sprintf("{\n  \"error\": \"%s\"\n}", msg)
}
