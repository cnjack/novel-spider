package spider

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis"
)

type IImageSpider interface {
	WriteWithUrl(string, http.ResponseWriter) error
	DefaultImage() []byte
}

type ImageSpider struct {
	httpClient   *http.Client
	storage      *redis.Client
	defaultImage []byte
}

func NewImageSpider(client *redis.Client) IImageSpider {
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer
	defaultTransport.MaxIdleConns = 100
	defaultTransport.MaxIdleConnsPerHost = 10
	img, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAP///wAAACwAAAAAAQABAAACAkQBADs=")
	return &ImageSpider{
		httpClient: &http.Client{
			Transport: &defaultTransport,
			Timeout:   5 * time.Second,
		},
		storage:      client,
		defaultImage: img,
	}
}

func (i *ImageSpider) WriteWithUrl(url string, w http.ResponseWriter) error {
	resp, err := i.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	w.Header().Set("Content-type", resp.Header.Get("Content-type"))
	_, err = io.Copy(w, resp.Body)
	return err
}

func (i *ImageSpider) DefaultImage() []byte {
	return i.defaultImage
}
