package model

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"gitee.com/cnjack/novel-spider/spider"
	"gitee.com/cnjack/novel-spider/spider/kxs"
	"gitee.com/cnjack/novel-spider/spider/snwx"
	"github.com/jinzhu/gorm"
)

type Novel struct {
	gorm.Model
	Title        string `sql:"title" json:"title"`
	Auth         string `sql:"auth" json:"auth"`
	Style        string `sql:"style" json:"style"`
	Status       string `sql:"status" json:"status"`
	Cover        string `sql:"cover" json:"cover"`
	Introduction string `sql:"introduction" gorm:"type:text" json:"intrduction"`
	Chapter      string `sql:"-" json:"-"`
	Url          string `sql:"url" json:"from"`
}

type SearchNovel struct {
	ID    uint   `sql:"id" json:"id"`
	Title string `sql:"title" json:"title"`
	Auth  string `sql:"auth" json:"auth"`
}

var styles []string

func GetStyle(db *gorm.DB) ([]string, error) {
	if len(styles) != 0 {
		return styles, nil
	}
	novels := make([]*Novel, 0)
	err := db.Model(&Novel{}).Select([]string{"style"}).Group("style").Find(&novels).Error
	if err != nil {
		return nil, err
	}
	styles = make([]string, 0)
	for _, v := range novels {
		if v.Style != "" {
			styles = append(styles, v.Style)
		}
	}
	return styles, nil
}

func SearchByTitleOrAuth(db *gorm.DB, title, auth string, op *PageOption) ([]*SearchNovel, error) {
	var ns = make([]*SearchNovel, 0)
	var err error
	if op == nil {
		op = defaultPageOption
	}
	if err = db.Table("novels").Where("title LIKE ? OR auth = ?", "%"+title+"%", auth).Select([]string{"title", "id", "auth"}).Limit(op.Count).Offset(op.Page * op.Count).Order("id desc").Find(&ns).Error; err != nil {
		return nil, err
	}
	return ns, nil
}

func (n *Novel) Add(db *gorm.DB) error {
	return db.Create(n).Error
}

func FirstNovelByID(db *gorm.DB, id uint) (n *Novel, err error) {
	n = &Novel{}
	if err = db.Model(n).Where("id = ?", id).First(n).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return
}

func FirstNovelByUrl(db *gorm.DB, url string) (n *Novel, err error) {
	n = &Novel{}
	if err = db.Model(n).Where("url = ?", url).First(n).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return
}

func FindNovels(db *gorm.DB, op *PageOption) (ns []*Novel, err error) {
	if op == nil {
		op = defaultPageOption
	}
	if err = db.Model(&Novel{}).Order("id " + op.Sort).Limit(op.Count).Offset(op.Page * op.Count).Find(&ns).Error; err != nil {
		return nil, err
	}
	return
}

func FindNovelsWithStyle(db *gorm.DB, style string, op *PageOption) (ns []*Novel, err error) {
	if op == nil {
		op = defaultPageOption
	}
	if err = db.Model(&Novel{}).Where("style LIKE '%" + style + "%'").Order("id " + op.Sort).Limit(op.Count).Offset(op.Page * op.Count).Find(&ns).Error; err != nil {
		return nil, err
	}
	return
}

func FirstChapterByID(db *gorm.DB, id uint) ([]*NovelChapter, error) {
	n := &Novel{}
	if err := db.Model(&Novel{}).Select([]string{"id", "url"}).Where("id = ?", id).Limit(1).Find(n).Error; err != nil {
		return nil, err
	}
	return n.GetChapter()
}

type NovelData struct {
	ID           uint            `json:"id"`
	CreateAt     string          `json:"create_at"`
	Title        string          `json:"title"`
	Auth         string          `json:"auth"`
	Style        string          `json:"style"`
	Status       string          `json:"status"`
	Cover        string          `json:"cover"`
	Introduction string          `json:"introduction"`
	Chapter      []*NovelChapter `json:"chapters"`
	Url          string          `json:"from"`
}

