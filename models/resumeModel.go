package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Resume struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name             string             `bson:"name" json:"name"`
	Email            string             `bson:"email" json:"email"`
	PhoneNumber      string             `bson:"phone_number" json:"phone_number"`
	LinkedInLink     *string            `bson:"linkedin_link,omitempty" json:"linkedin_link,omitempty"`
	GitHubLink       *string            `bson:"github_link,omitempty" json:"github_link,omitempty"`
	PortfolioLink    *string            `bson:"portfolio_link,omitempty" json:"portfolio_link,omitempty"`
	Location         string             `bson:"location" json:"location"`
	Summary          *string            `bson:"summary,omitempty" json:"summary,omitempty"`
	Skills           []string           `bson:"skills" json:"skills"`
	Education        []Education        `bson:"education" json:"education"`
	WorkExperience   []WorkExperience   `bson:"work_experience" json:"work_experience"`
	Projects         []Project          `bson:"projects" json:"projects"`
	Certifications   []Certification    `bson:"certifications" json:"certifications"`
	Languages        []string           `bson:"languages,omitempty" json:"languages,omitempty"`
	HonorsAwards     []HonorAward       `bson:"honors_awards,omitempty" json:"honors_awards,omitempty"`
	Extracurriculars []Extracurricular  `bson:"extracurriculars,omitempty" json:"extracurriculars,omitempty"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

type Education struct {
	Name      string    `bson:"name" json:"name"`
	Location  string    `bson:"location,omitempty" json:"location"`
	StartDate time.Time `bson:"start_date" json:"start_date"`
	EndDate   time.Time `bson:"end_date" json:"end_date"`
	IsEnrolled bool      `bson:"is_enrolled,omitempty" json:"is_enrolled"`
	ExpectedGraduationDate time.Time `bson:"expected_graduation_date,omitempty" json:"expected_graduation_date,omitempty"`
	GPA       float64   `bson:"gpa,omitempty" json:"gpa,omitempty"`
}

type WorkExperience struct {
	CompanyName  string    `bson:"company_name" json:"company_name"`
	RoleTitle    string    `bson:"role_title" json:"role_title"`
	StartDate    time.Time `bson:"start_date" json:"start_date"`
	EndDate      time.Time `bson:"end_date,omitempty" json:"end_date,omitempty"`
	IsWorking    bool      `bson:"is_working,omitempty" json:"is_working"`
	Location     string    `bson:"location" json:"location"`
	BulletPoints []string  `bson:"bullet_points" json:"bullet_points"`
}

type Project struct {
	Name         string   `bson:"name" json:"name"`
	Description  *string  `bson:"description,omitempty" json:"description,omitempty"`
	ProjectUrl   *string   `bson:"project_url,omitempty" json:"project_url,omitempty"`
	Technologies []string `bson:"technologies,omitempty" json:"technologies,omitempty"`
	BulletPoints []string `bson:"bullet_points" json:"bullet_points"`
	StartDate    time.Time `bson:"start_date" json:"start_date"`
	EndDate      time.Time `bson:"end_date,omitempty" json:"end_date,omitempty"`
}

type Certification struct {
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	CertificateLink string `bson:"certificate_link,omitempty" json:"certificate_link"`
}

type HonorAward struct {
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
}

type Extracurricular struct {
	ActivityName string `bson:"activity_name" json:"activity_name"`
	Description  string `bson:"description" json:"description"`
}
