package util

import (
	"time"
	"io"
	"compress/gzip"
	"io/ioutil"
	"net/url"
	"bytes"
	"net/http"
	"github.com/kusora/dlog"
)

func HttpPostUrlValuesRawResult(c *http.Client, link string, params url.Values) (int, []byte, error) {
	//defer statsdTrace(link)()
	start := time.Now()
	postDataStr := params.Encode()
	dlog.Println(link, postDataStr)
	postDataBytes := []byte(postDataStr)
	reqest, err := http.NewRequest("POST", link, bytes.NewReader(postDataBytes))
	if err != nil {
		return 200, nil, err
	}

	reqest.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	response, err := c.Do(reqest)

	if err != nil {
		if response == nil {
			return http.StatusInternalServerError, nil, err
		}
		return response.StatusCode, nil, err
	}

	defer func() {
		if response != nil && response.Body != nil {
			response.Body.Close()
		}
		dlog.Info("call [%s] in [%d] nanoseconds", link, time.Since(start).Nanoseconds())
	}()
	if response.Body == nil {
		return http.StatusInternalServerError, nil, err
	}
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(response.Body)
		defer reader.Close()
	default:
		reader = response.Body
	}

	html, err := ioutil.ReadAll(reader)

	dlog.Info("response %s, request %s", string(html), link)
	if err != nil {
		return response.StatusCode, nil,err
	}
	return response.StatusCode, html, nil
}
