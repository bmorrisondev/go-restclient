package restclient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Options struct {
	UrlBase  *string
	Method   string
	Headers  map[string]string
	Body     *string
	FormData *url.Values
}

type Response struct {
	RawResponse *http.Response
}

func (r *Response) MarshalJson(out interface{}) error {
	b, err := io.ReadAll(r.RawResponse.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &out)
	if err != nil {
		return err
	}
	return nil
}

func (r *Response) Text() (*string, error) {
	b, err := io.ReadAll(r.RawResponse.Body)
	if err != nil {
		return nil, err
	}
	str := string(b)
	return &str, nil
}

func Execute(url string, options *Options) (*Response, error) {
	if options == nil {
		options = &Options{
			Method: "GET",
		}
	}

	var req *http.Request
	if options.UrlBase != nil {
		url = *options.UrlBase + url
	}
	if options.Body != nil {
		bodyReader := strings.NewReader(*options.Body)
		r, err := http.NewRequest(options.Method, url, bodyReader)
		if err != nil {
			return nil, err
		}
		req = r
	} else if options.FormData != nil {
		data := *options.FormData
		r, err := http.NewRequest(options.Method, url, strings.NewReader(data.Encode()))
		if err != nil {
			return nil, err
		}
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req = r
	} else {
		r, err := http.NewRequest(options.Method, url, nil)
		if err != nil {
			return nil, err
		}
		req = r
	}

	if options.Headers != nil {
		for key, val := range options.Headers {
			req.Header.Add(key, val)
		}
	}

	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	fr := Response{
		RawResponse: res,
	}

	return &fr, nil
}
