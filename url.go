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

func (u *URL) String() string {
	u.url.RawQuery = u.query.Encode()
	return u.url.String()
}

func (u *URL) Add(key, value string) {
	u.query.Add(key, value)
}

func (u *URL) Del(key string) {
	u.query.Del(key)
}

func (u *URL) Set(key, value string) {
	u.query.Set(key, value)
}

func (u *URL) Get(key string) string {
	return u.query.Get(key)
}

func (u *URL) Query() url.Values {
	return u.query
}

func (u *URL) URL() *url.URL {
	return u.url
}

func (u *URL) JoinPath(p ...string) {
	var np = path.Join(p...)
	u.url.Path = path.Join(u.url.Path, np)
}

func URLEncode(s string) string {
	s = url.QueryEscape(s)
	return s
}

func URLDecode(s string) string {
	s, _ = url.QueryUnescape(s)
	return s
}
