package model

import "github.com/jinzhu/gorm"

type Chapter struct {
	ID        uint `gorm:"primary_key"`
	NovelID uint   `sql:"novel_id"`
	Index   uint   `sql:"index"`
	Title   string `sql:"title"`
	Data    string `sql:"data" gorm:"type:longtext"`
	Url     string `sql:"url"`
}

func FirstChapterByID(db *gorm.DB, id uint) (c *Chapter, err error) {
	c = &Chapter{}
	if err = db.Model(c).Where("id = ?", id).First(c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return
}

type ChapterData struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Data    string `json:"data"`
	Status  uint8  `json:"status"`
	NovelID uint   `json:"novel_id"`
	Url     string `json:"url"`
}

func (c *Chapter) Todata() *ChapterData {
	return &ChapterData{
		ID:      c.ID,
		Title:   c.Title,
		Data:    c.Data,
		NovelID: c.NovelID,
		Url:     c.Url,
	}
}
