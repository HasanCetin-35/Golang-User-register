package models

type Exercise struct {
	ID string `bson:"_id"`
	Name *string `json:"name" validate:"required,min2,max100"`
	Exercise_Type *string `json:"exercise_type" validate:"required"`
	
}