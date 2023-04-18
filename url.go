package nhttp

import (
	"net/url"
	"path"
)

type URL struct {
	url   *url.URL
	query url.Values
}

func NewURL(s string) (u *URL, err error) {
	unescape, err := url.QueryUnescape(s)
	if err != nil {
		return nil, err
	}
	nURL, err := url.Parse(unescape)
	if err != nil {
		return nil, err
	}

	u = &URL{}
	u.url = nURL
	u.query = u.url.Query()
	return u, nil
}

func MustURL(s string) (u *URL) {
	u, err := NewURL(s)
	if err != nil {
		panic(err)
	}
	return u
}

func (this *URL) String() string {
	this.url.RawQuery = this.query.Encode()
	return this.url.String()
}

func (this *URL) Add(key, value string) {
	this.query.Add(key, value)
}

func (this *URL) Del(key string) {
	this.query.Del(key)
}

func (this *URL) Set(key, value string) {
	this.query.Set(key, value)
}

func (this *URL) Get(key string) string {
	return this.query.Get(key)
}

func (this *URL) Query() url.Values {
	return this.query
}

func (this *URL) URL() *url.URL {
	return this.url
}

func (this *URL) JoinPath(p ...string) {
	var np = path.Join(p...)
	this.url.Path = path.Join(this.url.Path, np)
}

func URLEncode(s string) string {
	s = url.QueryEscape(s)
	return s
}

func URLDecode(s string) string {
	s, _ = url.QueryUnescape(s)
	return s
}
