package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExperienceLevel string
type UserType string

const (
	Fresher     ExperienceLevel = "Fresher"
	EntryLevel  ExperienceLevel = "Entry-level"
	MidLevel    ExperienceLevel = "Mid-level"
	SeniorLevel ExperienceLevel = "Senior-level"

	Student      UserType = "Student"
	Professional UserType = "Professional"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id"`
	First_name      *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_name       *string            `json:"last_name" validate:"required,min=2,max=100"`
	Date_of_birth   *string            `json:"date_of_birth" validate:"required"`
	Password        *string            `json:"password" validate:"required,min=6"`
	Email           *string            `json:"email" validate:"required"`
	UserType        UserType           `json:"user_type" validate:"required,oneof=Student Professional"`
	Experience      ExperienceLevel    `json:"experience_level" validate:"required,oneof=Fresher Entry-level Mid-level Senior-level"`
	College         *string            `json:"college"`
	Current_company *string            `json:"current_company"`
	ResumeURLs      []string           `json:"resume_urls" validate:"dive,url"`

	Token         *string   `json:"token"`
	Refresh_Token *string   `json:"refresh_token"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}
