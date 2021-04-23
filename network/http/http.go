package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func buildRequest(method string, url string, header map[string]string, query map[string]string, body io.Reader) (*http.Request, error){
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("cant create http request: %s", err.Error())
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	if query != nil {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}
	return req, nil
}

func Get(c *http.Client, url string, header map[string]string, query map[string]string) (*http.Response, error){
	req, err := buildRequest(http.MethodGet, url, header, query, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func PostJson(c *http.Client, url string, header map[string]string, query map[string]string, body interface{}) (*http.Response, error){
	bodyValue,err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("cant encode json: %s", err.Error())
	}
	req, err := buildRequest(http.MethodPost, url, header, query, bytes.NewReader(bodyValue))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	return c.Do(req)
}

func PostForm(c *http.Client, urlStr string, header map[string]string, query map[string]string, formData url.Values) (*http.Response, error){
	req, err := buildRequest(http.MethodPost, urlStr, header, query, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.Do(req)
}

func ReadBody(resp *http.Response) (string, error){
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}