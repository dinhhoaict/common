//Build http client with cookie and proxy

package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Options struct {
	CookieJar *cookiejar.Jar
	Transport *http.Transport
}

type Option func(options *http.Client)

func WithCookie(cookie *cookiejar.Jar) Option {
	return func(args *http.Client) {
		args.Jar = cookie
	}
}

func WithProxy(host string, port int) Option {
	trans := http.Transport{Proxy: http.ProxyURL(&url.URL{
		Scheme: "http",
		Opaque: "",
		Host:   fmt.Sprintf("%s:%d", host, port),
	}),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return WithTransport(&trans)
}

func WithAuthProxy(host string, port int, username string, password string) Option {
	var user *url.Userinfo
	if username != "" && password != "" {
		user = url.UserPassword(username, password)
	}
	proxy := http.ProxyURL(&url.URL{
		Scheme: "http",
		Opaque: "",
		User:   user,
		Host:   fmt.Sprintf("%s:%d", host, port),
	})
	trans := http.Transport{Proxy: proxy, TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	return WithTransport(&trans)
}

func WithTransport(trans *http.Transport) Option {
	return func(args *http.Client) {
		args.Transport = trans
	}
}

func New(options ...Option) *http.Client {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	for _, opt := range options {
		opt(&client)
	}
	return &client
}

type ProxyFunc func(r *http.Request) (*url.URL, error)

type buider struct {
	jar   *cookiejar.Jar
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
		Scheme: "http",
		Opaque: "",
		Host:   fmt.Sprintf("%s:%d", host, port),
	})
	return b
}

func (b *buider) SetHTTPProxyWithAuth(host string, port int, username string, password string) *buider {
	var user *url.Userinfo
	if username != "" && password != "" {
		user = url.UserPassword(username, password)
	}
	b.proxy = http.ProxyURL(&url.URL{
		Scheme: "http",
		Opaque: "",
		User:   user,
		Host:   fmt.Sprintf("%s:%d", host, port),
	})
	return b
}

func (b *buider) Build() *http.Client {
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
