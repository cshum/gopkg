package rest

import (
	"io/ioutil"
	"net/http"
	"time"
)

func GetURL(urlstr string, timeout time.Duration) ([]byte, error) {
	req, err := http.NewRequest("GET", urlstr, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	bytes, err := ioutil.ReadAll(resp.Body)
	return bytes, err
}
