package projects

import "time"

type Project struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	ProfileID       string    `json:"profileId" bson:"profileId"`
	Name            string    `json:"name" bson:"name"`
	Description     string    `json:"description" bson:"description"`
	TechStack       []string  `json:"techStack" bson:"techStack"`
	GitHubURL       *string   `json:"githubUrl,omitempty" bson:"githubUrl,omitempty"`
	LiveURL         *string   `json:"liveUrl,omitempty" bson:"liveUrl,omitempty"`
	ImageDiagramURL *string   `json:"imageDiagramUrl,omitempty" bson:"imageDiagramUrl,omitempty"`
	Visible         bool      `json:"visible" bson:"visible"`
	CreatedAt       time.Time `json:"createdAt" bson:"createdAt"`
}

type Response struct {
	Projects []Project `json:"projects"`
}
