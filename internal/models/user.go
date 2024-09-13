package models

// User data model
type User struct {
	Models
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique;not null"`
	Age   int    `json:"age"`
}
