package entities

type Group struct {
	Id int `json:"id"`
	Title string `json:"title" valid:"required,type(string),length(1|255)" gorm:"size:255"`
}
