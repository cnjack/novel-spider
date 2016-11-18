package spider

type Novel struct {
	Title        string
	Auth         string
	Style        string
	Introduction string
	From         string
	Status       string
	Cover        string
	Chapter      []*Chapter
}

type Chapter struct {
	Novel *Novel
	Index uint
	Title string
	From  string
	Data  string
}

type Search struct {
	Name       string
	SearchName string
	From       string
}

type Spider interface {
	Name() string
	Match(string) bool
	Gain() (interface{}, error)
}
