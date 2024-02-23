package goget

import (
	"net/http"
	"regexp"
	"time"
)

type Option struct {
	F func(o *Options)
}

type Options struct {
	Timeout          time.Duration
	Transport        http.RoundTripper
	ProxyUrl         string
	ProxyRedirectUrl string
}

func (o *Options) Apply(opts []Option) {
	for _, op := range opts {
		op.F(o)
	}
}

func WithProxyUrl(proxyUrl string) Option {
	return Option{F: func(o *Options) {
		o.ProxyUrl = proxyUrl
	}}
}
func WithProxyRedirectUrl(proxyRedirectUrl string) Option {
	return Option{F: func(o *Options) {
		o.ProxyRedirectUrl = proxyRedirectUrl
	}}
}

func WithTimeout(timeout time.Duration) Option {
	return Option{F: func(o *Options) {
		o.Timeout = timeout
	}}
}

func WithTransport(transport http.RoundTripper) Option {
	return Option{F: func(o *Options) {
		o.Transport = transport
	}}
}

var UrlReg = regexp.MustCompile(`http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`)
var DigitReg = regexp.MustCompile(`\d+`)
