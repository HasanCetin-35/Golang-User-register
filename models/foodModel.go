package models

type Food struct {
	ID    string   `bson:"_id"`
	Name  *string  `json:"name" validate:"required,min=2,max=100"`
	Price *float64 `json:"price" validate:"required"`
}
