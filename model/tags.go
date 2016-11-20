package model

import "github.com/jinzhu/gorm"

type Tags struct {
	ID      int    `sql:"id" json:"id"`
	TagName string `sql:"tag_name" json:"tag_name"`
}

func GetTags(db *gorm.DB) (*[]Tags, error) {
	var tags []Tags
	if err := db.Model(&Tags{}).Find(&tags).Error; err != nil {
		return nil, err
	}
	return &tags, nil
}

func FirstTagsByID(db *gorm.DB, id int) (*Tags, error) {
	tag := &Tags{}
	if err := db.Model(&Tags{}).Where("id = ?", id).First(&tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}
