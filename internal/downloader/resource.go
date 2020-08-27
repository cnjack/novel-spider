package downloader

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

type Resource struct {
	resp *http.Response
	doc  *goquery.Document
	body []byte
}

func (r *Resource) analysis() (body []byte, err error) {
	encoding := r.resp.Header.Get(ContentEncoding)
	ctype := r.resp.Header.Get(ContentType)
	var reader io.Reader

	switch encoding {
	case "gzip":
		gzipReader, err := gzip.NewReader(r.resp.Body)
		if err != nil {
			return nil, err
		}
		defer gzipReader.Close()
		reader, err = charset.NewReader(gzipReader, ctype)
	default:
		reader, err = charset.NewReader(r.resp.Body, ctype)
	}
	if err != nil {
		return
	}
	defer r.resp.Body.Close()

	body, err = ioutil.ReadAll(reader)
	r.body = body
	return
}

func (r *Resource) Document() (*goquery.Document, error) {
	//if exist
	if r.doc != nil {
		return r.doc, nil
	}
	body, err := r.analysis()
	if err != nil {
		return nil, err
	}
	r.doc, err = goquery.NewDocumentFromReader(bytes.NewReader(body))
	return r.doc, err
}

func (r *Resource) Body() ([]byte, error) {
	if r.body != nil {
		return r.body, nil
	}
	return r.analysis()
}
