package spider

import (
	"strings"
	"net/url"
)

type SnwxChapter struct {
	Url    string
	BookID string
	Data   interface{}
}

func (s *SnwxChapter) Name() string {
	return "snwx.com"
}

func (s *SnwxChapter) Match(urlString string) bool {
	s.Url = urlString
	u, err := url.Parse(urlString)
	if err != nil {
		return false
	}
	if u.Host != "www.snwx.com" {
		return false
	}
	u.Path = strings.TrimRight(u.Path, ".html")
	u.Path = strings.Trim(u.Path, `/`)
	paths := strings.Split(u.Path, `/`)
	if len(paths) == 0 {
		return false
	}
	if len(paths) == 3 && paths[0] == "book" {
		s.BookID = paths[1] + "/" + paths[2]
		if err != nil {
			return false
		}
		return true
	}
	return false
	return true
}

func (s *SnwxChapter) Gain() (interface{}, error) {
	return nil, nil
}
