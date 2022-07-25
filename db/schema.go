package db

type User struct {
	Name string `json:"name" gorm:"index"`
	Age  int    `json:"age"`
}