func (n *Novel) Todata(more bool) *NovelData {
	cover, err := url.Parse(n.Cover)
	coverStr := ""
	if err != nil || cover.Host != "spider-img.nightc.com" {
		coverStr = "http://spider-img.nightc.com/cover.jpg"
	} else {
		coverStr = cover.String()
	}
	resp := NovelData{
		ID:           n.ID,
		CreateAt:     n.CreatedAt.Format(time.RFC3339),
		Title:        n.Title,
		Auth:         n.Auth,
		Style:        n.Style,
		Status:       n.Status,
		Introduction: n.Introduction,
		Cover:        coverStr,
		Url:          n.Url,
	}
	if more {
		resp.Chapter, _ = n.GetChapter()
	}
	return &resp
}

type NovelChapter struct {
	Title string `json:"title"`
	Index uint   `json:"index"`
	Url   string `json:"url"`
}

const (
	ChaptersRedisKey = "chapters_"
	ChapterRedisKey  = "chapter_"
)

func (n *Novel) GetChapter() ([]*NovelChapter, error) {
	redisClient := MustGetRedisClient()
	key := ChaptersRedisKey + fmt.Sprintf("%x", md5.Sum([]byte(n.Url)))
	n.Chapter = redisClient.Get(key).Val()
	if n.Chapter == "" {
		var chaptersSpider spider.Spider
		for _, v := range []spider.Spider{
			&snwx.Novel{},
			&kxs.Novel{},
		} {
			if v.Match(n.Url) {
				chaptersSpider = v
				break
			}
		}
		if chaptersSpider != nil {
			novel, err := chaptersSpider.Gain()
			if err != nil {
				return nil, err
			}
			if novel != nil {
				novelChapters := make([]*NovelChapter, 0)
				for key, v := range novel.(spider.Novel).Chapter {
					novelChapters = append(novelChapters, &NovelChapter{
						Title: v.Title,
						Index: uint(key),
						Url:   v.From,
					})
				}
				chapterStr, _ := json.Marshal(novelChapters)
				redisClient.Set(key, chapterStr, 3*time.Hour)
				return novelChapters, nil
			}
		}
	}
	novelChapters := make([]*NovelChapter, 0)
	if n.Chapter != "" {
		err := json.Unmarshal([]byte(n.Chapter), &novelChapters)
		if err != nil {
			return nil, err
		}
	}
	return novelChapters, nil
}

func GetNovelFromUrl(url string) (*Novel, error) {
	novel := &Novel{}
	var novelSpider spider.Spider
	for _, v := range []spider.Spider{
		&snwx.Novel{},
		&kxs.Novel{},
	} {
		if v.Match(url) {
			novelSpider = v
			break
		}
	}
	if novelSpider != nil {
		novelSrc, err := novelSpider.Gain()
		if err != nil {
			return nil, err
		}
		if novelSrc != nil {
			novel.Url = url
			novel.Style = novelSrc.(spider.Novel).Style
			novel.Title = novelSrc.(spider.Novel).Title
			novel.Cover = novelSrc.(spider.Novel).Cover
			novel.Status = novelSrc.(spider.Novel).Status
			novel.Auth = novelSrc.(spider.Novel).Auth
			novel.Introduction = novelSrc.(spider.Novel).Introduction
		}

	}
	err := db.Model(&Novel{}).Save(novel).Error
	return novel, err
}

func GetChapter(url string) (string, error) {
	redisClient := MustGetRedisClient()
	key := ChapterRedisKey + fmt.Sprintf("%x", md5.Sum([]byte(url)))
	if chapter := redisClient.Get(key).Val(); chapter != "" {
		return chapter, nil
	}
	var chapterSpider spider.Spider
	for _, v := range []spider.Spider{
		&snwx.Chapter{},
		&kxs.Chapter{},
	} {
		if v.Match(url) {
			chapterSpider = v
			break
		}
	}
	if chapterSpider != nil {
		chapter, err := chapterSpider.Gain()
		if err != nil {
			return "", err
		}
		if chapter != nil {
			chapterStr := chapter.(string)
			redisClient.Set(key, chapterStr, 7*24*time.Hour)
			return chapterStr, nil
		}
	}
	return "", errors.New("spider err")
}
