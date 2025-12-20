package contacts

import "time"

type Contact struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	ProfileID   string    `json:"profileId" bson:"profileId"`
	Name        string    `json:"name" bson:"name" validate:"required,min=2,max=100"`
	Email       string    `json:"email" bson:"email" validate:"required,email"`
	Message     string    `json:"message" bson:"message" validate:"required,min=10,max=1000"`
	Contacted   bool      `json:"contacted" bson:"contacted"`
	ContactedAt time.Time `json:"contactedAt" bson:"contactedAt"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}

type Request struct {
	Name    string `json:"name" validate:"required,min=2,max=100"`
	Email   string `json:"email" validate:"required,email"`
	Message string `json:"message" validate:"required,min=10,max=1000"`
}

type Response struct {
	ID          string    `json:"id"`
	Message     string    `json:"message"`
	ContactedAt time.Time `json:"contactedAt"`
}
