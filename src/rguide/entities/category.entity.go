package entities

type Category struct {
	Id int `json:"id"`
	Title string `json:"title" valid:"required,type(string),length(1|255)" gorm:"size:255;index:unique" valid:"required,type(string),length(1|255)"`
}

func (faculty *Category) TableName() string {
	return "categories"
}
