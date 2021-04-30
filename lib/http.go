package lib

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

// HTTPClient - HTTP Client
type HTTPClient struct {
	Client *http.Client
	URL    string
}

// NewHTTPClient - Returns new HTTPClient
func NewHTTPClient(http *http.Client, url string) *HTTPClient {
	return &HTTPClient{
		Client: http,
		URL:    url,
	}
}

// DownloadFile - Download a file from a URL
func (c *HTTPClient) DownloadFile(fs FileSystem) (string, error) {
	response, err := c.Client.Get(c.URL)
	if err != nil {
		return "", err
	}
	defer func() {
		err = response.Body.Close()
	}()
	if response.StatusCode != 200 {
		return "", errors.New("Received non 200 response code")
	}
	thisUUID := getUUID()
	fileName := fmt.Sprintf("/domi/%s.zip", thisUUID)
	file, err := fs.Create(fileName)
	if err != nil {
		return "", err
	}
	defer func() {
		err = file.Close()
	}()
	_, err = fs.Copy(file, response.Body)
	if err != nil {
		return "", err
	}
	log.Printf("Downloaded %s as %s\n", c.URL, fileName)
	return thisUUID, nil
}
