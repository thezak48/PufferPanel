package pufferpanel

import (
	"github.com/cavaliergopher/grab/v3"
	"github.com/pufferpanel/pufferpanel/v3/files"
	"net/http"
	"os"
)

var httpClient = &http.Client{}

func Http() *http.Client {
	return httpClient
}

func HttpGet(requestUrl string) (*http.Response, error) {
	return httpClient.Get(requestUrl)
}

func HttpExtract(requestUrl, directory string) error {
	//we will write this to temp so we can not keep so much in memory
	response, err := grab.Get(os.TempDir(), requestUrl)
	if err != nil {
		return err
	}
	defer os.Remove(response.Filename)

	err = files.Extract(nil, response.Filename, directory, "*", false, nil)
	return err
}
