package dif

import "rguide/entities"

func Migrate() {
	DB.AutoMigrate(&entities.Category{}, &entities.Product{}, &entities.Group{})
}
