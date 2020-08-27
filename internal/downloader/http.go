package downloader

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

const (
	ContentType     string = `Content-Type`
	ContentEncoding string = `Content-Encoding`
)

type HttpDownloader struct {
	req  *http.Request
	e    error
	p    *url.URL
	auth interface{}
	res  *Resource
}

func NewHttpDownloaderFromUrl(u *url.URL) Downloader {
	header := make(http.Header)
	header.Add("Connection", "keep-alive")
	header.Add("Cache-Control", "max-age=0")
	req, err := http.NewRequest("GET", u.String(), nil)
	d := &HttpDownloader{
		req: req,
		e:   err,
	}
	return d
}

func NewHttpDownloaderFromRequest(req *http.Request) Downloader {
	d := &HttpDownloader{
		req: req,
	}
	return d
}

//SetProxy set the http proxy, param proxyUrl: the proxy host
func (d *HttpDownloader) SetProxy(proxyUrl string) {
	proxy, err := url.Parse(proxyUrl)
	if err != nil {
		return
	}
	d.p = proxy
}

//SetAuthParam set the cookieJar
func (d *HttpDownloader) SetAuthParam(param interface{}) {
	d.auth = param
}

func (d *HttpDownloader) Request() *http.Request {
	return d.req
}

func (d *HttpDownloader) Download() Downloader {
	resp, err := d.download()
	if err != nil {
		d.e = err
	}
	d.res = &Resource{resp: resp}
	return d
}

func (d *HttpDownloader) Error() error {
	return d.e
}

func (d *HttpDownloader) download() (resp *http.Response, err error) {
	client := &http.Client{
		// set the timeout 5s
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	//check cookie
	if d.auth != nil {
		cookieJar, ok := d.auth.(http.CookieJar)
		if ok {
			client.Jar = cookieJar
		}
	}
	//check proxy
	if d.p != nil {
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(d.p),
		}
	}
	resp, err = client.Do(d.req)
	return
}

func (d *HttpDownloader) Resource() *Resource {
	return d.res
}
