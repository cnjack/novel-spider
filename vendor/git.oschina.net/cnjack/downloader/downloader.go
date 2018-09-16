package downloader

import (
	"net/http"
)

type Downloader interface {
	Download() Downloader
	Request() *http.Request
	Error() error
	SetProxy(string)
	SetAuthParam(interface{})
	Resource() *Resource
}
