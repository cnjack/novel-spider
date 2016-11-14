package model

import "github.com/jinzhu/gorm"

type Chapter struct {
	gorm.Model
	NovelID uint   `sql:"novel_id"`
	Index   uint   `sql:"index"`
	Title   string `sql:"title"`
	Data    string `sql:"data" gorm:"type:longtext"`
	Status  uint8  `sql:"status"`
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

func (c *Chapter) Todata() interface{} {
	resp := map[string]interface{}{
		"id":       c.ID,
		"title":    c.Title,
		"data":     c.Data,
		"status":   c.Status,
		"novel_id": c.NovelID,
		"url":      c.Url,
	}
	return resp
}
