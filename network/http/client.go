//Build http client with cookie and proxy

package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type ProxyFunc func(r *http.Request) (*url.URL, error)

type buider struct {
	jar *cookiejar.Jar
	proxy ProxyFunc
}

func ClientBuilder() *buider {
	return &buider{}
}

func (b *buider) SetCookieJar(c *cookiejar.Jar) *buider {
	b.jar = c
	return b
}

func (b *buider) SetProxy(p ProxyFunc) *buider {
	b.proxy = p
	return b
}

func (b *buider) SetHTTPProxy(host string, port int) *buider {
	b.proxy = http.ProxyURL(&url.URL{
		Scheme:      "http",
		Opaque:      "",
		Host:        fmt.Sprintf("%s:%d", host, port),
	})
	return b
}


func (b *buider) SetHTTPProxyWithAuth(host string, port int, username string, password string) *buider {
	var user *url.Userinfo
	if username != "" && password != "" {
		user = url.UserPassword(username, password)
	}
	b.proxy = http.ProxyURL(&url.URL{
		Scheme:      "http",
		Opaque:      "",
		User:        user,
		Host:        fmt.Sprintf("%s:%d", host, port),
	})
	return b
}

func (b *buider) Build() *http.Client{
	if b.jar == nil {
		return &http.Client{
			Transport: &http.Transport{
				Proxy:           b.proxy,
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}
	return &http.Client{
		Transport: &http.Transport{
			Proxy:           b.proxy,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Jar: b.jar,
	}
}
