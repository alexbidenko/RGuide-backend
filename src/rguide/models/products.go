package models

type Product struct {
	Id int `json:"id"`
	Title string `json:"title" valid:"required,type(string),length(1|255)"`
	Description string `json:"description" valid:"required,type(string)"`
	Preview string `json:"preview"`
	Model string `json:"model"`
}
