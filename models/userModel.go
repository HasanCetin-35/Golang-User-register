package models

type User struct {
	ID       string `bson:"_id"`
	Email    string `json:"email" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=5,max=20"`
}
