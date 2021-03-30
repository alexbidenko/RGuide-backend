package models

import (
	"rguide/dif"
	"rguide/entities"
)

type CategoryModel struct {
}

func (categoryModel CategoryModel) All() []entities.Category {
	var categories []entities.Category
	dif.DB.Model(&entities.Category{}).Find(&categories)
	return categories
}

func (categoryModel CategoryModel) Create(model *entities.Category) {
	dif.DB.Model(&entities.Category{}).Create(&model)
}
