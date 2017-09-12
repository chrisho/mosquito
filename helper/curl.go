package helper

import (
	"net/http"
	"errors"
	"io/ioutil"
)

func CurlGet(url string) (r string, err error) {
	request := &http.Client{}
	response, err := request.Get(url)
	if err != nil {
		return r, err
	}
	if response.StatusCode != 200 {
		return r, errors.New("request fail")
	}
	byteSlice, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return r, err
	}
	return string(byteSlice), nil
}


