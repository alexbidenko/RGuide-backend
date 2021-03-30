package models

import (
	"rguide/dif"
	"rguide/entities"
)

type ProductModel struct {
}

func (productModel ProductModel) GetById(id int) (entities.Product, error) {
	var product entities.Product
	err := dif.DB.Model(&entities.Product{}).Preload("Category").First(&product, id).Error

	if err != nil {
		return product, err
	}

	return product, nil
}

func (productModel ProductModel) FindAll(q string, parameters map[string]interface{}) []entities.Product {
	var products []entities.Product
	dif.DB.Model(&entities.Product{}).Where("title LIKE ?", "%" + q +  "%").Where(parameters).Find(&products)
	return products
}

func (productModel ProductModel) Create(model *entities.Product) {
	dif.DB.Model(&entities.Product{}).Create(&model)
}

func (productModel ProductModel) Update(id int, model *entities.Product) {
	dif.DB.Model(&entities.Product{}).Where("id = ?", id).Updates(&model)
}
