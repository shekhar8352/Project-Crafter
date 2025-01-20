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
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName       string             `bson:"first_name" json:"first_name" validate:"required,min=2,max=100"`
	LastName        string             `bson:"last_name" json:"last_name" validate:"required,min=2,max=100"`
	DateOfBirth     time.Time          `bson:"date_of_birth" json:"date_of_birth" validate:"required"`
	Password        string             `bson:"password" json:"password" validate:"required,min=6"`
	Email           string             `bson:"email" json:"email" validate:"required,email"`
	UserType        UserType           `bson:"user_type" json:"user_type" validate:"required,oneof=Student Professional"`
	Experience      ExperienceLevel    `bson:"experience_level" json:"experience_level" validate:"required,oneof=Fresher Entry-level Mid-level Senior-level"`
	College         *string            `bson:"college,omitempty" json:"college,omitempty"`
	CurrentCompany  *string            `bson:"current_company,omitempty" json:"current_company,omitempty"`
	ResumeURLs      []string           `bson:"resume_urls,omitempty" json:"resume_urls,omitempty" validate:"dive,url"`
	Token           *string            `bson:"token,omitempty" json:"token,omitempty"`
	RefreshToken    *string            `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}
