package questions

import "time"

type Question struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	ProfileID string    `json:"profileId" bson:"profileId"`
	Message   string    `json:"message" bson:"message"`
	IP        string    `json:"ip" bson:"ip"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type Request struct {
	Message string `json:"message" validate:"required,min=5,max=500"`
}

type Response struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}
