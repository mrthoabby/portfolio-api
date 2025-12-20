package skills

// Category constants
const (
	CategoryBackend    = "backend"
	CategoryFrontend   = "frontend"
	CategoryTools      = "tools"
	CategorySoftSkills = "softSkills"
)

// Proficiency constants
const (
	ProficiencyAdvanced   = "advanced"
	ProficiencyOccasional = "occasional"
	ProficiencyPast       = "past"
)

type Skill struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	ProfileID   string `json:"profileId" bson:"profileId"`
	Name        string `json:"name" bson:"name"`
	Category    string `json:"category" bson:"category"`
	Proficiency string `json:"proficiency" bson:"proficiency"`
}

type Response struct {
	Skills []Skill `json:"skills"`
}

