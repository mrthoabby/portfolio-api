package profile

import "time"

type Profile struct {
	ID                  string     `json:"id" bson:"_id,omitempty"`
	Name                string     `json:"name" bson:"name"`
	PhotoURL            string     `json:"photoUrl" bson:"photoUrl"`
	ProfessionTittle    string     `json:"title" bson:"title"`
	AboutMe             string     `json:"aboutMe" bson:"aboutMe"`
	FirstExperienceDate *time.Time `json:"firstExperienceDate,omitempty" bson:"firstExperienceDate,omitempty"`
	CreatedAt           time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt" bson:"updatedAt"`
}
