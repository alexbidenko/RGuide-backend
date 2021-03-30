package entities

type Product struct {
	Id int `json:"id"`
	Title string `json:"title" valid:"required,type(string),length(1|255)"`
	Description string `json:"description" valid:"required,type(string)"`
	Specification string `json:"specification" valid:"type(string)"`
	Preview string `json:"preview"`
	Model string `json:"model"`
	CategoryId int `json:"category_id"`
	Category *Category `json:"category"`
	GroupId int `json:"group_id"`
	Group *Group `json:"group"`
}
